package dataloader

import (
	"canvas-server/infra/datastore/frame"
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
)

type FrameLoader struct {
	loader *dataloader.Loader
}

func (l *FrameLoader) Load(ctx context.Context, id string) ([]*frame.Entity, error) {
	thunk := l.loader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if es, ok := result.([]*frame.Entity); ok {
		return es, nil
	}

	return make([]*frame.Entity, 0), nil
}

func NewFrameLoader(frameRepo frame.Repository) *FrameLoader {
	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := make([]string, 0, len(keys))
		for _, key := range keys {
			ids = append(ids, key.String())
		}

		frameEntityMap := make(map[string][]*frame.Entity, 0)
		for _, id := range ids {
			es, _ := frameRepo.GetAllByWork(ctx, id)
			frameEntityMap[id] = es
		}

		results := make([]*dataloader.Result, len(keys))

		for i, id := range ids {
			results[i] = &dataloader.Result{Data: frameEntityMap[id], Error: nil}
		}

		return results
	}

	return &FrameLoader{
		loader: dataloader.NewBatchedLoader(batchFn),
	}
}
