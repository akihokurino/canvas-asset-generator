package thumbnail

import (
	"canvas-server/infra/datastore"
	"context"

	"github.com/pkg/errors"
	"go.mercari.io/datastore/boom"
)

type Repository interface {
	GetAllByWork(ctx context.Context, workID string) ([]*Entity, error)
	Put(tx *boom.Transaction, item *Entity) error
	Delete(tx *boom.Transaction, id string) error
}

func NewRepository(df datastore.DSFactory) Repository {
	return &repository{
		df: df,
	}
}

type repository struct {
	df datastore.DSFactory
}

func (r *repository) GetAllByWork(ctx context.Context, workID string) ([]*Entity, error) {
	b := boom.FromClient(ctx, r.df(ctx))
	q := b.Client.NewQuery(kind).Filter("WorkID =", workID)

	var entities []*Entity
	if _, err := b.GetAll(q, &entities); err != nil {
		return nil, errors.WithStack(err)
	}

	return entities, nil
}

func (r *repository) Put(tx *boom.Transaction, item *Entity) error {
	if _, err := tx.Put(item); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(tx *boom.Transaction, id string) error {
	if err := tx.Delete(&Entity{
		ID: id,
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
