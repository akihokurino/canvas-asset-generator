package subscriber

import (
	"bufio"
	"bytes"
	"canvas-asset-generator/config"
	"canvas-asset-generator/infra/cloud_storage"
	"canvas-asset-generator/infra/datastore"
	"canvas-asset-generator/infra/datastore/fcm_token"
	"canvas-asset-generator/infra/datastore/frame"
	"canvas-asset-generator/infra/datastore/work"
	"canvas-asset-generator/infra/ffmpeg"
	"canvas-asset-generator/infra/firebase"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/image/draw"

	"go.mercari.io/datastore/boom"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type SplitVideo func(w http.ResponseWriter, r *http.Request)

func NewSplitVideo(
	gcsClient cloud_storage.Client,
	ffmpegClient ffmpeg.Client,
	fireClient firebase.Client,
	tx datastore.Transaction,
	workRepo work.Repository,
	frameRepo frame.Repository,
	fcmTokenRepo fcm_token.Repository) SplitVideo {
	split := func(ctx context.Context, path string) error {
		now := time.Now()
		videoName := strings.Replace(path, ".mp4", "", -1)

		log.Printf("video path = %s", path)
		log.Printf("video name = %s", videoName)

		log.Printf("---------- download video ----------")

		buf, err := gcsClient.Download(ctx, config.VideoBucketName, path)
		if err != nil {
			return errors.WithStack(err)
		}

		file, err := os.CreateTemp("", "video")
		if err != nil {
			return errors.WithStack(err)
		}

		writer := bufio.NewWriter(file)
		if _, err := writer.Write(buf.Bytes()); err != nil {
			return errors.WithStack(err)
		}

		log.Printf("---------- delete current frames ----------")

		currents, err := gcsClient.List(ctx, config.FrameBucketName, videoName)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, current := range currents {
			log.Printf("deleted frame path = %s", current)
			if err := gcsClient.Delete(ctx, config.FrameBucketName, current); err != nil {
				return errors.WithStack(err)
			}
		}

		log.Printf("---------- prepare split video ----------")

		tmpPath := file.Name()
		durationSecond, err := ffmpegClient.DurationSecond(ctx, tmpPath)
		if err != nil {
			return errors.WithStack(err)
		}
		totalSecond := int(math.Min(float64(durationSecond), 30))

		log.Printf("total second %d", totalSecond)

		workEntity := work.NewEntity(videoName, gcsClient.FullPath(config.VideoBucketName, path), now)
		frameEntities := make([]*frame.Entity, 0)

		for i := 0; i < totalSecond; i++ {
			log.Printf("---------- split video of %d ----------", i)

			data, err := ffmpegClient.Video2Thumbnail(tmpPath, i)
			if err != nil {
				return errors.WithStack(err)
			}

			imgSource, _, err := image.Decode(data)
			if err != nil {
				return errors.WithStack(err)
			}

			// 正方形で切り取る
			width := imgSource.Bounds().Dx()
			height := imgSource.Bounds().Dy()
			length := math.Min(float64(width), float64(height))
			startX := (width - int(length)) / 2
			startY := (height - int(length)) / 2
			endX := startX + int(length)
			endY := startY + int(length)
			subImage := imgSource.(SubImager).SubImage(image.Rect(startX, startY, endX, endY))

			rect := subImage.Bounds()
			resizedImage := image.NewRGBA(image.Rect(0, 0, rect.Dx()/5, rect.Dy()/5))
			draw.CatmullRom.Scale(resizedImage, resizedImage.Bounds(), subImage, rect, draw.Over, nil)

			orgBuf := bytes.NewBuffer(nil)
			if err := jpeg.Encode(orgBuf, subImage, &jpeg.Options{Quality: 100}); err != nil {
				return errors.WithStack(err)
			}

			orgURL, err := gcsClient.Save(
				ctx,
				config.FrameBucketName,
				fmt.Sprintf("%s/%d", videoName, i),
				orgBuf.Bytes(),
				"image/jpeg")
			if err != nil {
				return errors.WithStack(err)
			}

			resizedBuf := bytes.NewBuffer(nil)
			if err := jpeg.Encode(resizedBuf, resizedImage, &jpeg.Options{Quality: 80}); err != nil {
				return errors.WithStack(err)
			}

			resizedURL, err := gcsClient.Save(
				ctx,
				config.FrameBucketName,
				fmt.Sprintf("%s/%d/resized", videoName, i),
				resizedBuf.Bytes(),
				"image/jpeg")
			if err != nil {
				return errors.WithStack(err)
			}

			log.Printf("splited frame url %s", orgURL.String())
			frameEntities = append(frameEntities, frame.NewEntity(workEntity.ID, orgURL, resizedURL, i, now))
		}

		currentFrameEntities, err := frameRepo.GetAllByWork(ctx, workEntity.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		log.Printf("---------- delete and create datastore entities ----------")

		eg := errgroup.Group{}

		for i := range currentFrameEntities {
			e := currentFrameEntities[i]
			eg.Go(func() error {
				return tx(ctx, func(tx *boom.Transaction) error {
					return frameRepo.Delete(tx, e.ID)
				})
			})
		}

		eg.Go(func() error {
			return tx(ctx, func(tx *boom.Transaction) error {
				return workRepo.Put(tx, workEntity)
			})
		})

		for i := range frameEntities {
			e := frameEntities[i]
			eg.Go(func() error {
				return tx(ctx, func(tx *boom.Transaction) error {
					return frameRepo.Put(tx, e)
				})
			})
		}

		if err := eg.Wait(); err != nil {
			return errors.WithStack(err)
		}

		log.Printf("---------- send push notification for complete ----------")

		tokens, err := fcmTokenRepo.GetAll(ctx)
		if err != nil {
			log.Printf("failed get fcm tokens, %+v", err)
			return nil
		}

		pushBody := fmt.Sprintf("%sのフレームの生成が完了しました", workEntity.ID)
		for _, token := range tokens {
			if err := fireClient.SendPushNotification(ctx, token.Token, "", pushBody, 0, map[string]string{}, func(t string) {
				_ = tx(ctx, func(tx *boom.Transaction) error {
					return fcmTokenRepo.Delete(tx, token.ID)
				})
			}); err != nil {
				log.Printf("failed send push, %+v", err)
			}
		}

		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type Payload struct {
			Path string `json:"path"`
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ReadAll: %v", err)
			http.Error(w, "Internal Error, cannot read body", http.StatusInternalServerError)
			return
		}

		var payload Payload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Unmarshal: %v", err)
			http.Error(w, "Internal Error, cannot parse json body", http.StatusInternalServerError)
			return
		}

		if err := split(ctx, payload.Path); err != nil {
			log.Printf("SplitVideo: %v", err)
			http.Error(w, "Internal Error, cannot split video", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
