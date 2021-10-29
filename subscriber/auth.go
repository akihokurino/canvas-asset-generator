package subscriber

import (
	"net/http"
)

const (
	authKey = "Authorization"
)

type Authenticate func(base http.Handler) http.Handler

func NewAuthenticate(privateKey string) Authenticate {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(authKey)
			if token != privateKey {
				http.Error(w, "unauthorized", 401)
				return
			}

			base.ServeHTTP(w, r)
		})
	}
}
