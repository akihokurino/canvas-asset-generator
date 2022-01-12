package batch

import (
	"bytes"
	"canvas-server/config"
	"canvas-server/infra/cloud_storage"
	"canvas-server/infra/datastore/thumbnail"
	"canvas-server/infra/datastore/work"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
)

type Server func(mux *http.ServeMux)

func NewServer(
	gcsClient cloud_storage.Client,
	workRepo work.Repository,
	thumbnailRepo thumbnail.Repository) Server {
	export := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		workEntities, err := workRepo.GetAll(ctx)
		if err != nil {
			log.Printf("ExportError: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		workRecords := [][]string{
			{"ID", "VideoPath"},
		}
		for _, entity := range workEntities {
			workRecords = append(workRecords, []string{
				entity.ID, entity.VideoPath,
			})
		}

		var workCSVBuf bytes.Buffer
		workCSVWriter := csv.NewWriter(&workCSVBuf)
		for _, record := range workRecords {
			if err := workCSVWriter.Write(record); err != nil {
				log.Printf("ExportError: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		workCSVWriter.Flush()

		_, err = gcsClient.Save(
			ctx,
			config.CSVBucketName,
			"work.csv",
			workCSVBuf.Bytes(),
			"text/csv")
		if err != nil {
			log.Printf("ExportError: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		thumbnailEntities, err := thumbnailRepo.GetAll(ctx)
		if err != nil {
			log.Printf("ExportError: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		thumbnailRecords := [][]string{
			{"ID", "WorkID", "ImagePath", "Order"},
		}
		for _, entity := range thumbnailEntities {
			thumbnailRecords = append(thumbnailRecords, []string{
				entity.ID, entity.WorkID, entity.ImagePath, fmt.Sprintf("%d", entity.Order),
			})
		}

		var thumbnailCSVBuf bytes.Buffer
		thumbnailCSVWriter := csv.NewWriter(&thumbnailCSVBuf)
		for _, record := range thumbnailRecords {
			if err := thumbnailCSVWriter.Write(record); err != nil {
				log.Printf("ExportError: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		thumbnailCSVWriter.Flush()

		_, err = gcsClient.Save(
			ctx,
			config.CSVBucketName,
			"thumbnail.csv",
			thumbnailCSVBuf.Bytes(),
			"text/csv")
		if err != nil {
			log.Printf("ExportError: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	noAuth := func(server http.Handler) http.Handler {
		return applyMiddleware(server)
	}

	return func(mux *http.ServeMux) {
		mux.Handle("/export", noAuth(http.HandlerFunc(export)))
	}
}

func applyMiddleware(target http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	h := target
	for _, mw := range handlers {
		h = mw(h)
	}
	return h
}
