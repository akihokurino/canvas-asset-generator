package graph

import (
	"canvas-server/graph/dataloader"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/fcm_token"
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"canvas-server/infra/firebase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	contextProvider ContextProvider
	fireClient      firebase.Client
	gcsClient       cloud_storage.Client
	transaction     datastore.Transaction
	workRepo        work.Repository
	thumbnailRepo   thumbnail.Repository
	fcmTokenRepo    fcm_token.Repository
	workLoader      *dataloader.WorkLoader
	thumbnailLoader *dataloader.ThumbnailLoader
}

func NewResolver(
	contextProvider ContextProvider,
	fireClient firebase.Client,
	gcsClient cloud_storage.Client,
	transaction datastore.Transaction,
	workRepo work.Repository,
	thumbnailRepo thumbnail.Repository,
	fcmTokenRepo fcm_token.Repository,
	workLoader *dataloader.WorkLoader,
	thumbnailLoader *dataloader.ThumbnailLoader) *Resolver {
	return &Resolver{
		contextProvider: contextProvider,
		fireClient:      fireClient,
		gcsClient:       gcsClient,
		transaction:     transaction,
		workRepo:        workRepo,
		thumbnailRepo:   thumbnailRepo,
		fcmTokenRepo:    fcmTokenRepo,
		workLoader:      workLoader,
		thumbnailLoader: thumbnailLoader,
	}
}
