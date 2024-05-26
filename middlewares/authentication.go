package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
)

func WithAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		loggedIn := ctx.Value("LoggedIn").(bool)
		if !loggedIn {
			redirectUrl := r.URL
			http.Redirect(w, r, fmt.Sprintf("/signin?redirect_to=%s", redirectUrl), http.StatusSeeOther)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func WithNoAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		loggedIn := ctx.Value("LoggedIn").(bool)
		if loggedIn {
			userID := int64(ctx.Value("UserID").(int))
			http.Redirect(w, r, "/users/"+strconv.FormatInt(userID, 10), http.StatusSeeOther)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
