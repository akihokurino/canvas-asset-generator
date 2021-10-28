package dataloader

import (
	"canvas-server/infra/datastore/thumbnail"
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
)

type ThumbnailLoader struct {
	loader *dataloader.Loader
}

func (l *ThumbnailLoader) Load(ctx context.Context, id string) ([]*thumbnail.Entity, error) {
	thunk := l.loader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if es, ok := result.([]*thumbnail.Entity); ok {
		return es, nil
	}

	return make([]*thumbnail.Entity, 0), nil
}

func NewThumbnailLoader(thumbnailRepo thumbnail.Repository) *ThumbnailLoader {
	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := make([]string, 0, len(keys))
		for _, key := range keys {
			ids = append(ids, key.String())
		}

		thumbnailEntityMap := make(map[string][]*thumbnail.Entity, 0)
		for _, id := range ids {
			es, _ := thumbnailRepo.GetAllByWork(ctx, id)
			thumbnailEntityMap[id] = es
		}

		results := make([]*dataloader.Result, len(keys))

		for i, id := range ids {
			results[i] = &dataloader.Result{Data: thumbnailEntityMap[id], Error: nil}
		}

		return results
	}

	return &ThumbnailLoader{
		loader: dataloader.NewBatchedLoader(batchFn),
	}
}
