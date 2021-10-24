package handler

import (
	"canvas-server/usecase"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Subscriber func(mux *http.ServeMux)

func NewSubscriber(splitVideo usecase.SplitVideo) Subscriber {
	videoSplit := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type Payload struct {
			Path string `json:"path"`
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ReadAll: %v", err)
			http.Error(w, "Internal Error, cannot read body", http.StatusInternalServerError)
			return
		}

		var payload Payload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Unmarshal: %v", err)
			http.Error(w, "Internal Error, cannot parse json body", http.StatusInternalServerError)
			return
		}

		if err := splitVideo(ctx, payload.Path); err != nil {
			log.Printf("SplitVideo: %v", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	return func(mux *http.ServeMux) {
		mux.HandleFunc("/video_split", videoSplit)
	}
}
