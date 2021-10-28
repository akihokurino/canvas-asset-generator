package work

import (
	"canvas-server/infra/datastore"
	"context"

	w "go.mercari.io/datastore"

	"github.com/pkg/errors"
	"go.mercari.io/datastore/boom"
)

type Repository interface {
	GetMulti(ctx context.Context, ids []string) ([]*Entity, error)
	Put(tx *boom.Transaction, item *Entity) error
}

func NewRepository(df datastore.DSFactory) Repository {
	return &repository{
		df: df,
	}
}

type repository struct {
	df datastore.DSFactory
}

func (r *repository) GetMulti(ctx context.Context, ids []string) ([]*Entity, error) {
	entities := make([]*Entity, 0, len(ids))
	for _, id := range ids {
		entities = append(entities, &Entity{
			ID: id,
		})
	}

	b := boom.FromClient(ctx, r.df(ctx))

	if err := b.GetMulti(entities); err != nil {
		_, ok := err.(w.MultiError)
		if !ok {
			return nil, errors.WithStack(err)
		}
	}

	return entities, nil
}

func (r *repository) Put(tx *boom.Transaction, item *Entity) error {
	if _, err := tx.Put(item); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
