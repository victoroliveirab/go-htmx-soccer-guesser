package infra

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

var Db *sql.DB
var maxRetries = 5

func DbConnect(url string) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	var journalMode string
	// Enable WAL Mode
	err = db.QueryRow("PRAGMA journal_mode=WAL;").Scan(&journalMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to enable WAL mode: %s", err)
		os.Exit(1)
	}
	if journalMode != "wal" {
		fmt.Fprintf(os.Stderr, "failed to set WAL mode, current mode = %s", journalMode)
		os.Exit(1)
	}

	Db = db
}

func DbExecuteWithRetry(queryFunc func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = queryFunc()
		if err == nil {
			return nil
		}
		if err.Error() != "database is locked" {
			return err
		}
		time.Sleep(time.Duration((i + 1)) * time.Second) // Exponential backoff
	}
	return err
}
