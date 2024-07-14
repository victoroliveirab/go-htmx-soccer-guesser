package models

import (
	"database/sql"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/constants"
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
	Outcome   constants.Outcome
	Counted   int
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
	Counted   int
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
const (
	QUERY_INSERT_GUESS = `
        INSERT INTO Guesses(
            group_id, fixture_id, user_id, locked, home_goals, away_goals
        ) VALUES(?, ?, ?, ?, ?, ?)
    `
)
const (
	QUERY_GET_POSSIBLE_GUESSES_BY_FIXTURE_ID_AND_USER_ID = `
        SELECT
            "group".id as group_id,
            "group".name as group_name,
            fixture.id AS fixture_id, fixture.home_team_id, home.name AS home_name,
            home.logo_url AS home_logo_url, fixture.away_team_id,
            away.name AS away_name, away.logo_url AS away_logo_url,
            fixture.home_score, fixture.away_score, fixture.home_winner,
            fixture.away_winner, fixture.timestamp_numb,
            guess.id, guess.locked, guess.home_goals, guess.away_goals,
            guess.points, guess.created_at, guess.updated_at, guess.outcome
        FROM
            Groups "group"
        JOIN
            User_Groups user_group ON "group".id = user_group.group_id
        CROSS JOIN
            Fixtures fixture
        LEFT JOIN
            Guesses guess
                ON user_group.user_id = guess.user_id
                AND "group".id = guess.group_id
                AND fixture.id = guess.fixture_id
        JOIN
            Teams home ON fixture.home_team_id = home.id
        JOIN
            Teams away ON fixture.away_team_id = away.id
        WHERE
            fixture.id = ?
            AND user_group.user_id = ?
        `
)

func GetGuessById(db *sql.DB, id int64) (*Guess, error) {
	var guess Guess
	guess.Fixture = &Fixture{
		HomeTeam: Team{},
		AwayTeam: Team{},
	}
	var homeScore, awayScore sql.NullInt64
	var homeWinner, awayWinner sql.NullInt64
	var lockedInt int
	var outcome sql.NullInt64
	row := db.QueryRow(`
        SELECT guess.id, guess.group_id, fixture.id, home.name, home.logo_url,
        away.name, away.logo_url, fixture.home_score, fixture.away_score,
        fixture.home_winner, fixture.away_winner, guess.locked, guess.home_goals,
        guess.away_goals, guess.points, guess.created_at, guess.updated_at, guess.outcome
        FROM Guesses guess
        JOIN Fixtures fixture ON fixture.id = guess.fixture_id
        JOIN Teams home ON fixture.home_team_id = home.id
        JOIN Teams away ON fixture.away_team_id = away.id
        WHERE guess.id = ?
    `, id)
	if err := row.Scan(
		&guess.Id, &guess.GroupId, &guess.Fixture.Id, &guess.Fixture.HomeTeam.Name,
		&guess.Fixture.HomeTeam.LogoUrl, &guess.Fixture.AwayTeam.Name,
		&guess.Fixture.AwayTeam.LogoUrl, &homeScore, &awayScore,
		&homeWinner, &awayWinner, &lockedInt, &guess.HomeGoals, &guess.AwayGoals,
		&guess.Points, &guess.CreatedAt, &guess.UpdatedAt, &outcome,
	); err != nil {
		return nil, err
	}

	if homeWinner.Int64 == 1 {
		guess.Fixture.Winner = "Home"
	}
	if awayWinner.Int64 == 1 {
		guess.Fixture.Winner = "Away"
	}

	guess.Locked = lockedInt == 1
	guess.Outcome = constants.NormalizeOutcome(outcome)

	return &guess, nil
}

func GetPossibleGuessesByFixtureId(db *sql.DB, userId, fixtureId int64) ([]*Guess, error) {
	rows, err := db.Query(
		QUERY_GET_POSSIBLE_GUESSES_BY_FIXTURE_ID_AND_USER_ID,
		fixtureId, userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guesses := make([]*Guess, 0)

	now := time.Now().Unix()

	for rows.Next() {
		var guess Guess
		var guessId sql.NullInt64
		var guessHomeGoals, guessAwayGoals sql.NullInt64
		var homeScore, awayScore sql.NullInt64
		var homeWinner, awayWinner sql.NullInt64
		var lockedInt sql.NullInt64
		var sqlPoints, sqlCreatedAt, sqlUpdatedAt sql.NullInt64
		var outcome sql.NullInt64

		guess.Fixture = &Fixture{
			HomeTeam: Team{},
			AwayTeam: Team{},
		}

		err := rows.Scan(
			&guess.GroupId, &guess.GroupName, &guess.Fixture.Id, &guess.Fixture.HomeTeam.Id,
			&guess.Fixture.HomeTeam.Name, &guess.Fixture.HomeTeam.LogoUrl,
			&guess.Fixture.AwayTeam.Id, &guess.Fixture.AwayTeam.Name,
			&guess.Fixture.AwayTeam.LogoUrl, &homeScore, &awayScore,
			&homeWinner, &awayWinner, &guess.Fixture.TimestampNumb,
			&guessId, &lockedInt, &guessHomeGoals,
			&guessAwayGoals, &sqlPoints, &sqlCreatedAt, &sqlUpdatedAt,
			&outcome,
		)
		if err != nil {
			return nil, err
		}

		guess.Id = int(guessId.Int64)
		guess.UserId = int(userId)
		guess.HomeGoals = int(guessHomeGoals.Int64)
		guess.AwayGoals = int(guessAwayGoals.Int64)
		guess.Points = int(sqlPoints.Int64)
		guess.CreatedAt = int(sqlCreatedAt.Int64)
		guess.UpdatedAt = int(sqlUpdatedAt.Int64)
		guess.Outcome = constants.NormalizeOutcome(outcome)

		if homeWinner.Int64 == 1 {
			guess.Fixture.Winner = "Home"
		}
		if awayWinner.Int64 == 1 {
			guess.Fixture.Winner = "Away"
		}

		timestamp := guess.Fixture.TimestampNumb
		if !guessId.Valid {
			guess.Locked = now > int64(timestamp)
		} else {
			guess.Locked = lockedInt.Int64 == 1
		}

		guesses = append(guesses, &guess)
	}
	return guesses, nil
}

func GetGuessesByFixtureId(db *sql.DB, userId, fixtureId int64) ([]*Guess, error) {
	rows, err := db.Query(QUERY_GET_GUESSES_BY_FIXTURE_ID_AND_USER_ID, fixtureId, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guesses := make([]*Guess, 0)

	for rows.Next() {
		var guess Guess
		var lockedInt int
		var outcome sql.NullInt64
		err := rows.Scan(
			&guess.Id, &guess.GroupId, &guess.GroupName, &lockedInt,
			&guess.HomeGoals, &guess.AwayGoals, &guess.Points,
			&guess.CreatedAt, &guess.UpdatedAt, &outcome,
		)

		if err != nil {
			return nil, err
		}

		guess.Locked = lockedInt == 1
		guess.Outcome = constants.NormalizeOutcome(outcome)
		guesses = append(guesses, &guess)
	}

	return guesses, nil
}

func SaveGuess(db *sql.DB, guess *SQLGuess) (int64, error) {
	locked := 0
	if guess.Locked {
		locked = 1
	}
	row, err := db.Exec(
		QUERY_INSERT_GUESS,
		guess.GroupId, guess.FixtureId, guess.UserId, locked, guess.HomeGoals,
		guess.AwayGoals,
	)

	if err != nil {
		return -1, err
	}

	return row.LastInsertId()
}
