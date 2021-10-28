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
	thunk := l.loader.LoadMany(ctx, dataloader.NewKeysFromStrings([]string{id}))
	results, errs := thunk()
	if errs != nil && len(errs) > 0 {
		return nil, errors.WithStack(errs[0])
	}

	items := make([]*thumbnail.Entity, 0, len(results))
	for _, result := range results {
		items = append(items, result.(*thumbnail.Entity))
	}

	return items, nil
}

func NewThumbnailLoader(thumbnailRepo thumbnail.Repository) *ThumbnailLoader {
	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := make([]string, 0, len(keys))
		for _, key := range keys {
			ids = append(ids, key.String())
		}

		thumbnailEntities := make([]*thumbnail.Entity, 0)
		for _, id := range ids {
			es, _ := thumbnailRepo.GetAllByWork(ctx, id)
			thumbnailEntities = append(thumbnailEntities, es...)
		}

		results := make([]*dataloader.Result, len(keys))

		for i, entity := range thumbnailEntities {
			results[i] = &dataloader.Result{Data: entity, Error: nil}
		}

		return results
	}

	return &ThumbnailLoader{
		loader: dataloader.NewBatchedLoader(batchFn),
	}
}
