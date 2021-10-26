package graph

import "net/http"

type CROS func(base http.Handler) http.Handler

func NewCROS() CROS {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, X-User-Id, Authorization, X-Requested-With, X-Requested-By")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
			if r.Method == "OPTIONS" {
				w.WriteHeader(200)
				return
			}
			base.ServeHTTP(w, r)
		})
	}
}
