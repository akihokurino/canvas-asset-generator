package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"canvas-server/graph/generated"
	"canvas-server/graph/model"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/fcm_token"
	"context"

	"go.mercari.io/datastore/boom"
)

func (r *mutationResolver) RegisterFCMToken(ctx context.Context, input model.RegisterFCMToken) (bool, error) {
	uid := r.contextProvider.MustAuthUID(ctx)

	if err := r.transaction(ctx, func(tx *boom.Transaction) error {
		tokenEntity := fcm_token.NewEntity(uid.String(), input.Device, input.Token)

		if err := r.fcmTokenRepo.Put(tx, tokenEntity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Thumbnails(ctx context.Context, page int, limit int) ([]*model.Thumbnail, error) {
	thumbnails, err := r.thumbnailRepo.GetWithPager(ctx, datastore.NewPager(page, limit))
	if err != nil {
		return nil, err
	}

	resItems := make([]*model.Thumbnail, 0, len(thumbnails))
	for _, t := range thumbnails {
		resItems = append(resItems, &model.Thumbnail{
			ID:        t.ID,
			WorkID:    t.WorkID,
			ImagePath: t.ImagePath,
		})
	}

	return resItems, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
