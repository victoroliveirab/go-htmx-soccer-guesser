package user

import (
	"net/http"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/config"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
)

var Logout http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(config.SessionCookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := sessionCookie.Value

	lib.DeleteSession(cookie)

	invalidCookie := http.Cookie{
		Name:     config.SessionCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     config.SessionCookiePath,
		HttpOnly: true,
	}
	http.SetCookie(w, &invalidCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
})
