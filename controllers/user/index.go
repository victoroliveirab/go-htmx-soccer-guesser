package user

import (
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := models.GetUserById(infra.Db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	lib.RenderTemplate(w, "user.html", user)

}
