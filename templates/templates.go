package templates

import (
	"html/template"
	"log"
	"net/http"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/config"
)

type AppTemplate struct {
	name string
	tmpl *template.Template
}

var templates map[string]*AppTemplate
var partials map[string]*AppTemplate

var templFuncs template.FuncMap = template.FuncMap{
	"Mod": func(a, b int) int {
		return a % b
	},
}

func LoadPartial(name string, templatesList ...string) *AppTemplate {
	if partials == nil {
		partials = make(map[string]*AppTemplate)
	}
	_, exists := partials[name]
	if !exists {
		files := []string{config.TemplatesPath + "/base.html"}
		for _, templName := range templatesList {
			files = append(files, config.TemplatesPath+"/"+templName)
		}
		tmpl := template.New(name)
		tmpl = tmpl.Funcs(templFuncs)
		tmpl, _ = tmpl.ParseFiles(files...)
		partials[name] = &AppTemplate{
			name: name,
			tmpl: tmpl,
		}
		log.Printf("loaded %s partial for the first time\n", name)
		log.Println("templates list:")
		log.Println(templatesList)
	}

	return partials[name]
}

func LoadTemplate(name string, templatesList ...string) *AppTemplate {
	if templates == nil {
		templates = make(map[string]*AppTemplate)
	}
	_, exists := templates[name]
	if !exists {
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
		log.Printf("loaded %s template for the first time\n", name)
		log.Println("templates list:")
		log.Println(templatesList)
	}

	return GetTemplate(name)
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

	log.Printf("executing template %s\n", t.name)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.tmpl.ExecuteTemplate(w, "base", data)

	return nil
}

func (t *AppTemplate) ExecutePartial(w http.ResponseWriter, r *http.Request, block string, data interface{}) error {
	log.Printf("executing partial %s\n", t.name)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.tmpl.ExecuteTemplate(w, block, data)

	return nil
}
