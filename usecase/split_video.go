package usecase

import (
	"bufio"
	"bytes"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/fcm_token"
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"canvas-server/infra/ffmpeg"
	"canvas-server/infra/firebase"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"go.mercari.io/datastore/boom"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type SplitVideo func(ctx context.Context, path string) error

func NewSplitVideo(
	gcsClient cloud_storage.Client,
	ffmpegClient ffmpeg.Client,
	fireClient firebase.Client,
	tx datastore.Transaction,
	workRepo work.Repository,
	thumbnailRepo thumbnail.Repository,
	fcmTokenRepo fcm_token.Repository) SplitVideo {
	return func(ctx context.Context, path string) error {
		now := time.Now()
		videoName := strings.Replace(path, "Video/", "", -1)
		videoName = strings.Replace(videoName, ".mp4", "", -1)

		log.Printf("video path = %s", path)
		log.Printf("video name = %s", videoName)

		buf, err := gcsClient.Download(ctx, path)
		if err != nil {
			return errors.WithStack(err)
		}

		file, err := ioutil.TempFile("", "video")
		if err != nil {
			return errors.WithStack(err)
		}

		writer := bufio.NewWriter(file)

		if _, err := writer.Write(buf.Bytes()); err != nil {
			return errors.WithStack(err)
		}

		tmpPath := file.Name()

		durationSecond, err := ffmpegClient.DurationSecond(ctx, tmpPath)
		if err != nil {
			return errors.WithStack(err)
		}

		currents, err := gcsClient.List(ctx, fmt.Sprintf("Thumbnail/%s", videoName))
		if err != nil {
			return errors.WithStack(err)
		}

		for _, current := range currents {
			log.Printf("deleted thumbnail path = %s", current)
			if err := gcsClient.Delete(ctx, current); err != nil {
				return errors.WithStack(err)
			}
		}

		workEntity := work.NewEntity(videoName, gcsClient.FullPath(path), now)
		thumbnailEntities := make([]*thumbnail.Entity, 0)

		for i := 0; i < durationSecond; i++ {
			data, err := ffmpegClient.Video2Thumbnail(tmpPath, i)
			if err != nil {
				return errors.WithStack(err)
			}

			imgSource, _, err := image.Decode(data)
			if err != nil {
				return errors.WithStack(err)
			}

			subImage := imgSource.(SubImager).SubImage(
				image.Rect(0, 250, imgSource.Bounds().Dx(), imgSource.Bounds().Dy()),
			)

			buf := bytes.NewBuffer(nil)
			if err := jpeg.Encode(buf, subImage, &jpeg.Options{Quality: 100}); err != nil {
				return errors.WithStack(err)
			}

			u, err := gcsClient.Save(ctx, fmt.Sprintf("Thumbnail/%s/%d", videoName, i), buf.Bytes())
			if err != nil {
				return errors.WithStack(err)
			}

			log.Printf("thumbnail url %s", u.String())
			thumbnailEntities = append(thumbnailEntities, thumbnail.NewEntity(workEntity.ID, u.String(), i, now))
		}

		currentThumbnailEntities, err := thumbnailRepo.GetAllByWork(ctx, workEntity.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		eg := errgroup.Group{}

		for i := range currentThumbnailEntities {
			e := currentThumbnailEntities[i]
			eg.Go(func() error {
				return tx(ctx, func(tx *boom.Transaction) error {
					return thumbnailRepo.Delete(tx, e.ID)
				})
			})
		}

		eg.Go(func() error {
			return tx(ctx, func(tx *boom.Transaction) error {
				return workRepo.Put(tx, workEntity)
			})
		})

		for i := range thumbnailEntities {
			e := thumbnailEntities[i]
			eg.Go(func() error {
				return tx(ctx, func(tx *boom.Transaction) error {
					return thumbnailRepo.Put(tx, e)
				})
			})
		}

		if err := eg.Wait(); err != nil {
			return errors.WithStack(err)
		}

		tokens, err := fcmTokenRepo.GetAll(ctx)
		if err != nil {
			log.Printf("failed get fcm tokens, %+v", err)
			return nil
		}

		pushBody := fmt.Sprintf("%sのサムネイルの生成が完了しました", workEntity.ID)
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
}
