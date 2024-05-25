package models

import "database/sql"

type SQLLeague struct {
	Id             int
	ApiFootballId  int
	Name           string
	LogoUrl        sql.NullString
	Country        sql.NullString
	CountryFlagUrl sql.NullString
	LeagueType     sql.NullString
	Meta           sql.NullString
}

type League struct {
	Id             int
	ApiFootballId  int
	Name           string
	LogoUrl        string
	Country        string
	CountryFlagUrl string
	LeagueType     string
	Meta           string
}

func FromSQLLeagueToLeague(sqlLeague *SQLLeague) *League {
	var league League
	league.Id = sqlLeague.Id
	league.ApiFootballId = sqlLeague.ApiFootballId
	league.Name = sqlLeague.Name
	league.LogoUrl = sqlLeague.LogoUrl.String
	league.Country = sqlLeague.Country.String
	league.CountryFlagUrl = sqlLeague.CountryFlagUrl.String
	league.LeagueType = sqlLeague.LeagueType.String
	league.Meta = sqlLeague.Meta.String
	return &league
}
