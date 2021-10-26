package graph

import (
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"canvas-server/infra/firebase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	fireClient    firebase.Client
	workRepo      work.Repository
	thumbnailRepo thumbnail.Repository
}

func NewResolver(
	fireClient firebase.Client,
	workRepo work.Repository,
	thumbnailRepo thumbnail.Repository) *Resolver {
	return &Resolver{
		fireClient:    fireClient,
		workRepo:      workRepo,
		thumbnailRepo: thumbnailRepo,
	}
}
