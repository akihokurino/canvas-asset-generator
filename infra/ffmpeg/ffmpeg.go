package ffmpeg

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pkg/errors"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type Client interface {
	Video2Thumbnail(fromPath string, frameNum int) (*bytes.Buffer, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) Video2Thumbnail(fromPath string, frameNum int) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	if err := ffmpeg_go.Input(fromPath).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run(); err != nil {
		return nil, errors.WithStack(err)
	}

	return buf, nil
}
