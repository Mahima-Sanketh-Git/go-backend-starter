package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

// chain
func Chain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}

// a function that takes a handler and returns a handler (middleware)
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ==== Bfore ====
		start := time.Now()

		// Hand control to whatever is "next" in the chain
		next.ServeHTTP(w, r)

		// ==== After ===
		log.Printf("method=%s path=%s status=%d duration=%s", r.Method, r.URL.Path, http.StatusOK, time.Since(start))

	})
}

// Auth middleware
func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})

}

// Timeout middleware with context
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})

	}

}
