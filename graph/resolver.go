package graph

import (
	"canvas-asset-generator/infra/cloud_storage"
	"canvas-asset-generator/infra/datastore"
	"canvas-asset-generator/infra/datastore/fcm_token"
	"canvas-asset-generator/infra/datastore/frame"
	"canvas-asset-generator/infra/datastore/work"
	"canvas-asset-generator/infra/firebase"
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
	frameRepo       frame.Repository
	fcmTokenRepo    fcm_token.Repository
}

func NewResolver(
	contextProvider ContextProvider,
	fireClient firebase.Client,
	gcsClient cloud_storage.Client,
	transaction datastore.Transaction,
	workRepo work.Repository,
	frameRepo frame.Repository,
	fcmTokenRepo fcm_token.Repository) *Resolver {
	return &Resolver{
		contextProvider: contextProvider,
		fireClient:      fireClient,
		gcsClient:       gcsClient,
		transaction:     transaction,
		workRepo:        workRepo,
		frameRepo:       frameRepo,
		fcmTokenRepo:    fcmTokenRepo,
	}
}
