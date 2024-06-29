package fixture

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

var partials = []string{"fixtures/_fixture-form.html", "fixtures/_fixture.html"}

func getFixtureAndGuesses(fixtureId, userId int64) (*models.Fixture, []*models.Guess, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	errs := make(chan error, 2)

	var guesses []*models.Guess
	var fixture *models.Fixture

	go func() {
		defer wg.Done()

		rows, err := models.GetPossibleGuessesByFixtureId(infra.Db, userId, fixtureId)

		if err != nil {
			errs <- err
			return
		}

		guesses = rows
	}()

	go func() {
		defer wg.Done()

		row, err := models.GetFixtureById(infra.Db, fixtureId)

		if err != nil {
			errs <- err
			return
		}

		fixture = row
	}()

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return nil, nil, err
		}
	}

	for guessIdx := range guesses {
		guesses[guessIdx].Fixture = fixture
	}

	return fixture, guesses, nil
}

var ViewFixture http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := int64(r.Context().Value("UserID").(int))

	query := r.URL.Query()
	isModal := query.Get("modal") == "1"

	fixture, guesses, err := getFixtureAndGuesses(id, userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Fixture": fixture,
		"Guesses": guesses,
	}

	if isModal {
		tmpl := templates.LoadPartial("show-fixture", partials...)
		tmpl.ExecutePartial(w, r, "fixture-information", data)
		return
	}

	tmpl := templates.LoadTemplate("show-fixture", "fixtures/show.html")
	tmpl.Execute(w, r, data)
})
