package user

import (
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

var RegisterPage http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	tmpl := templates.LoadTemplate("signup", "signup.html")
	data := map[string]interface{}{
		"HideNav": true,
	}
	tmpl.Execute(w, r, data)
})

var RegisterPost http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	id, err := models.CreateUser(infra.Db, username, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectUrl := "/users/" + strconv.FormatInt(id, 10)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)

})
