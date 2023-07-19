package graph

import (
	"canvas-asset-generator/graph/generated"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type Server func(mux *http.ServeMux)

func NewServer(resolver *Resolver, authenticate Authenticate, cros CROS, dataloader Dataloader) Server {
	mw := func(server *handler.Server) http.Handler {
		return applyMiddleware(
			server,
			dataloader,
			authenticate,
			cros)
	}

	return func(mux *http.ServeMux) {
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", mw(srv))
	}
}

func applyMiddleware(target http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	h := target
	for _, mw := range handlers {
		h = mw(h)
	}
	return h
}
