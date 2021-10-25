package usecase

import (
	"bufio"
	"bytes"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/ffmpeg"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type SplitVideo func(ctx context.Context, path string) error

func NewSplitVideo(gcsClient cloud_storage.Client, ffmpegClient ffmpeg.Client) SplitVideo {
	return func(ctx context.Context, path string) error {
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

		for i := 0; i < durationSecond; i++ {
			thumbnail, err := ffmpegClient.Video2Thumbnail(tmpPath, i)
			if err != nil {
				return errors.WithStack(err)
			}

			imgSource, _, err := image.Decode(thumbnail)
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
		}

		return nil
	}
}
