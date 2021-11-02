package graph

import (
	"canvas-server/graph/dataloader"
	"canvas-server/infra/firebase"
	"context"
)

const authUidStoreKey = "__auth_uid_store_key__"
const thumbnailDataLoaderStoreKey = "__thumbnail_dataloader_store_key__"
const workDataLoaderStoreKey = "__work_dataloader_store_key__"

type ContextProvider interface {
	WithAuthUID(ctx context.Context, uid firebase.UID) context.Context
	MustAuthUID(ctx context.Context) firebase.UID

	WithThumbnailDataloader(ctx context.Context, loader *dataloader.ThumbnailLoader) context.Context
	MustThumbnailDataloader(ctx context.Context) *dataloader.ThumbnailLoader

	WithWorkDataloader(ctx context.Context, loader *dataloader.WorkLoader) context.Context
	MustWorkDataloader(ctx context.Context) *dataloader.WorkLoader
}

type contextProvider struct {
}

func NewContextProvider() ContextProvider {
	return &contextProvider{}
}

func (u *contextProvider) WithAuthUID(ctx context.Context, uid firebase.UID) context.Context {
	return context.WithValue(ctx, authUidStoreKey, uid)
}

func (u *contextProvider) MustAuthUID(ctx context.Context) firebase.UID {
	v, ok := ctx.Value(authUidStoreKey).(firebase.UID)
	if !ok {
		panic("not found uid")
	}
	return v
}

func (u *contextProvider) WithThumbnailDataloader(ctx context.Context, loader *dataloader.ThumbnailLoader) context.Context {
	return context.WithValue(ctx, thumbnailDataLoaderStoreKey, loader)
}

func (u *contextProvider) MustThumbnailDataloader(ctx context.Context) *dataloader.ThumbnailLoader {
	v, ok := ctx.Value(thumbnailDataLoaderStoreKey).(*dataloader.ThumbnailLoader)
	if !ok {
		panic("not found thumbnail dataloader")
	}
	return v
}

func (u *contextProvider) WithWorkDataloader(ctx context.Context, loader *dataloader.WorkLoader) context.Context {
	return context.WithValue(ctx, workDataLoaderStoreKey, loader)
}

func (u *contextProvider) MustWorkDataloader(ctx context.Context) *dataloader.WorkLoader {
	v, ok := ctx.Value(workDataLoaderStoreKey).(*dataloader.WorkLoader)
	if !ok {
		panic("not found work dataloader")
	}
	return v
}
