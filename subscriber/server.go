package subscriber

import (
	"net/http"
)

type Server func(mux *http.ServeMux)

func NewServer(splitVideo SplitVideo, authenticate Authenticate) Server {
	auth := func(server http.Handler) http.Handler {
		return applyMiddleware(
			server,
			authenticate)
	}

	return func(mux *http.ServeMux) {
		mux.Handle("/split-video", auth(http.HandlerFunc(splitVideo)))
	}
}

func applyMiddleware(target http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	h := target
	for _, mw := range handlers {
		h = mw(h)
	}
	return h
}
