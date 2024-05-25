package models

import (
	"database/sql"
)

type SQLTeam struct {
	Id            int
	ApiFootballId int
	Name          string
	LogoUrl       sql.NullString
}

type Team struct {
	Id            int
	ApiFootballId int
	Name          string
	LogoUrl       string
}

func FromSQLTeamToTeam(sqlTeam *SQLTeam) *Team {
	var team Team
	team.Id = sqlTeam.Id
	team.ApiFootballId = sqlTeam.ApiFootballId
	team.Name = sqlTeam.Name
	team.LogoUrl = sqlTeam.LogoUrl.String
	return &team
}

func CreateTeam(db *sql.DB, apiFootballId int, name, logoUrl string) error {
	_, err := db.Exec(
		"INSERT INTO Teams(api_football_id, name, logo_url) VALUES($1, $2, $3)",
		apiFootballId,
		name,
		logoUrl,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetAllTeams(db *sql.DB) ([]Team, error) {
	rows, err := db.Query("SELECT * FROM Teams")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var teams []Team

	for rows.Next() {
		var team Team
		var optionalLogoUrl sql.NullString
		if err := rows.Scan(&team.Id, &team.ApiFootballId, &team.Name, &optionalLogoUrl); err != nil {
			return nil, err
		}
		team.LogoUrl = optionalLogoUrl.String
		teams = append(teams, team)
	}

	return teams, nil
}
