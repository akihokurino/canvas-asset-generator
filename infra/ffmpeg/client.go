package ffmpeg

import (
	"bytes"
	"context"

	"gopkg.in/vansante/go-ffprobe.v2"

	"github.com/pkg/errors"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type Client interface {
	Video2Thumbnail(fromPath string, start int) (*bytes.Buffer, error)
	DurationSecond(ctx context.Context, fromPath string) (int, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) Video2Thumbnail(fromPath string, start int) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	if err := ffmpeg_go.Input(fromPath, ffmpeg_go.KwArgs{"ss": start}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run(); err != nil {
		return nil, errors.WithStack(err)
	}

	return buf, nil
}

func (c *client) DurationSecond(ctx context.Context, fromPath string) (int, error) {
	metadata, err := ffprobe.ProbeURL(ctx, fromPath)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	result := int(metadata.Format.Duration().Seconds())

	return result, nil
}
