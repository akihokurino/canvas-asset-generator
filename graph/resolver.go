package graph

import (
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	workRepo      work.Repository
	thumbnailRepo thumbnail.Repository
}

func NewResolver(
	workRepo work.Repository,
	thumbnailRepo thumbnail.Repository) *Resolver {
	return &Resolver{
		workRepo:      workRepo,
		thumbnailRepo: thumbnailRepo,
	}
}
