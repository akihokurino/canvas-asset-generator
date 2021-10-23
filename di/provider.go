// +build wireinject

package di

import (
	"canvas-server/handlers"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	handlers.NewSubscriber,
	handlers.NewAPI,
)

func ResolveSubscriber() handlers.Subscriber {
	panic(wire.Build(providerSet))
}

func ResolveAPI() handlers.API {
	panic(wire.Build(providerSet))
}
