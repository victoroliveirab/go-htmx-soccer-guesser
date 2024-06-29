package lib

import (
	"errors"
	"html/template"
	"net/http"
)

var templates map[string]*template.Template

var funcMap = template.FuncMap{
	"Mod": func(a, b int) int {
		return a % b
	},
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

func RenderPartial(w http.ResponseWriter, paths []string, block string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles(paths...)).ExecuteTemplate(w, block, data)
}
