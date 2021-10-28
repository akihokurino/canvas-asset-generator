package cloud_storage

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"google.golang.org/api/iterator"

	"github.com/pkg/errors"

	"cloud.google.com/go/storage"
)

type Client interface {
	List(ctx context.Context, path string) ([]string, error)
	Download(ctx context.Context, path string) (*bytes.Buffer, error)
	Save(ctx context.Context, path string, data []byte) (*url.URL, error)
	Delete(ctx context.Context, path string) error
	FullPath(path string) string
}

type client struct {
	bucketName string
}

func NewClient(bucketName string) Client {
	return &client{
		bucketName: bucketName,
	}
}

func (c *client) List(ctx context.Context, path string) ([]string, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = cli.Close()
	}()

	results := make([]string, 0)
	it := cli.Bucket(c.bucketName).Objects(ctx, &storage.Query{Prefix: path})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		results = append(results, attrs.Name)
	}

	return results, nil
}

func (c *client) Download(ctx context.Context, path string) (*bytes.Buffer, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = cli.Close()
	}()

	reader, err := cli.Bucket(c.bucketName).Object(path).NewReader(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = reader.Close()
	}()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, errors.WithStack(err)
	}

	return &buf, nil
}

func (c *client) Save(ctx context.Context, path string, data []byte) (*url.URL, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	writer := cli.Bucket(c.bucketName).Object(path).NewWriter(ctx)
	defer func() {
		_ = writer.Close()
	}()

	writer.ContentType = "image/jpeg"

	if _, err := writer.Write(data); err != nil {
		return nil, errors.WithStack(err)
	}

	u, err := url.Parse(fmt.Sprintf("gs://%s/%s", c.bucketName, path))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u, err
}

func (c *client) Delete(ctx context.Context, path string) error {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = cli.Close()
	}()

	if err := cli.Bucket(c.bucketName).Object(path).Delete(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) FullPath(path string) string {
	return fmt.Sprintf("gs://%s/%s", c.bucketName, path)
}
