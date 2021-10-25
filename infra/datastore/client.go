package datastore

import (
	"context"

	"github.com/pkg/errors"

	"go.mercari.io/datastore/boom"

	"cloud.google.com/go/datastore"
	"go.mercari.io/datastore/clouddatastore"

	w "go.mercari.io/datastore"
)

type DSFactory func(ctx context.Context) w.Client

func NewDSFactory(projectID string) DSFactory {
	return func(ctx context.Context) w.Client {
		dc, err := datastore.NewClient(ctx, projectID)
		if err != nil {
			panic(err)
		}

		client, err := clouddatastore.FromClient(ctx, dc)
		if err != nil {
			panic(err)
		}

		return client
	}
}

type Transaction func(ctx context.Context, fn func(tx *boom.Transaction) error) error

func NewTransaction(df DSFactory) Transaction {
	return func(ctx context.Context, fn func(tx *boom.Transaction) error) error {
		b := boom.FromClient(ctx, df(ctx))
		if _, err := b.RunInTransaction(fn); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}
