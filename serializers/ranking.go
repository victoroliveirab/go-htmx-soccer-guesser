package serializers

import (
	"encoding/json"
	"sort"
	"strconv"
)

type UserIdWithPoint struct {
	UserId int
	Points int
}

type OrderedRanking []*UserIdWithPoint

// Sort interface - descending order
func (a OrderedRanking) Len() int           { return len(a) }
func (a OrderedRanking) Less(i, j int) bool { return a[i].Points > a[j].Points }
func (a OrderedRanking) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func ParseGroupRanking(marshelledRanking string) (OrderedRanking, error) {
	var object map[string]int
	err := json.Unmarshal([]byte(marshelledRanking), &object)
	if err != nil {
		return nil, err
	}
	var newObject = make(OrderedRanking, 0)

	for keyStr, value := range object {
		key, err := strconv.ParseInt(keyStr, 10, 64)
		if err != nil {
			return nil, err
		}
		newObject = append(newObject, &UserIdWithPoint{
			UserId: int(key),
			Points: value,
		})
	}

	sort.Sort(OrderedRanking(newObject))

	return newObject, nil
}
