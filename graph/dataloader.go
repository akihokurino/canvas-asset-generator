package graph

import (
	"canvas-asset-generator/graph/dataloader"
	"canvas-asset-generator/infra/datastore/frame"
	"canvas-asset-generator/infra/datastore/work"
	"net/http"
)

type Dataloader func(base http.Handler) http.Handler

func NewDataloader(
	contextProvider ContextProvider,
	workRepo work.Repository,
	frameRepo frame.Repository) Dataloader {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = contextProvider.WithWorkDataloader(ctx, dataloader.NewWorkLoader(workRepo))
			ctx = contextProvider.WithFrameDataloader(ctx, dataloader.NewFrameLoader(frameRepo))

			base.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
