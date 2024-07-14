package router

import (
	"net/http"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/controllers/fixture"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/controllers/group"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/controllers/guess"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/controllers/league"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/controllers/user"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/router/middlewares"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

func New() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/favicon.ico", http.StripPrefix("/", fileServer))

	// Users

	mux.Handle("GET /signin", middlewares.WithNoAuth(user.LoginPage))
	mux.Handle("POST /signin", middlewares.WithNoAuth(user.LoginPost))

	mux.Handle("GET /signout", middlewares.WithAuth(user.Logout))

	templates.LoadTemplate("signup", "signup.html")
	mux.Handle("GET /signup", middlewares.WithNoAuth(user.RegisterPage))

	mux.Handle("GET /users/{id}", middlewares.WithAuth(user.Index))

	mux.Handle("POST /users", middlewares.WithNoAuth(user.RegisterPost))

	// Groups
	mux.Handle("GET /groups/{id}", middlewares.WithAuth(group.Show))

	// Fixtures
	mux.Handle("GET /fixtures/{id}", middlewares.WithAuth(fixture.ViewFixture))
	templates.LoadTemplate("fixtures", "fixtures/index.html")
	mux.Handle("GET /fixtures", middlewares.WithAuth(fixture.FixturesByDate))

	// Guesses
	mux.Handle("POST /guesses", middlewares.WithAuth(guess.Create))
	mux.Handle("GET /guesses/group/{groupId}/user/{userId}", middlewares.WithAuth(guess.GetUserGuessesByGroup))

	// Leagues
	mux.Handle("GET /leagues/{id}", middlewares.WithAuth(league.ViewLeagueWithStandings))

	// Index

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := templates.LoadTemplate("index", "index.html")
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			w.Write([]byte("Resource not found"))
			return
		}

		// Find a way to default HideNav to false (middleware?)
		data := map[string]interface{}{
			"Env":   "DEV",
			"Title": "Home",
		}

		tmpl.Execute(w, r, data)
	})

	session := &lib.Session{
		ID:        "593f4fd9c72267eb31e4060859b75102cea329c63a75107754e58f40ef76b1c4",
		UserID:    5,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	lib.AddSession(session)

	muxWithSession := middlewares.WithSession(mux)
	muxWithLogging := middlewares.WithLogging(muxWithSession)
	return muxWithLogging
}
