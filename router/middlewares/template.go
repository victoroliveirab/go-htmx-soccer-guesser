package middlewares

import (
	"net/http"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
)

func WithTemplate(name string, data map[string]interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lib.RenderTemplate(w, r, name, data)
	})
}
