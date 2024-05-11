package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/controllers/user"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/middlewares"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

func init() {
	lib.RegisterTemplates()
	infra.DbConnect("file:local.db")
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
	mux.Handle("/favicon.ico", http.StripPrefix("/", fileServer))

	// Users

	mux.HandleFunc("GET /signin", middlewares.WithTemplate("signin.html", nil))
	mux.HandleFunc("POST /signin", user.Login)

	mux.HandleFunc("GET /signup", middlewares.WithTemplate("signup.html", nil))

	mux.HandleFunc("GET /users/{id}", user.Index)
	mux.HandleFunc("POST /users", user.Register)

	// Index

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			w.Write([]byte("Resource not found"))
			return
		}

		teams, err := models.GetAllTeams(infra.Db)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		data := struct {
			Env   string
			Teams []models.Team
		}{
			Env:   "DEV",
			Teams: teams,
		}

		lib.RenderTemplate(w, "index.html", data)
	})

	fmt.Println("Listening on port", port)
	muxWithLogging := middlewares.WithLogging(mux)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), muxWithLogging))
}
