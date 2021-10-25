package work

import (
	"github.com/pkg/errors"
	"go.mercari.io/datastore/boom"
)

type Repository interface {
	Put(tx *boom.Transaction, item *Entity) error
}

func NewRepository() Repository {
	return &repository{}
}

type repository struct {
}

func (r *repository) Put(tx *boom.Transaction, item *Entity) error {
	if _, err := tx.Put(item); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
