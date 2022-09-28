package batch

import (
	"bytes"
	"canvas-server/config"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore"
	"canvas-server/infra/datastore/frame"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"go.mercari.io/datastore/boom"

	"golang.org/x/image/draw"
)

type ResizeFrame func(w http.ResponseWriter, r *http.Request)

func NewResizeFrame(
	gcsClient cloud_storage.Client,
	tx datastore.Transaction,
	frameRepo frame.Repository) ResizeFrame {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		frameEntities, err := frameRepo.GetAllByNotResized(ctx)
		if err != nil {
			log.Printf("ResizeFrame: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, frameEntity := range frameEntities {
			data, err := gcsClient.Download(ctx, config.FrameBucketName, fmt.Sprintf("%s/%d", frameEntity.WorkID, frameEntity.Order))
			if err != nil {
				log.Printf("ResizeFrame: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			imgSource, _, err := image.Decode(data)
			if err != nil {
				log.Printf("ResizeFrame: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			rect := imgSource.Bounds()
			resizedImage := image.NewRGBA(image.Rect(0, 0, rect.Dx()/5, rect.Dy()/5))
			draw.CatmullRom.Scale(resizedImage, resizedImage.Bounds(), imgSource, rect, draw.Over, nil)

			resizeBuf := bytes.NewBuffer(nil)
			if err := jpeg.Encode(resizeBuf, resizedImage, &jpeg.Options{Quality: 80}); err != nil {
				log.Printf("ResizeFrame: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			resizeURL, err := gcsClient.Save(
				ctx,
				config.FrameBucketName,
				fmt.Sprintf("%s/%d/resized", frameEntity.WorkID, frameEntity.Order),
				resizeBuf.Bytes(),
				"image/jpeg")
			if err != nil {
				log.Printf("ResizeFrame: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			frameEntity.ResizedImagePath = resizeURL.String()

			if err := tx(ctx, func(tx *boom.Transaction) error {
				return frameRepo.Put(tx, frameEntity)
			}); err != nil {
				log.Printf("ResizeFrame: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
