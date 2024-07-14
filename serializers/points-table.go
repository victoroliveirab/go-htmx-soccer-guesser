package serializers

import (
	"encoding/json"
)

func ParsePointsTable(marshelledPointsTable string) (map[string]int, error) {
	var object map[string]int
	err := json.Unmarshal([]byte(marshelledPointsTable), &object)
	if err != nil {
		return nil, err
	}

	return object, nil
}
