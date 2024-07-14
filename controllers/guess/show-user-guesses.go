package guess

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/constants"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/serializers"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

const (
	GET_USER_RESULTS_BY_GROUP = `
        SELECT fixture.id, home.name, home.logo_url, away.name, away.logo_url, fixture.home_score,
               fixture.away_score, guess.home_goals, guess.away_goals, guess.points,
               guess.created_at, guess.updated_at, guess.outcome, guess.counted
        FROM Guesses guess
        JOIN Fixtures fixture ON fixture.id = guess.fixture_id
        JOIN Teams home ON fixture.home_team_id = home.id
        JOIN Teams away ON fixture.away_team_id = away.id
        WHERE guess.group_id = ? AND guess.user_id = ? AND guess.locked = 1
    `
)

type ResultsStats map[string]int

var GetUserGuessesByGroup http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.ParseInt(r.PathValue("groupId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: use goroutines instead of doing sequentially
	guesses, err := getGuesses(groupId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pointsTable, err := getPointsTable(groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PointsTable": pointsTable,
	}

	query := r.URL.Query()
	isPartial := query.Get("partial") == "1"

	if isPartial {
		stats := processGuessesToPartial(guesses)
		data["Statistics"] = stats
		fmt.Println(data)
		tmpl := templates.LoadPartial("user-guesses-by-group", "groups/show/_user-results.html")
		tmpl.ExecutePartial(w, r, "partial", data)
		return
	}

	fmt.Println(guesses[0])

	// http.Error(w, "", http.StatusNotImplemented)
})

func getGuesses(groupId, userId int64) ([]*models.Guess, error) {
	rows, err := infra.Db.Query(GET_USER_RESULTS_BY_GROUP, groupId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guesses = make([]*models.Guess, 0)
	for rows.Next() {
		var guess models.Guess
		var sqlPoints sql.NullInt64
		var outcome sql.NullInt64
		guess.Fixture = &models.Fixture{
			HomeTeam: models.Team{},
			AwayTeam: models.Team{},
		}
		err := rows.Scan(
			&guess.Fixture.Id, &guess.Fixture.HomeTeam.Name, &guess.Fixture.HomeTeam.LogoUrl,
			&guess.Fixture.AwayTeam.Name, &guess.Fixture.AwayTeam.LogoUrl, &guess.Fixture.HomeScore,
			&guess.Fixture.AwayScore, &guess.HomeGoals, &guess.AwayGoals, &sqlPoints,
			&guess.CreatedAt, &guess.UpdatedAt, &outcome, &guess.Counted,
		)
		if err != nil {
			return nil, err
		}

		guess.Points = int(sqlPoints.Int64)
		guess.Outcome = constants.NormalizeOutcome(outcome)
		guess.Locked = false

		guesses = append(guesses, &guess)
	}

	return guesses, nil
}

func getPointsTable(groupId int64) (map[string]int, error) {
	var rawPointsTable string
	row := infra.Db.QueryRow("SELECT points_table FROM Groups WHERE id = ?", groupId)
	if err := row.Scan(&rawPointsTable); err != nil {
		return nil, err
	}
	return serializers.ParsePointsTable(rawPointsTable)
}

func processGuessesToPartial(guesses []*models.Guess) ResultsStats {
	stats := make(ResultsStats, 0)
	for key := range constants.OutcomesMap {
		stats[key] = 0
	}

	for _, guess := range guesses {
		if guess.Counted == 0 {
			continue
		}
		stats[guess.Outcome] += 1
	}

	return stats
}
