package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"canvas-server/graph/generated"
	"canvas-server/graph/model"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/fcm_token"
	"context"
	"net/url"

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

func (r *queryResolver) Works(ctx context.Context, page int, limit int) (*model.WorkConnection, error) {
	workEntities, hasNext, err := r.workRepo.GetWithPager(ctx, datastore.NewPager(page, limit))
	if err != nil {
		return nil, err
	}

	edges := make([]*model.WorkEdge, 0, len(workEntities))
	for _, entity := range workEntities {
		videoPath, _ := url.Parse(entity.VideoPath)
		signedVideoURL, _ := r.gcsClient.Signature(videoPath)
		edges = append(edges, &model.WorkEdge{
			Node: &model.Work{
				ID:          entity.ID,
				VideoUrl:    signedVideoURL.String(),
				VideoGsPath: entity.VideoPath,
			},
		})
	}

	count, err := r.workRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}

	return &model.WorkConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			TotalCount:  int(count),
			HasNextPage: hasNext,
		},
	}, nil
}

func (r *queryResolver) Work(ctx context.Context, id string) (*model.Work, error) {
	workEntity, err := r.workRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	videoPath, _ := url.Parse(workEntity.VideoPath)
	signedVideoURL, _ := r.gcsClient.Signature(videoPath)

	return &model.Work{
		ID:          workEntity.ID,
		VideoUrl:    signedVideoURL.String(),
		VideoGsPath: workEntity.VideoPath,
	}, nil
}

func (r *queryResolver) Thumbnails(ctx context.Context, page int, limit int) (*model.ThumbnailConnection, error) {
	thumbnailEntities, hasNext, err := r.thumbnailRepo.GetWithPager(ctx, datastore.NewPager(page, limit))
	if err != nil {
		return nil, err
	}

	edges := make([]*model.ThumbnailEdge, 0, len(thumbnailEntities))
	for _, entity := range thumbnailEntities {
		imagePath, _ := url.Parse(entity.ImagePath)
		signedImageURL, _ := r.gcsClient.Signature(imagePath)

		edges = append(edges, &model.ThumbnailEdge{
			Node: &model.Thumbnail{
				ID:          entity.ID,
				WorkID:      entity.WorkID,
				ImageUrl:    signedImageURL.String(),
				ImageGsPath: entity.ImagePath,
			},
		})
	}

	count, err := r.thumbnailRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}

	return &model.ThumbnailConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			TotalCount:  int(count),
			HasNextPage: hasNext,
		},
	}, nil
}

func (r *thumbnailResolver) Work(ctx context.Context, obj *model.Thumbnail) (*model.Work, error) {
	workLoader := r.contextProvider.MustWorkDataloader(ctx)

	workEntity, err := workLoader.Load(ctx, obj.WorkID)
	if err != nil {
		return nil, err
	}

	videoPath, _ := url.Parse(workEntity.VideoPath)
	signedVideoURL, _ := r.gcsClient.Signature(videoPath)

	return &model.Work{
		ID:          workEntity.ID,
		VideoUrl:    signedVideoURL.String(),
		VideoGsPath: workEntity.VideoPath,
	}, nil
}

func (r *workResolver) Thumbnails(ctx context.Context, obj *model.Work) ([]*model.Thumbnail, error) {
	thumbnailLoader := r.contextProvider.MustThumbnailDataloader(ctx)

	thumbnailEntities, err := thumbnailLoader.Load(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	resItems := make([]*model.Thumbnail, 0, len(thumbnailEntities))
	for _, entity := range thumbnailEntities {
		imagePath, _ := url.Parse(entity.ImagePath)
		signedImageURL, _ := r.gcsClient.Signature(imagePath)

		resItems = append(resItems, &model.Thumbnail{
			ID:          entity.ID,
			WorkID:      entity.WorkID,
			ImageUrl:    signedImageURL.String(),
			ImageGsPath: entity.ImagePath,
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
