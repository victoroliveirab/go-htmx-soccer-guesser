package league

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/serializers"
)

const (
	GET_LEAGUE_WITH_STANDINGS_QUERY = `
        SELECT l.id, l.name, l.logo_url, l.country, l.country_flag_url, l.league_type, ls.standings
        FROM Leagues_Seasons ls
        JOIN Leagues l ON ls.league_id = l.id
        WHERE ls.id = ? AND ls.season = ?
    `
)

var ViewLeagueWithStandings http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sqlLeague models.SQLLeague
	var rawStandings sql.NullString

	row := infra.Db.QueryRow(GET_LEAGUE_WITH_STANDINGS_QUERY, id, 2024)
	if err := row.Scan(
		&sqlLeague.Id, &sqlLeague.Name, &sqlLeague.LogoUrl, &sqlLeague.Country,
		&sqlLeague.CountryFlagUrl, &sqlLeague.LeagueType, &rawStandings,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	league := models.FromSQLLeagueToLeague(&sqlLeague)

	data := map[string]interface{}{
		"League":             league,
		"StandingsAvailable": false,
	}

	if !rawStandings.Valid {
		lib.RenderTemplate(w, r, "leagues/show.html", data)
		return
	}

	standings, err := serializers.ParseStandings(rawStandings.String)
	if err != nil {
		lib.RenderTemplate(w, r, "leagues/show.html", data)
		return
	}

	orderClauses := ""
	for i, standing := range standings {
		if i > 0 {
			orderClauses += ","
		}
		orderClauses += fmt.Sprintf("(%d, %d)", standing.TeamId, i+1)
	}

	query := fmt.Sprintf(`
            WITH ordered_teams (id, seq) AS (VALUES %s)
            SELECT team.id, team.name, team.logo_url
            FROM Teams team
            JOIN ordered_teams o ON team.id = o.id
            ORDER BY o.seq
        `, orderClauses)

	var teams []*models.Team

	rows, err := infra.Db.Query(query)
	if err != nil {
		lib.RenderTemplate(w, r, "leagues/show.html", data)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var team models.Team
		err = rows.Scan(&team.Id, &team.Name, &team.LogoUrl)
		if err != nil {
			lib.RenderTemplate(w, r, "leagues/show.html", data)
			return
		}
		teams = append(teams, &team)
	}

	for i, team := range teams {
		standings[i].Team = *team
	}

	data["StandingsAvailable"] = true
	data["Standings"] = standings

	lib.RenderTemplate(w, r, "leagues/show.html", data)
})
