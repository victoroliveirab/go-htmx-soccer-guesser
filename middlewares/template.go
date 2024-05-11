package middlewares

import (
	"net/http"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
)

func WithTemplate(name string, data interface{}) func(w http.ResponseWriter, r *http.Request) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lib.RenderTemplate(w, name, data)
	})
	return handler.ServeHTTP
}
