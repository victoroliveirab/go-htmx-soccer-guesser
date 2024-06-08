package main

import (
	"github.com/victoroliveirab/go-htmx-soccer-guesser/constants"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

func init() {
	infra.DbConnect("file:local.db")
}

func main() {
	db := infra.Db

	rows, err := db.Query(`
        SELECT guess.id, guess.home_goals, guess.away_goals, fixture.home_score, fixture.away_score
        FROM Guesses guess
        JOIN Fixtures fixture ON fixture.id = guess.fixture_id
        WHERE guess.locked = 1
        AND fixture.status IN (0, 2)
    `)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var guess models.Guess
		var fixture models.Fixture
		err := rows.Scan(&guess.Id, &guess.HomeGoals, &guess.AwayGoals, &fixture.HomeScore, &fixture.AwayScore)
		if err != nil {
			panic(err)
		}

		outcome := lib.DefineOutcome(&guess, &fixture)
		points := constants.DefaultOutcomePointsMap[outcome]

		_, err = tx.Exec("UPDATE guesses SET outcome = ?, points = ? WHERE id = ?", outcome, points, guess.Id)

		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
