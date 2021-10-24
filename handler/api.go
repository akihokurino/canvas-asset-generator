package handler

import (
	"fmt"
	"net/http"
)

type API func(mux *http.ServeMux)

func NewAPI() API {
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello, World!")
	}

	return func(mux *http.ServeMux) {
		mux.HandleFunc("/", apiHandler)
	}
}
