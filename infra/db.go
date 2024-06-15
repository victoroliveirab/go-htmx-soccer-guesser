package infra

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

var Db *sql.DB
var maxRetries = 5

func DbConnect(url string) error {
	db, err := sql.Open("libsql", url)
	if err != nil {
		return fmt.Errorf("failed to open db %s: %w", url, err)
	}

	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode=WAL;").Scan(&journalMode)
	if err != nil {
		return fmt.Errorf("failed to enable WAL mode: %w", err)
	}
	if journalMode != "wal" {
		return fmt.Errorf("failed to set WAL mode, current mode = %s", journalMode)
	}

	Db = db
	return nil
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
