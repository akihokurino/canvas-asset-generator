package work

import (
	"canvas-server/infra/datastore"
	"context"

	w "go.mercari.io/datastore"

	"github.com/pkg/errors"
	"go.mercari.io/datastore/boom"
)

type Repository interface {
	GetWithPager(ctx context.Context, pager *datastore.Pager) ([]*Entity, error)
	GetMulti(ctx context.Context, ids []string) ([]*Entity, error)
	Get(ctx context.Context, id string) (*Entity, error)
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

func (r *repository) GetWithPager(ctx context.Context, pager *datastore.Pager) ([]*Entity, error) {
	b := boom.FromClient(ctx, r.df(ctx))
	q := b.Client.NewQuery(kind).
		Offset(pager.Offset()).
		Limit(pager.Limit()).
		Order("-CreatedAt")

	var entities []*Entity
	if _, err := b.GetAll(q, &entities); err != nil {
		return nil, errors.WithStack(err)
	}

	return entities, nil
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

func (r *repository) Get(ctx context.Context, id string) (*Entity, error) {
	entity := &Entity{
		ID: id,
	}

	b := boom.FromClient(ctx, r.df(ctx))

	if err := b.Get(entity); err != nil {
		if err == w.ErrNoSuchEntity {
			return nil, errors.WithStack(errors.New("entity not found"))
		}
		return nil, errors.WithStack(err)
	}

	return entity, nil
}

func (r *repository) Put(tx *boom.Transaction, item *Entity) error {
	if _, err := tx.Put(item); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
