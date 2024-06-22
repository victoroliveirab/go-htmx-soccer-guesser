package models

import (
	"database/sql"
	"fmt"
)

type SQLFixture struct {
	Id             int
	ApiFootballId  int
	LeagueSeasonId sql.NullInt64
	HomeTeamId     int
	AwayTeamId     int
	TimestampNumb  sql.NullInt64
	VenueId        sql.NullInt64
	Status         int
	Referee        sql.NullString
	HomeScore      sql.NullInt64
	AwayScore      sql.NullInt64
	HomeWinner     sql.NullInt64
	AwayWinner     sql.NullInt64
	Round          sql.NullString
	CreatedAt      int
	UpdatedAt      int
}

type Fixture struct {
	Id            int
	League        League
	Season        int
	HomeTeam      Team
	AwayTeam      Team
	TimestampNumb int
	Status        string
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
        SELECT fixture.id, league_season.league_id, league.name, league.logo_url,
               league.country, league.country_flag_url, league_season.season,
               fixture.home_team_id, home.name, home.logo_url,
               fixture.away_team_id, away.name, away.logo_url,
               fixture.timestamp_numb, fixture.venue_id,
               fixture.status, fixture.referee, fixture.home_score,
               fixture.away_score, fixture.home_winner, fixture.away_winner,
               fixture.round, fixture.created_at, fixture.updated_at
        FROM Fixtures fixture
        JOIN Leagues_Seasons league_season ON league_season.id = fixture.league_season_id
        JOIN Leagues league ON league_season.league_id = league.id
        JOIN Teams home ON fixture.home_team_id = home.id
        JOIN Teams away ON fixture.away_team_id = away.id
    `
)

var STATUS_MAP map[int]string
var isStatusMapInit = false

func initStatusMap() {
	if isStatusMapInit {
		return
	}
	STATUS_MAP = make(map[int]string)
	STATUS_MAP[0] = "FIN"
	STATUS_MAP[1] = "PST"
	STATUS_MAP[2] = "PEN"
	STATUS_MAP[3] = "NST"
}

func FromSQLFixtureToFixture(sqlFixture *SQLFixture) *Fixture {
	var fixture Fixture
	fixture.Id = sqlFixture.Id
	fixture.TimestampNumb = int(sqlFixture.TimestampNumb.Int64)
	fixture.Status = FixtureTranslateStatus(sqlFixture.Status)
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
	var season int
	var sqlHomeTeam SQLTeam
	var sqlAwayTeam SQLTeam
	err := row.Scan(
		&sqlFixture.Id, &sqlFixture.LeagueSeasonId, &league.Name, &league.LogoUrl,
		&league.Country, &league.CountryFlagUrl, &season, &sqlFixture.HomeTeamId,
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
	fixture.Season = season

	return fixture, nil
}

func FixtureTranslateStatus(status int) string {
	initStatusMap()

	stringStatus, exists := STATUS_MAP[status]
	if !exists {
		return "UNK"
	}
	return stringStatus
}
