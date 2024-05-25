package models

import (
	"database/sql"
	"fmt"
)

type SQLFixture struct {
	Id            int
	ApiFootballId int
	LeagueId      sql.NullInt64
	Season        sql.NullInt64
	HomeTeamId    int
	AwayTeamId    int
	TimestampNumb sql.NullInt64
	VenueId       sql.NullInt64
	Status        int
	Referee       sql.NullString
	HomeScore     sql.NullInt64
	AwayScore     sql.NullInt64
	HomeWinner    sql.NullInt64
	AwayWinner    sql.NullInt64
	Round         sql.NullString
	CreatedAt     int
	UpdatedAt     int
}

type Fixture struct {
	Id            int
	League        League
	Season        int
	HomeTeam      Team
	AwayTeam      Team
	TimestampNumb int
	Status        int
	Referee       string
	HomeScore     int
	AwayScore     int
	Winner        string
	Round         string
	CreatedAt     int
	UpdatedAt     int
}

const (
	QUERY_GET_COMPLETE_FIXTURE = `
        SELECT fixture.id, fixture.league_id, league.name, league.logo_url,
               league.country, league.country_flag_url, fixture.season,
               fixture.home_team_id, home.name, home.logo_url,
                fixture.away_team_id, away.name, away.logo_url,
               fixture.timestamp_numb, fixture.venue_id,
               fixture.status, fixture.referee, fixture.home_score,
               fixture.away_score, fixture.home_winner, fixture.away_winner,
               fixture.round, fixture.created_at, fixture.updated_at
        FROM Fixtures fixture
        JOIN Leagues league ON fixture.league_id = league.id
        JOIN League_Seasons league_season ON fixture.league_season_id = league_season.id
        JOIN Teams home ON fixture.home_team_id = home.id
        JOIN Teams away ON fixture.away_team_id = away.id
    `
)

func FromSQLFixtureToFixture(sqlFixture *SQLFixture) *Fixture {
	var fixture Fixture
	fixture.Id = sqlFixture.Id
	fixture.TimestampNumb = int(sqlFixture.TimestampNumb.Int64)
	fixture.Status = sqlFixture.Status
	fixture.Referee = sqlFixture.Referee.String
	fixture.HomeScore = int(sqlFixture.HomeScore.Int64)
	fixture.AwayScore = int(sqlFixture.AwayScore.Int64)
	fixture.Winner = ""
	fixture.Round = sqlFixture.Round.String
	fixture.CreatedAt = sqlFixture.CreatedAt
	fixture.UpdatedAt = sqlFixture.UpdatedAt
	if sqlFixture.HomeWinner.Int64 == 1 {
		fixture.Winner = "Home"
	}
	if sqlFixture.AwayWinner.Int64 == 1 {
		fixture.Winner = "Away"
	}
	return &fixture
}

func GetFixtureById(db *sql.DB, id int64) (*Fixture, error) {
	query := fmt.Sprintf(`%s WHERE fixture.id = $1`, QUERY_GET_COMPLETE_FIXTURE)
	row := db.QueryRow(query, id)
	var sqlFixture SQLFixture
	var league SQLLeague
	var sqlHomeTeam SQLTeam
	var sqlAwayTeam SQLTeam
	err := row.Scan(
		&sqlFixture.Id, &sqlFixture.LeagueId, &league.Name, &league.LogoUrl,
		&league.Country, &league.CountryFlagUrl, &sqlFixture.Season,
		&sqlFixture.HomeTeamId,
		&sqlHomeTeam.Name, &sqlHomeTeam.LogoUrl, &sqlFixture.AwayTeamId,
		&sqlAwayTeam.Name, &sqlAwayTeam.LogoUrl, &sqlFixture.TimestampNumb,
		&sqlFixture.VenueId, &sqlFixture.Status,
		&sqlFixture.Referee, &sqlFixture.HomeScore, &sqlFixture.AwayScore,
		&sqlFixture.HomeWinner, &sqlFixture.AwayWinner,
		&sqlFixture.Round, &sqlFixture.CreatedAt, &sqlFixture.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	sqlHomeTeam.Id = sqlFixture.HomeTeamId
	sqlAwayTeam.Id = sqlFixture.AwayTeamId

	homeTeam := FromSQLTeamToTeam(&sqlHomeTeam)
	awayTeam := FromSQLTeamToTeam(&sqlAwayTeam)
	fixture := FromSQLFixtureToFixture(&sqlFixture)
	fixture.HomeTeam = *homeTeam
	fixture.AwayTeam = *awayTeam

	return fixture, nil
}
