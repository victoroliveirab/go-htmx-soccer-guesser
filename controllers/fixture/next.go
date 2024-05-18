package fixture

import (
	"net/http"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

var NextFixtures http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	tomorrow := now.Add(24 * time.Hour)

	startTime := now.Unix() * 1000
	endTime := tomorrow.Unix() * 1000

	query := `
        SELECT f.id, f.league_id, f.league_season_id, f.home_team_id, ho.name,
               ho.logo_url, f.away_team_id, aw.name, aw.logo_url,
               f.timestamp_numb, f.match_date, f.status
        FROM Fixtures f
        JOIN Teams ho ON f.home_team_id = ho.id
        JOIN Teams aw ON f.away_team_id = aw.id
        WHERE timestamp_numb BETWEEN ? AND ?
    `
	rows, err := infra.Db.Query(query, startTime, endTime)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fixtures := []models.Fixture{}

	for rows.Next() {
		var fixture models.Fixture
		// var homeTeam models.Team
		// var awayTeam models.Team

		err := rows.Scan(&fixture.Id, &fixture.LeagueId, &fixture.LeagueSeasonId, &fixture.HomeTeam.Id, &fixture.HomeTeam.Name, &fixture.HomeTeam.LogoUrl, &fixture.AwayTeam.Id, &fixture.AwayTeam.Name, &fixture.AwayTeam.LogoUrl, &fixture.TimestampNumb, &fixture.MatchDate, &fixture.Status)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fixtures = append(fixtures, fixture)
	}

	data := struct {
		Fixtures []models.Fixture
		HideNav  bool
	}{
		Fixtures: fixtures,
		HideNav:  false,
	}

	lib.RenderTemplate(w, "fixtures/next.html", data)
})
