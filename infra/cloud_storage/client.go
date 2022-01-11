package cloud_storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"google.golang.org/api/iterator"

	"github.com/pkg/errors"

	"cloud.google.com/go/storage"
)

type Client interface {
	List(ctx context.Context, bucket string, path string) ([]string, error)
	Download(ctx context.Context, bucket string, path string) (*bytes.Buffer, error)
	Save(ctx context.Context, bucket string, path string, data []byte, contentType string) (*url.URL, error)
	Delete(ctx context.Context, bucket string, path string) error
	FullPath(bucket string, path string) string
	Signature(gsURL *url.URL) (*url.URL, error)
}

type client struct {
	projectID         string
	encodedPrivateKey string
}

func NewClient(
	projectID string,
	encodedPrivateKey string) Client {
	return &client{
		projectID:         projectID,
		encodedPrivateKey: encodedPrivateKey,
	}
}

func (c *client) List(ctx context.Context, bucket string, path string) ([]string, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = cli.Close()
	}()

	results := make([]string, 0)
	it := cli.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: path})
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

func (c *client) Download(ctx context.Context, bucket string, path string) (*bytes.Buffer, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = cli.Close()
	}()

	reader, err := cli.Bucket(bucket).Object(path).NewReader(ctx)
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

func (c *client) Save(ctx context.Context, bucket string, path string, data []byte, contentType string) (*url.URL, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	writer := cli.Bucket(bucket).Object(path).NewWriter(ctx)
	defer func() {
		_ = writer.Close()
	}()

	writer.ContentType = contentType

	if _, err := writer.Write(data); err != nil {
		return nil, errors.WithStack(err)
	}

	u, err := url.Parse(fmt.Sprintf("gs://%s/%s", bucket, path))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u, err
}

func (c *client) Delete(ctx context.Context, bucket string, path string) error {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = cli.Close()
	}()

	if err := cli.Bucket(bucket).Object(path).Delete(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) FullPath(bucket string, path string) string {
	return fmt.Sprintf("gs://%s/%s", bucket, path)
}

func (c *client) Signature(gsURL *url.URL) (*url.URL, error) {
	if gsURL == nil {
		return nil, nil
	}

	if gsURL.String() == "" {
		return nil, nil
	}

	paths := strings.Split(gsURL.Path, "/")
	bucketID := gsURL.Host
	objectID := strings.Join(paths[1:], "/")

	expires := time.Now().Add(time.Hour * 1)

	privateKey, _ := base64.StdEncoding.DecodeString(c.encodedPrivateKey)

	urlStringWithSignature, err := storage.SignedURL(bucketID, objectID, &storage.SignedURLOptions{
		GoogleAccessID: fmt.Sprintf("%s@appspot.gserviceaccount.com", c.projectID),
		PrivateKey:     privateKey,
		Method:         "GET",
		Expires:        expires,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := gsURL.Parse(urlStringWithSignature)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return urlWithSignature, nil
}
