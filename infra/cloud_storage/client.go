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
	sc, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = sc.Close()
	}()

	results := make([]string, 0)
	it := sc.Bucket(c.bucketName).Objects(ctx, &storage.Query{Prefix: path})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		results = append(results, attrs.Name)
	}

	return results, nil
}

func (c *client) Download(ctx context.Context, path string) (*bytes.Buffer, error) {
	sc, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = sc.Close()
	}()

	reader, err := sc.Bucket(c.bucketName).Object(path).NewReader(ctx)
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
	s, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	writer := s.Bucket(c.bucketName).Object(path).NewWriter(ctx)
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
	sc, err := storage.NewClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = sc.Close()
	}()

	if err := sc.Bucket(c.bucketName).Object(path).Delete(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
