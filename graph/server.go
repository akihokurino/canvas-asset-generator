package graph

import (
	"canvas-server/graph/generated"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type Server func(mux *http.ServeMux)

func NewServer(resolver *Resolver) Server {
	return func(mux *http.ServeMux) {
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", srv)
	}
}
