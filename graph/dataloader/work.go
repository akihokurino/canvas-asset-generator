package dataloader

import (
	"canvas-server/infra/datastore/work"
	"context"

	"github.com/pkg/errors"

	"github.com/graph-gophers/dataloader"
)

type WorkLoader struct {
	loader *dataloader.Loader
}

func (l *WorkLoader) Load(ctx context.Context, id string) (*work.Entity, error) {
	thunk := l.loader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result.(*work.Entity), nil
}

func NewWorkLoader(workRepo work.Repository) *WorkLoader {
	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := make([]string, 0, len(keys))
		for _, key := range keys {
			ids = append(ids, key.String())
		}

		workEntities, _ := workRepo.GetMulti(ctx, ids)

		results := make([]*dataloader.Result, len(keys))

		for i, id := range ids {
			results[i] = &dataloader.Result{Data: nil, Error: nil}

			for _, e := range workEntities {
				if e.ID == id {
					results[i].Data = e
					continue
				}
			}

			if results[i].Data == nil {
				results[i].Error = errors.New("entity not found")
			}
		}

		return results
	}

	return &WorkLoader{
		loader: dataloader.NewBatchedLoader(batchFn),
	}
}
