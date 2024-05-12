package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/config"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
)

func WithAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(config.SessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		cookie := sessionCookie.Value

		if !lib.IsValidSession(cookie) {
			lib.DeleteSession(cookie)
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		userID := lib.GetUserByCookie(cookie).UserID

		ctx := context.WithValue(r.Context(), "userID", userID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithNoAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(config.SessionCookieName)
		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		cookie := sessionCookie.Value

		if !lib.IsValidSession(cookie) {
			invalidCookie := http.Cookie{
				Name:     config.SessionCookieName,
				Value:    "",
				MaxAge:   -1,
				Path:     config.SessionCookiePath,
				HttpOnly: true,
			}
			http.SetCookie(w, &invalidCookie)

			handler.ServeHTTP(w, r)
			return
		}

		userID := lib.GetUserByCookie(cookie).UserID
		http.Redirect(w, r, "/users/"+strconv.FormatInt(int64(userID), 10), http.StatusSeeOther)
	})
}
