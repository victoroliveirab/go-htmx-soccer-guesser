package serializers

import (
	"encoding/json"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

func ParseStandings(marshelledStandings string) ([]*models.Standing, error) {
	var object []*models.Standing
	err := json.Unmarshal([]byte(marshelledStandings), &object)
	if err != nil {
		return nil, err
	}
	return object, nil
}
