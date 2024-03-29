package graph

import (
	"canvas-asset-generator/graph/dataloader"
	"canvas-asset-generator/infra/firebase"
	"context"
)

const authUidStoreKey = "__auth_uid_store_key__"
const frameDataLoaderStoreKey = "__frame_dataloader_store_key__"
const workDataLoaderStoreKey = "__work_dataloader_store_key__"

type ContextProvider interface {
	WithAuthUID(ctx context.Context, uid firebase.UID) context.Context
	MustAuthUID(ctx context.Context) firebase.UID

	WithFrameDataloader(ctx context.Context, loader *dataloader.FrameLoader) context.Context
	MustFrameDataloader(ctx context.Context) *dataloader.FrameLoader

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

func (u *contextProvider) WithFrameDataloader(ctx context.Context, loader *dataloader.FrameLoader) context.Context {
	return context.WithValue(ctx, frameDataLoaderStoreKey, loader)
}

func (u *contextProvider) MustFrameDataloader(ctx context.Context) *dataloader.FrameLoader {
	v, ok := ctx.Value(frameDataLoaderStoreKey).(*dataloader.FrameLoader)
	if !ok {
		panic("not found frame dataloader")
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
