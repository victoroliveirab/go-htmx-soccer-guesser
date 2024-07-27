package guess

import (
	"net/http"
	"strconv"
	"time"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

var Update http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	tmpl := templates.LoadPartial("form-fixture", "fixtures/_fixture-form.html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	homeGoalsStr := r.FormValue("home-goals")
	awayGoalsStr := r.FormValue("away-goals")

	if homeGoalsStr == "" || awayGoalsStr == "" {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	homeGoals, err := strconv.ParseInt(homeGoalsStr, 10, 64)
	if err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	awayGoals, err := strconv.ParseInt(awayGoalsStr, 10, 64)
	if err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	userId := int64(r.Context().Value("UserID").(int))
	now := time.Now().Unix()

	updatedRow, err := infra.Db.Exec(`
       UPDATE guesses SET home_goals = ?, away_goals = ?, updated_at = ?
       WHERE id = ? AND user_id = ?
    `, homeGoals, awayGoals, now, id, userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := updatedRow.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	var guess models.Guess
	guess.Fixture = &models.Fixture{
		HomeTeam: models.Team{},
		AwayTeam: models.Team{},
	}

	row := infra.Db.QueryRow(`
        SELECT guess.id, fixture.id, guess.group_id, "group".name, guess.locked,
               home.name, home.logo_url, away.name, away.logo_url,
               guess.home_goals, guess.away_goals, guess.points,
               guess.created_at, guess.updated_at
        FROM Guesses guess
        JOIN Groups "group" ON guess.group_id = "group".id
        JOIN Fixtures fixture ON guess.fixture_id = fixture.id
        JOIN Teams home ON fixture.home_team_id = home.id
        JOIN Teams away ON fixture.away_team_id = away.id
        WHERE guess.id = ?
    `, id)

	var locked int

	if err = row.Scan(&guess.Id, &guess.Fixture.Id, &guess.GroupId, &guess.GroupName, &locked,
		&guess.Fixture.HomeTeam.Name, &guess.Fixture.HomeTeam.LogoUrl,
		&guess.Fixture.AwayTeam.Name, &guess.Fixture.AwayTeam.LogoUrl,
		&guess.HomeGoals, &guess.AwayGoals, &guess.Points, &guess.CreatedAt, &guess.UpdatedAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	guess.Locked = locked == 1

	tmpl.ExecutePartial(w, r, "fixture-form", guess)
})
