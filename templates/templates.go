package templates

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/config"
)

type AppTemplate struct {
	name string
	tmpl *template.Template
}

var templates map[string]*AppTemplate

var templFuncs template.FuncMap = template.FuncMap{
	"Mod": func(a, b int) int {
		return a % b
	},
}

func LoadTemplate(name string, templatesList ...string) {
	if templates == nil {
		templates = make(map[string]*AppTemplate)
	}
	_, exists := templates[name]
	if exists {
		return
	}

	files := []string{config.TemplatesPath + "/base.html"}
	for _, templName := range templatesList {
		files = append(files, config.TemplatesPath+"/"+templName)
	}
	tmpl := template.New(name)
	tmpl = tmpl.Funcs(templFuncs)
	tmpl, _ = tmpl.ParseFiles(files...)
	templates[name] = &AppTemplate{
		name: name,
		tmpl: tmpl,
	}
}

func GetTemplate(name string) *AppTemplate {
	return templates[name]
}

func (t *AppTemplate) Execute(w http.ResponseWriter, r *http.Request, data map[string]interface{}) error {
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

	fmt.Println(data)
	fmt.Println(t.name)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.tmpl.ExecuteTemplate(w, "base", data)

	return nil
}
