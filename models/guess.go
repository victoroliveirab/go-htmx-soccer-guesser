package models

import (
	"database/sql"
	"fmt"
)

type SQLGuess struct {
	Id        int
	UserId    int
	GroupId   int
	FixtureId int
	Locked    bool
	HomeGoals int
	AwayGoals int
	Points    int
	CreatedAt int
	UpdatedAt int
	Outcome   string
}

type Guess struct {
	Id        int
	UserId    int
	GroupId   int
	GroupName string
	Fixture   *Fixture
	Locked    bool
	HomeGoals int
	AwayGoals int
	Points    int
	CreatedAt int
	UpdatedAt int
	Outcome   string
}

const (
	QUERY_GET_GUESSES_BY_FIXTURE_ID_AND_USER_ID = `
        SELECT guess.id, guess.group_id, "group".name, guess.locked,
               guess.home_goals, guess.away_goals, guess.points,
               guess.created_at, guess.updated_at, guess.outcome
        FROM Guesses guess
        JOIN Groups "group" ON guess.group_id = "group".id
        WHERE guess.fixture_id = ?
        AND guess.user_id = ?
    `
)

func GetGuessById(db *sql.DB, id int64) *Guess {
	var guess Guess
	var lockedInt int
	var outcome sql.NullString
	row := db.QueryRow("SELECT * FROM Guesses WHERE id = $1", id)
	if err := row.Scan(
		&guess.Id, &guess.UserId, &guess.GroupId, &lockedInt, &guess.HomeGoals,
		&guess.AwayGoals, &guess.Points, &guess.CreatedAt, &guess.UpdatedAt,
		&outcome,
	); err != nil {
		return nil
	}

	guess.Locked = lockedInt == 1
	guess.Outcome = outcome.String

	return &guess
}

func GetGuessesByFixtureId(db *sql.DB, userId, fixtureId int64) ([]*Guess, error) {
	rows, err := db.Query(QUERY_GET_GUESSES_BY_FIXTURE_ID_AND_USER_ID, fixtureId, userId)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	guesses := make([]*Guess, 0)

	for rows.Next() {
		var guess Guess
		var lockedInt int
		err := rows.Scan(
			&guess.Id, &guess.GroupId, &guess.GroupName, &lockedInt,
			&guess.HomeGoals, &guess.AwayGoals, &guess.Points,
			&guess.CreatedAt, &guess.UpdatedAt, &guess.Outcome,
		)

		if err != nil {
			return nil, err
		}

		guess.Locked = lockedInt == 1

		fmt.Println(guess)
		fmt.Println("=============")

		guesses = append(guesses, &guess)
	}

	return guesses, nil
}
