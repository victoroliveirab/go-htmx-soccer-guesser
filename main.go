package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
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

	mux.HandleFunc("GET /signin", func(w http.ResponseWriter, r *http.Request) {
		lib.RenderTemplate(w, "signin.html", nil)
	})

	mux.HandleFunc("POST /signin", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		user := models.GetLoggingInUser(infra.Db, username, password)

		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := lib.NewSession(user.Id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionCookie := http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			MaxAge:   int(time.Hour),
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, &sessionCookie)

		redirectUrl := "/users/" + strconv.FormatInt(int64(user.Id), 10)
		http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	})

	mux.HandleFunc("GET /signup", func(w http.ResponseWriter, r *http.Request) {
		lib.RenderTemplate(w, "signup.html", nil)
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
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
		lib.RenderTemplate(w, "user.html", user)
	})

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
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
