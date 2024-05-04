package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w http.ResponseWriter, name string, data interface{}) {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err == nil {
		return
	}
	fmt.Println(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index.html"] = template.Must(template.ParseFiles("views/index.html", "components/env-card.html", "views/base.html"))
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return errors.New("Template not found")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "base", data)

	return nil
}

func main() {
	// env := os.Getenv("APP_ENV")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/favicon.ico", http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Env string
		}{
			Env: "DEV",
		}

		renderTemplate(w, "index.html", data)
	})

	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
