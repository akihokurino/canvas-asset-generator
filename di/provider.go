// +build wireinject

package di

import (
	"canvas-server/handler"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"canvas-server/infra/ffmpeg"
	"canvas-server/usecase"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	provideGCSClient,
	provideDSFactory,
	datastore.NewTransaction,
	work.NewRepository,
	thumbnail.NewRepository,
	ffmpeg.NewClient,
	usecase.NewSplitVideo,
	handler.NewSubscriber,
	handler.NewAPI,
)

func provideGCSClient() cloud_storage.Client {
	return cloud_storage.NewClient("canvas-329810.appspot.com")
}

func provideDSFactory() datastore.DSFactory {
	return datastore.NewDSFactory("canvas-329810")
}

func ResolveSubscriber() handler.Subscriber {
	panic(wire.Build(providerSet))
}

func ResolveAPI() handler.API {
	panic(wire.Build(providerSet))
}
