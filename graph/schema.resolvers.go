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

func (r *queryResolver) Works(ctx context.Context, page int, limit int) ([]*model.Work, error) {
	workEntities, err := r.workRepo.GetWithPager(ctx, datastore.NewPager(page, limit))
	if err != nil {
		return nil, err
	}

	resItems := make([]*model.Work, 0, len(workEntities))
	for _, entity := range workEntities {
		resItems = append(resItems, &model.Work{
			ID:        entity.ID,
			VideoPath: entity.VideoPath,
		})
	}

	return resItems, nil
}

func (r *queryResolver) Work(ctx context.Context, id string) (*model.Work, error) {
	workEntity, err := r.workRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Work{
		ID:        workEntity.ID,
		VideoPath: workEntity.VideoPath,
	}, nil
}

func (r *queryResolver) Thumbnails(ctx context.Context, page int, limit int) ([]*model.Thumbnail, error) {
	thumbnailEntities, err := r.thumbnailRepo.GetWithPager(ctx, datastore.NewPager(page, limit))
	if err != nil {
		return nil, err
	}

	resItems := make([]*model.Thumbnail, 0, len(thumbnailEntities))
	for _, entity := range thumbnailEntities {
		resItems = append(resItems, &model.Thumbnail{
			ID:        entity.ID,
			WorkID:    entity.WorkID,
			ImagePath: entity.ImagePath,
		})
	}

	return resItems, nil
}

func (r *thumbnailResolver) Work(ctx context.Context, obj *model.Thumbnail) (*model.Work, error) {
	workEntity, err := r.workLoader.Load(ctx, obj.WorkID)
	if err != nil {
		return nil, err
	}

	return &model.Work{
		ID:        workEntity.ID,
		VideoPath: workEntity.VideoPath,
	}, nil
}

func (r *workResolver) Thumbnails(ctx context.Context, obj *model.Work) ([]*model.Thumbnail, error) {
	thumbnailEntities, err := r.thumbnailLoader.Load(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	resItems := make([]*model.Thumbnail, 0, len(thumbnailEntities))
	for _, entity := range thumbnailEntities {
		resItems = append(resItems, &model.Thumbnail{
			ID:        entity.ID,
			WorkID:    entity.WorkID,
			ImagePath: entity.ImagePath,
		})
	}

	return resItems, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Thumbnail returns generated.ThumbnailResolver implementation.
func (r *Resolver) Thumbnail() generated.ThumbnailResolver { return &thumbnailResolver{r} }

// Work returns generated.WorkResolver implementation.
func (r *Resolver) Work() generated.WorkResolver { return &workResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type thumbnailResolver struct{ *Resolver }
type workResolver struct{ *Resolver }
