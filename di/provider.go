//go:build wireinject
// +build wireinject

package di

import (
	"canvas-server/batch"
	"canvas-server/graph"
	"canvas-server/grpc"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/fcm_token"
	"canvas-server/infra/datastore/frame"
	"canvas-server/infra/datastore/work"
	"canvas-server/infra/ffmpeg"
	"canvas-server/infra/firebase"
	"canvas-server/subscriber"
	"os"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	firebase.NewClient,
	provideGCSClient,
	provideDSFactory,
	datastore.NewTransaction,
	work.NewRepository,
	frame.NewRepository,
	fcm_token.NewRepository,
	ffmpeg.NewClient,
	subscriber.NewSplitVideo,
	subscriber.NewServer,
	provideSubscriberAuthenticate,
	batch.NewExportCSV,
	batch.NewResizeFrame,
	batch.NewServer,
	graph.NewResolver,
	graph.NewServer,
	graph.NewContextProvider,
	provideAPIAuthenticate,
	graph.NewDataloader,
	graph.NewCROS,
	grpc.NewAPI,
	grpc.NewServer,
	provideGRPCAuthenticate,
)

func provideGCSClient() cloud_storage.Client {
	return cloud_storage.NewClient(
		"canvas-329810",
		os.Getenv("SERVICE_ACCOUNT_PEM"))
}

func provideDSFactory() datastore.DSFactory {
	return datastore.NewDSFactory("canvas-329810")
}

func provideSubscriberAuthenticate() subscriber.Authenticate {
	return subscriber.NewAuthenticate(os.Getenv("INTERNAL_TOKEN"))
}

func provideAPIAuthenticate(
	contextProvider graph.ContextProvider,
	fireClient firebase.Client) graph.Authenticate {
	return graph.NewAuthenticate(os.Getenv("INTERNAL_TOKEN"), contextProvider, fireClient)
}

func provideGRPCAuthenticate() grpc.Authenticate {
	return grpc.NewAuthenticate(os.Getenv("INTERNAL_TOKEN"))
}

func ResolveSubscriber() subscriber.Server {
	panic(wire.Build(providerSet))
}

func ResolveBatch() batch.Server {
	panic(wire.Build(providerSet))
}

func ResolveGraphQL() graph.Server {
	panic(wire.Build(providerSet))
}

func ResolveGRPC() grpc.Server {
	panic(wire.Build(providerSet))
}
