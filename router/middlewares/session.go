package middlewares

import (
	"context"
	"net/http"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/config"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
)

// Session middleware.
func WithSession(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "LoggedIn", false)
		ctx = context.WithValue(ctx, "UserID", -1)

		sessionCookie, err := r.Cookie(config.SessionCookieName)
		if err != nil {
			handler.ServeHTTP(w, r.WithContext(ctx))
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
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		userID := lib.GetUserByCookie(sessionCookie.Value).UserID
		ctx = context.WithValue(ctx, "LoggedIn", true)
		ctx = context.WithValue(ctx, "UserID", userID)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
