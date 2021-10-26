package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"canvas-server/graph/generated"
	"canvas-server/graph/model"
	"context"
)

func (r *mutationResolver) RegisterFCMToken(ctx context.Context, input model.RegisterFCMToken) (bool, error) {
	return true, nil
}

func (r *queryResolver) Thumbnails(ctx context.Context) ([]*model.Thumbnail, error) {
	return []*model.Thumbnail{
		{
			ID:        "1",
			WorkID:    "1",
			ImagePath: "test",
		},
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
