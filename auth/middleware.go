package middleware

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENABLE_AUTH") != "true" {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(auth, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(authParts[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || !validateCredentials(pair[0], pair[1]) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateCredentials(username, password string) bool {
	return username == os.Getenv("API_USERNAME") && password == os.Getenv("API_PASSWORD")
}
