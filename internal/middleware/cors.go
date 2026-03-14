package middleware

import (
	"net/http"
	"strings"
)

// CORS returns a middleware that sets CORS headers.
// allowedOrigins is a comma-separated list of allowed origins.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	origins := make(map[string]bool)
	for _, o := range strings.Split(allowedOrigins, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			origins[trimmed] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestOrigin := r.Header.Get("Origin")

			if origins["*"] || origins[requestOrigin] {
				w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
