package handler

import (
	"crypto/sha256"
	"net/http"
)

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sha256.Sum256([]byte(username))
		sha256.Sum256([]byte(password))
	})
}

func AdminHandler(w http.ResponseWriter, req *http.Request) {
}
