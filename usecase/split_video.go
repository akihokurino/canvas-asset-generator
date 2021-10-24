package usecase

import (
	"bufio"
	"canvas-server/infra/cloud_storage"
	"context"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

type SplitVideo func(ctx context.Context, path string) error

func NewSplitVideo(gcsClient cloud_storage.Client) SplitVideo {
	return func(ctx context.Context, path string) error {
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
		log.Println(tmpPath)

		return nil
	}
}
