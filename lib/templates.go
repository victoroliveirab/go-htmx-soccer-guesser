package lib

import (
	"errors"
	"html/template"
	"net/http"
)

var templates map[string]*template.Template

func RegisterTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index.html"] = template.Must(
		template.ParseFiles(
			"views/base.html",
			"views/teams/index.html",
			"views/index.html",
		),
	)
	templates["signup.html"] = template.Must(
		template.ParseFiles(
			"views/base.html",
			"views/signup.html",
		),
	)
	templates["user.html"] = template.Must(
		template.ParseFiles(
			"views/base.html",
			"views/user.html",
		),
	)
	templates["signin.html"] = template.Must(
		template.ParseFiles(
			"views/base.html",
			"views/signin.html",
		),
	)
	templates["fixtures/index.html"] = template.Must(
		template.ParseFiles(
			"views/base.html",
			"views/fixtures/index.html",
			// "views/guesses/fixture-modal.html",
			// "views/fixtures/next.html",
		),
	)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return errors.New("Template not found")
	}
	if data == nil {
		data = map[string]interface{}{}
	}

	ctx := r.Context()

	_, exists := data["HideNav"]

	if !exists {
		data["HideNav"] = false
	}
	data["LoggedIn"] = ctx.Value("LoggedIn")
	data["UserID"] = ctx.Value("UserID")
	data["CurrentPath"] = r.URL.Path

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "base", data)

	return nil
}

func RenderPartial(w http.ResponseWriter, path string, block string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles(path)).ExecuteTemplate(w, block, data)
}
