package firebase

import (
	"context"

	"firebase.google.com/go/messaging"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
)

type UID string

func (id UID) String() string {
	return string(id)
}

type PushPayload map[string]string

func (p PushPayload) IOS() map[string]interface{} {
	converted := map[string]interface{}{}
	for k, v := range p {
		converted[k] = v
	}
	return converted
}

type Client interface {
	VerifyToken(ctx context.Context, token string) (UID, error)
	SendPushNotification(
		ctx context.Context,
		token string,
		title string,
		body string,
		badge int,
		payload PushPayload,
		onExpired func(token string)) error
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

	cli, err := app.Auth(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cli, nil
}

func (c *client) messageClient(ctx context.Context) (*messaging.Client, error) {
	app, err := c.initApp(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cli, err := app.Messaging(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cli, nil
}

func (c *client) VerifyToken(ctx context.Context, token string) (UID, error) {
	cli, err := c.authClient(ctx)
	if err != nil {
		return "", errors.WithStack(err)
	}

	decoded, err := cli.VerifyIDToken(ctx, token)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return UID(decoded.UID), nil
}

func (c *client) SendPushNotification(
	ctx context.Context,
	token string,
	title string,
	body string,
	badge int,
	payload PushPayload,
	onExpired func(token string)) error {
	cli, err := c.messageClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	message := &messaging.Message{
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &badge,
				},
				CustomData: payload.IOS(),
			},
		},
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
		Android: &messaging.AndroidConfig{
			Data:         payload,
			Notification: &messaging.AndroidNotification{},
		},
	}

	if _, err := cli.Send(ctx, message); err != nil {
		if messaging.IsRegistrationTokenNotRegistered(err) {
			onExpired(token)
			return nil
		}
		return errors.WithStack(err)
	}

	return nil
}
