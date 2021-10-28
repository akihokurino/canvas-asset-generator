// +build wireinject

package di

import (
	"canvas-server/graph"
	"canvas-server/graph/dataloader"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/fcm_token"
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"canvas-server/infra/ffmpeg"
	"canvas-server/infra/firebase"
	"canvas-server/subscriber"
	"canvas-server/usecase"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	firebase.NewClient,
	provideGCSClient,
	provideDSFactory,
	datastore.NewTransaction,
	work.NewRepository,
	thumbnail.NewRepository,
	fcm_token.NewRepository,
	ffmpeg.NewClient,
	usecase.NewSplitVideo,
	subscriber.NewSubscriber,
	dataloader.NewWorkLoader,
	graph.NewResolver,
	graph.NewServer,
	graph.NewContextProvider,
	graph.NewAuthenticate,
	graph.NewCROS,
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
