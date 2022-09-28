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

// Work is the resolver for the work field.
func (r *frameResolver) Work(ctx context.Context, obj *model.Frame) (*model.Work, error) {
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

// RegisterFCMToken is the resolver for the registerFCMToken field.
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

// Works is the resolver for the works field.
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

// Work is the resolver for the work field.
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

// Frames is the resolver for the frames field.
func (r *queryResolver) Frames(ctx context.Context, page int, limit int) (*model.FrameConnection, error) {
	frameEntities, hasNext, err := r.frameRepo.GetWithPager(ctx, datastore.NewPager(page, limit))
	if err != nil {
		return nil, err
	}

	edges := make([]*model.FrameEdge, 0, len(frameEntities))
	for _, entity := range frameEntities {
		resizedImagePath, _ := url.Parse(entity.ResizedImagePath)
		signedImageURL, _ := r.gcsClient.Signature(resizedImagePath)

		edges = append(edges, &model.FrameEdge{
			Node: &model.Frame{
				ID:          entity.ID,
				WorkID:      entity.WorkID,
				ImageUrl:    signedImageURL.String(),
				ImageGsPath: entity.ImagePath,
			},
		})
	}

	count, err := r.frameRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}

	return &model.FrameConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			TotalCount:  int(count),
			HasNextPage: hasNext,
		},
	}, nil
}

// Frames is the resolver for the frames field.
func (r *workResolver) Frames(ctx context.Context, obj *model.Work, limit *int) ([]*model.Frame, error) {
	frameLoader := r.contextProvider.MustFrameDataloader(ctx)

	frameEntities, err := frameLoader.Load(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	if limit != nil {
		frameEntities = frameEntities[0:*limit]
	}

	resItems := make([]*model.Frame, 0, len(frameEntities))
	for _, entity := range frameEntities {
		resizedImagePath, _ := url.Parse(entity.ResizedImagePath)
		signedImageURL, _ := r.gcsClient.Signature(resizedImagePath)

		resItems = append(resItems, &model.Frame{
			ID:          entity.ID,
			WorkID:      entity.WorkID,
			ImageUrl:    signedImageURL.String(),
			ImageGsPath: entity.ImagePath,
		})
	}

	return resItems, nil
}

// Frame returns generated.FrameResolver implementation.
func (r *Resolver) Frame() generated.FrameResolver { return &frameResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Work returns generated.WorkResolver implementation.
func (r *Resolver) Work() generated.WorkResolver { return &workResolver{r} }

type frameResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type workResolver struct{ *Resolver }
