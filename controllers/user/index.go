package user

import (
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

var Index http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	groups := models.GetGroupsAssociatedWithUserId(infra.Db, id)

	data := map[string]interface{}{
		"User":       user,
		"UserGroups": groups,
	}

	tmpl := templates.LoadTemplate("show-user", "user/index.html")
	tmpl.Execute(w, r, data)
})
