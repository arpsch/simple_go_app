package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger logs the request duration and its info
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		log.Printf("logger:  %v - %s - %s\n", time.Since(start), r.Method, r.URL.Path)
	})
}
