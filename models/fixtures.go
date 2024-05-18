package models

type RawFixture struct {
	Id             int
	ApiFootballId  int
	LeagueId       int
	LeagueSeasonId int
	HomeTeamId     int
	AwayTeamId     int
	TimestampNumb  int
	MatchDate      string
	Status         int
	Referee        string
	HomeScore      int
	AwayScore      int
	Round          int
	CreatedAt      int
	UpdatedAt      int
}

type Fixture struct {
	Id             int
	League         League
	LeagueSeasonId int
	HomeTeam       Team
	AwayTeam       Team
	TimestampNumb  int
	MatchDate      string
	Status         int
	Referee        string
	HomeScore      int
	AwayScore      int
	Round          int
	CreatedAt      int
	UpdatedAt      int
}
