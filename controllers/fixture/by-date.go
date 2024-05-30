package fixture

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

type FixtureView struct {
	Id              int
	LeagueName      string
	Round           string
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
type RoundKey string
type FixtureMap map[LeagueNameKey]map[RoundKey][]FixtureView

// Date in format YYYY-mm-dd
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func getQueriedDate(r *http.Request) time.Time {
	queryParams := r.URL.Query()
	queriedDate := queryParams.Get("date")
	if !dateRegex.MatchString(queriedDate) {
		return time.Now().UTC().Truncate(24 * time.Hour)
	}
	parsedDate, err := time.Parse("2006-01-02", queriedDate)
	if err != nil {
		return time.Now().UTC().Truncate(24 * time.Hour)
	}
	return parsedDate.Truncate(24 * time.Hour)
}

var FixturesByDate http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	startOfDay := getQueriedDate(r)
	endOfDay := time.Date(startOfDay.Year(), startOfDay.Month(), startOfDay.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)

	// FIXME: generalize this behavior
	// Prevent 21:00 matches of yesterday from showing today (GMT-03)
	startTime := startOfDay.Add(3 * time.Hour).Unix()
	endTime := endOfDay.Unix()

	query := `
        SELECT f.id, f.league_id, le.name, f.season, f.home_team_id,
        ho.name, ho.logo_url, f.away_team_id, aw.name, aw.logo_url,
        f.timestamp_numb, f.status, f.round
        FROM Fixtures f
        JOIN Leagues le ON f.league_id = le.id
        JOIN Teams ho ON f.home_team_id = ho.id
        JOIN Teams aw ON f.away_team_id = aw.id
        WHERE timestamp_numb BETWEEN ? AND ?
        ORDER BY f.timestamp_numb ASC;
    `
	rows, err := infra.Db.Query(query, startTime, endTime)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fixturesView := make(FixtureMap)
	empty := true

	for rows.Next() {
		empty = false
		var fixture models.Fixture

		err := rows.Scan(
			&fixture.Id, &fixture.League.Id, &fixture.League.Name,
			&fixture.Season, &fixture.HomeTeam.Id,
			&fixture.HomeTeam.Name, &fixture.HomeTeam.LogoUrl,
			&fixture.AwayTeam.Id, &fixture.AwayTeam.Name,
			&fixture.AwayTeam.LogoUrl, &fixture.TimestampNumb,
			&fixture.Status, &fixture.Round,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		datetime := time.Unix(int64(fixture.TimestampNumb), 0)
		var formattedDate DateKey = DateKey(datetime.Format("02/01/2006"))
		fmt.Println(formattedDate)
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

		byLeagueName, existByLeagueName := fixturesView[leagueName]
		if !existByLeagueName {
			fixturesView[leagueName] = make(
				map[RoundKey][]FixtureView,
			)
			byLeagueName = fixturesView[leagueName]
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

	data := map[string]interface{}{
		"PrevDate":        startOfDay.AddDate(0, 0, -1).Format("2006-01-02"),
		"NextDate":        startOfDay.AddDate(0, 0, 1).Format("2006-01-02"),
		"Date":            startOfDay.Format("02/01/2006"),
		"FixturesViewMap": fixturesView,
		"NoMatches":       empty,
	}

	// cache page for 1 min
	cacheControlHeader := "public, max-age=60"
	expirationTime := time.Now().Add(time.Minute)

	w.Header().Set("Cache-Control", cacheControlHeader)
	w.Header().Set("Expires", expirationTime.Format(http.TimeFormat))

	lib.RenderTemplate(w, r, "fixtures/index.html", data)
})
