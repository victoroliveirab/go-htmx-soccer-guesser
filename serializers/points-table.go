package serializers

import (
	"encoding/json"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/constants"
)

func ParsePointsTable(marshelledPointsTable string) (map[string]constants.Outcome, error) {
	var object map[string]int
	err := json.Unmarshal([]byte(marshelledPointsTable), &object)
	if err != nil {
		return nil, err
	}

	var newObject = make(map[string]constants.Outcome, 0)

	for key, value := range object {
		newObject[key] = constants.Outcome(value)
	}
	return newObject, nil
}
