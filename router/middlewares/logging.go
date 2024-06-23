package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func WithLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("%s %s - %s\n", r.Method, r.URL, time.Since(timestamp))
	})
}
