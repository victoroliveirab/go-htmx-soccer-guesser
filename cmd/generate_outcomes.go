package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

type PointsMap map[string]int

func init() {
	err := infra.DbConnect("file:local.db?_busy_timeout=5000")
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		panic(err)
	}
}

func main() {
	db := infra.Db
	defer db.Close()

	var groupPointsMap map[int]PointsMap
	groupPointsMap = make(map[int]PointsMap, 0)

	rows, err := db.Query(`
		SELECT guess.id, guess.group_id, guess.home_goals, guess.away_goals, fixture.home_score, fixture.away_score
		FROM Guesses guess
		JOIN Fixtures fixture ON fixture.id = guess.fixture_id
		WHERE guess.locked = 1
		AND fixture.status IN (0, 2)
	`)

	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		panic(err)
	}

	defer rows.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
		panic(err)
	}

	log.Printf("starting updating guesses")

	i := 0

	for rows.Next() {
		var guess models.Guess
		var fixture models.Fixture
		err := rows.Scan(&guess.Id, &guess.GroupId, &guess.HomeGoals, &guess.AwayGoals, &fixture.HomeScore, &fixture.AwayScore)
		if err != nil {
			log.Fatalf("failed to scan row index %d: %v", i, err)
			err = tx.Rollback()
			if err != nil {
				log.Fatalf("failed rolling back transaction: %v", err)
			}
			panic(err)
		}

		pointsMap, exists := groupPointsMap[guess.GroupId]

		if !exists {
			pointsMap, err = getPointsMap(&guess)
			if err != nil {
				log.Fatalf("failed to fetch points table for group_id %d: %v", guess.GroupId, err)
				err = tx.Rollback()
				if err != nil {
					log.Fatalf("failed rolling back transaction: %v", err)
				}
				panic(err)
			}
		}

		outcome := lib.DefineOutcome(&guess, &fixture)
		points := pointsMap[outcome.String()]

		_, err = tx.Exec("UPDATE guesses SET outcome = ?, points = ? WHERE id = ?", outcome, points, guess.Id)
		if err != nil {
			log.Fatalf("failed adding update statement to transaction for guess_id %d: %v", guess.Id, err)
			err = tx.Rollback()
			if err != nil {
				log.Fatalf("failed rolling back transaction: %v", err)
			}
			panic(err)
		}

		log.Printf("executed update for guess = %d", guess.Id)
		i++
	}

	if i == 0 {
		log.Printf("nothing to update")
		return
	}

	err = tx.Commit()

	if err != nil {
		log.Fatalf("failed commiting transaction: %v", err)
		err = tx.Rollback()
		if err != nil {
			log.Fatalf("failed rolling back transaction: %v", err)
		}
		panic(err)
	}

	log.Printf("finished updating %d guesses", i+1)
}

func getPointsMap(guess *models.Guess) (map[string]int, error) {
	row := infra.Db.QueryRow(`
                    SELECT points_table FROM Groups WHERE id = ?
                `, guess.GroupId)

	var rawPointsMap sql.NullString
	err := row.Scan(&rawPointsMap)
	if err != nil {
		return nil, err
	}
	var pointsTable map[string]int
	err = json.Unmarshal([]byte(rawPointsMap.String), &pointsTable)
	if err != nil {
		return nil, err
	}
	return pointsTable, nil

}
