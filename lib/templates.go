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
			"views/index.html",
			"views/teams/index.html",
			"views/base.html",
		),
	)
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return errors.New("Template not found")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "base", data)

	return nil
}
