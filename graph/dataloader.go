package graph

import (
	"canvas-server/graph/dataloader"
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"net/http"
)

type Dataloader func(base http.Handler) http.Handler

func NewDataloader(
	contextProvider ContextProvider,
	workRepo work.Repository,
	thumbnailRepo thumbnail.Repository) Dataloader {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = contextProvider.WithWorkDataloader(ctx, dataloader.NewWorkLoader(workRepo))
			ctx = contextProvider.WithThumbnailDataloader(ctx, dataloader.NewThumbnailLoader(thumbnailRepo))

			base.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
