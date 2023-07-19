//go:build wireinject
// +build wireinject

package di

import (
	"canvas-asset-generator/batch"
	"canvas-asset-generator/graph"
	"canvas-asset-generator/grpc"
	"canvas-asset-generator/infra/cloud_storage"
	"canvas-asset-generator/infra/datastore"
	"canvas-asset-generator/infra/datastore/fcm_token"
	"canvas-asset-generator/infra/datastore/frame"
	"canvas-asset-generator/infra/datastore/work"
	"canvas-asset-generator/infra/ffmpeg"
	"canvas-asset-generator/infra/firebase"
	"canvas-asset-generator/subscriber"
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
