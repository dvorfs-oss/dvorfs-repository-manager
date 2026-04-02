package api

import (
	"net/http"
	"os"
	"strings"
)

func WithMiddleware(handler http.Handler, middleware func(http.Handler) http.Handler) http.Handler {
	if middleware == nil {
		return handler
	}

	return middleware(handler)
}

func CORSMiddleware(next http.Handler) http.Handler {
	allowedOrigin := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGIN"))
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		switch {
		case allowedOrigin == "*":
			w.Header().Set("Access-Control-Allow-Origin", "*")
		case origin != "" && origin == allowedOrigin:
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
