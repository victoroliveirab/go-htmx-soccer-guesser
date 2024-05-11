package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	user := models.GetLoggingInUser(infra.Db, username, password)

	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := lib.NewSession(user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionCookie := http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		MaxAge:   int(time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &sessionCookie)

	redirectUrl := "/users/" + strconv.FormatInt(int64(user.Id), 10)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
