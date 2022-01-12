package subscriber

import (
	"canvas-server/usecase"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Server func(mux *http.ServeMux)

func NewServer(splitVideo usecase.SplitVideo, authenticate Authenticate) Server {
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
			http.Error(w, "Internal Error, cannot split video", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	auth := func(server http.Handler) http.Handler {
		return applyMiddleware(
			server,
			authenticate)
	}

	return func(mux *http.ServeMux) {
		mux.Handle("/video_split", auth(http.HandlerFunc(videoSplit)))
	}
}

func applyMiddleware(target http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	h := target
	for _, mw := range handlers {
		h = mw(h)
	}
	return h
}
