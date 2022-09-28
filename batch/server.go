package batch

import (
	"net/http"
)

type Server func(mux *http.ServeMux)

func NewServer(exportCSV ExportCSV, resizeFrame ResizeFrame) Server {

	noAuth := func(server http.Handler) http.Handler {
		return applyMiddleware(server)
	}

	return func(mux *http.ServeMux) {
		mux.Handle("/export", noAuth(http.HandlerFunc(exportCSV)))
		mux.Handle("/resize", noAuth(http.HandlerFunc(resizeFrame)))
	}
}

func applyMiddleware(target http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	h := target
	for _, mw := range handlers {
		h = mw(h)
	}
	return h
}
