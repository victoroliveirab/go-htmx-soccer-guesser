package guess

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

// Create a new guess entity and return the partial for the modal view
var Create http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	tmpl := templates.LoadPartial("form-fixture", "fixtures/_fixture-form.html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	groupIdStr := r.FormValue("group-id")
	fixtureIdStr := r.FormValue("fixture-id")
	homeGoalsStr := r.FormValue("home-goals")
	awayGoalsStr := r.FormValue("away-goals")

	if groupIdStr == "" || fixtureIdStr == "" || homeGoalsStr == "" || awayGoalsStr == "" {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	fixtureId, err := strconv.ParseInt(fixtureIdStr, 10, 64)
	if err != nil {
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

	ctx := r.Context()

	sqlGuess := models.SQLGuess{
		GroupId:   int(groupId),
		FixtureId: int(fixtureId),
		UserId:    ctx.Value("UserID").(int),
		HomeGoals: int(homeGoals),
		AwayGoals: int(awayGoals),
		Locked:    false,
	}

	guessId, err := models.SaveGuess(infra.Db, &sqlGuess)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if guessId == -1 {
		http.Error(w, "Resource not found", http.StatusInternalServerError)
		return
	}

	guess, err := models.GetGuessById(infra.Db, guessId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.ExecutePartial(w, r, "fixture-form", guess)
})
