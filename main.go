package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/controllers/fixture"
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

	mux.Handle("GET /signin", middlewares.WithNoAuth(middlewares.WithTemplate("signin.html", map[string]interface{}{"HideNav": true})))
	mux.Handle("POST /signin", middlewares.WithNoAuth(user.Login))

	mux.Handle("GET /signout", middlewares.WithAuth(user.Logout))

	mux.Handle("GET /signup", middlewares.WithNoAuth(middlewares.WithTemplate("signup.html", nil)))

	mux.Handle("GET /users/{id}", middlewares.WithAuth(user.Index))

	mux.Handle("POST /users", middlewares.WithNoAuth(user.Register))

	// Fixtures
	mux.Handle("GET /fixtures/{id}", middlewares.WithNoAuth(fixture.ViewFixture))
	mux.Handle("GET /fixtures", middlewares.WithAuth(fixture.NextFixtures))

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

		// Find a way to default HideNav to false (middleware?)
		data := map[string]interface{}{
			"Env":   "DEV",
			"Teams": teams,
			"Title": "Home",
		}

		lib.RenderTemplate(w, r, "index.html", data)
	})

	session := &lib.Session{
		ID:        "4760a0753a07d3d53217afb028d68901cbf59c64e6913f02130b76612e1308c0",
		UserID:    5,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	lib.AddSession(session)

	fmt.Println("Listening on port", port)
	muxWithSession := middlewares.WithSession(mux)
	muxWithLogging := middlewares.WithLogging(muxWithSession)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), muxWithLogging))
}
