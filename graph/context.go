package graph

import (
	"canvas-server/infra/firebase"
	"context"
)

const authUidStoreKey = "__auth_uid_store_key__"

type ContextProvider interface {
	WithAuthUID(ctx context.Context, uid firebase.UID) (context.Context, error)
	MustAuthUID(ctx context.Context) firebase.UID
}

type contextProvider struct {
}

func NewContextProvider() ContextProvider {
	return &contextProvider{}
}

func (u *contextProvider) WithAuthUID(ctx context.Context, uid firebase.UID) (context.Context, error) {
	return context.WithValue(ctx, authUidStoreKey, uid), nil
}

func (u *contextProvider) MustAuthUID(ctx context.Context) firebase.UID {
	uid, ok := ctx.Value(authUidStoreKey).(firebase.UID)
	if !ok {
		panic("not found access token")
	}
	return uid
}
