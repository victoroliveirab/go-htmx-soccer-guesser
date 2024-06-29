package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/router"
)

func init() {
	infra.DbConnect("file:local.db")
}

func main() {
	// env := os.Getenv("APP_ENV")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	httpRouter := router.New()

	session := &lib.Session{
		ID:        "0c0ce7e041d521061f25f49f692cf0f6171543a284c35e8b03760a05b262141d",
		UserID:    5,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	lib.AddSession(session)

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), httpRouter))
}
