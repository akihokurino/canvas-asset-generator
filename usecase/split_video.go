package usecase

import (
	"bufio"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/ffmpeg"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type SplitVideo func(ctx context.Context, path string) error

func NewSplitVideo(gcsClient cloud_storage.Client, ffmpegClient ffmpeg.Client) SplitVideo {
	return func(ctx context.Context, path string) error {
		videoName := strings.Replace(path, "Video/", "", -1)
		videoName = strings.Replace(videoName, ".mp4", "", -1)
		log.Println(videoName)

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

		thumbnail, err := ffmpegClient.Video2Thumbnail(tmpPath, 1)
		if err != nil {
			return errors.WithStack(err)
		}

		u, err := gcsClient.Save(ctx, fmt.Sprintf("Thumbnail/%s", videoName), thumbnail.Bytes())
		if err != nil {
			return errors.WithStack(err)
		}

		log.Printf("url %s", u.String())

		return nil
	}
}
