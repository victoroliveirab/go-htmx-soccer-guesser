package fixture

import (
	"fmt"
	"net/http"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

type FixtureView struct {
	Id              int
	LeagueName      string
	Round           int
	FormattedDate   string
	FormattedTime   string
	HomeTeamId      int
	HomeTeamName    string
	HomeTeamLogoUrl string
	HomeTeamScore   int
	AwayTeamId      int
	AwayTeamName    string
	AwayTeamLogoUrl string
	AwayTeamScore   int
}

type DateKey string
type LeagueNameKey string
type RoundKey int
type FixtureViewMap map[DateKey]map[LeagueNameKey]map[RoundKey][]FixtureView

var NextFixtures http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	// threeDaysAhead := now.Add(3 * 24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	startTime := now.Unix() * 1000
	// endTime := threeDaysAhead.Unix() * 1000
	endTime := tomorrow.Unix() * 1000

	query := `
        SELECT f.id, f.league_id, le.name, f.league_season_id, f.home_team_id,
        ho.name, ho.logo_url, f.away_team_id, aw.name, aw.logo_url,
        f.timestamp_numb, f.match_date, f.status, f.round
        FROM Fixtures f
        JOIN Leagues le ON f.league_id = le.id
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

	fixturesView := make(FixtureViewMap)

	for rows.Next() {
		var fixture models.Fixture

		err := rows.Scan(
			&fixture.Id, &fixture.League.Id, &fixture.League.Name,
			&fixture.LeagueSeasonId, &fixture.HomeTeam.Id,
			&fixture.HomeTeam.Name, &fixture.HomeTeam.LogoUrl,
			&fixture.AwayTeam.Id, &fixture.AwayTeam.Name,
			&fixture.AwayTeam.LogoUrl, &fixture.TimestampNumb,
			&fixture.MatchDate, &fixture.Status, &fixture.Round,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		datetime := time.Unix(int64(fixture.TimestampNumb/1000), 0)
		var formattedDate DateKey = DateKey(datetime.Format("02/01/2006"))
		formattedTime := datetime.Format("15:04")
		leagueName := LeagueNameKey(fixture.League.Name)
		round := RoundKey(fixture.Round)

		fixtureView := FixtureView{
			Id:              fixture.Id,
			LeagueName:      fixture.League.Name,
			Round:           fixture.Round,
			FormattedDate:   fmt.Sprintf("%s - %s", formattedDate, formattedTime),
			FormattedTime:   formattedTime,
			HomeTeamId:      fixture.HomeTeam.Id,
			HomeTeamName:    fixture.HomeTeam.Name,
			HomeTeamLogoUrl: fixture.HomeTeam.LogoUrl,
			HomeTeamScore:   fixture.HomeScore,
			AwayTeamId:      fixture.AwayTeam.Id,
			AwayTeamName:    fixture.AwayTeam.Name,
			AwayTeamLogoUrl: fixture.AwayTeam.LogoUrl,
			AwayTeamScore:   fixture.AwayScore,
		}

		byDate, existByDate := fixturesView[formattedDate]
		if !existByDate {
			fixturesView[formattedDate] = make(
				map[LeagueNameKey]map[RoundKey][]FixtureView,
			)
			byDate = fixturesView[formattedDate]
		}

		byLeagueName, existByLeagueName := byDate[leagueName]
		if !existByLeagueName {
			byDate[leagueName] = make(
				map[RoundKey][]FixtureView,
			)
			byLeagueName = byDate[leagueName]
		}

		byRound, existByRound := byLeagueName[round]
		if !existByRound {
			byRound = make([]FixtureView, 0)
			byLeagueName[round] = byRound
		}

		byRound = append(
			byRound,
			fixtureView,
		)
		byLeagueName[round] = byRound
	}

	data := struct {
		FixturesViewMap FixtureViewMap
		HideNav         bool
	}{
		FixturesViewMap: fixturesView,
		HideNav:         false,
	}

	lib.RenderTemplate(w, "fixtures/next.html", data)
})
