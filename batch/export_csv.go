package batch

import (
	"bytes"
	"canvas-asset-generator/config"
	"canvas-asset-generator/infra/cloud_storage"
	"canvas-asset-generator/infra/datastore/frame"
	"canvas-asset-generator/infra/datastore/work"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
)

type ExportCSV func(w http.ResponseWriter, r *http.Request)

func NewExportCSV(
	gcsClient cloud_storage.Client,
	workRepo work.Repository,
	frameRepo frame.Repository) ExportCSV {
	return func(w http.ResponseWriter, r *http.Request) {
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

		frameEntities, err := frameRepo.GetAll(ctx)
		if err != nil {
			log.Printf("ExportError: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		frameRecords := [][]string{
			{"ID", "WorkID", "ImagePath", "Order"},
		}
		for _, entity := range frameEntities {
			frameRecords = append(frameRecords, []string{
				entity.ID, entity.WorkID, entity.ImagePath, fmt.Sprintf("%d", entity.Order),
			})
		}

		var frameCSVBuf bytes.Buffer
		frameCSVWriter := csv.NewWriter(&frameCSVBuf)
		for _, record := range frameRecords {
			if err := frameCSVWriter.Write(record); err != nil {
				log.Printf("ExportError: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		frameCSVWriter.Flush()

		_, err = gcsClient.Save(
			ctx,
			config.CSVBucketName,
			"frame.csv",
			frameCSVBuf.Bytes(),
			"text/csv")
		if err != nil {
			log.Printf("ExportError: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
