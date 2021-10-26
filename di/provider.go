// +build wireinject

package di

import (
	"canvas-server/graph"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"canvas-server/infra/ffmpeg"
	"canvas-server/subscriber"
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
	subscriber.NewSubscriber,
	graph.NewServer,
)

func provideGCSClient() cloud_storage.Client {
	return cloud_storage.NewClient("canvas-329810.appspot.com")
}

func provideDSFactory() datastore.DSFactory {
	return datastore.NewDSFactory("canvas-329810")
}

func ResolveSubscriber() subscriber.Subscriber {
	panic(wire.Build(providerSet))
}

func ResolveGraphQL() graph.Server {
	panic(wire.Build(providerSet))
}
