package lib

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

var DbConnection *sql.DB

func DbConnect(url string) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	DbConnection = db
}
