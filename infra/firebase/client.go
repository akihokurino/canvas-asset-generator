package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
)

type UID string

func (id UID) String() string {
	return string(id)
}

type Client interface {
	VerifyToken(ctx context.Context, token string) (UID, error)
}

func NewClient() Client {
	return &client{}
}

type client struct {
}

func (c *client) initApp(ctx context.Context) (*firebase.App, error) {
	var app *firebase.App
	var err error

	conf := &firebase.Config{}

	app, err = firebase.NewApp(ctx, conf)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return app, nil
}

func (c *client) authClient(ctx context.Context) (*auth.Client, error) {
	app, err := c.initApp(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client, nil
}

func (c *client) VerifyToken(ctx context.Context, token string) (UID, error) {
	cli, err := c.authClient(ctx)
	if err != nil {
		return "", err
	}

	decoded, err := cli.VerifyIDToken(ctx, token)
	if err != nil {
		return "", err
	}

	return UID(decoded.UID), nil
}
