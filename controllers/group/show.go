package group

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/serializers"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/templates"
)

const (
	GET_IS_USER_FROM_GROUP = `SELECT COUNT(*) FROM User_Groups WHERE user_id = ? AND group_id = ?`
)
const (
	GET_GROUP = `
        SELECT id, name, description, admin, points_table, ranking,
        ranking_up_to_date, created_at, updated_at
        FROM Groups WHERE id = ?
    `
)
const GET_GROUP_PARTICIPANTS = `
    SELECT user.id, user.username, user.name
    FROM Users user
    JOIN User_Groups ug ON ug.user_id = user.id
    WHERE ug.group_id = ?
`

type RankingEntry struct {
	Points int
	User   *models.User
}

var Show http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tmpl := templates.LoadTemplate("show-group", "groups/show.html")

	ctx := r.Context()
	userID := int64(ctx.Value("UserID").(int))

	row := infra.Db.QueryRow(GET_IS_USER_FROM_GROUP, userID, id)
	var count int
	if err := row.Scan(&count); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	row = infra.Db.QueryRow(GET_GROUP, id)
	var group models.GroupWithParticipants
	if err = row.Scan(&group.Id, &group.Name, &group.Description,
		&group.AdminId, &group.RawPointsTable, &group.RawRanking,
		&group.RankingUpToDate, &group.CreatedAt, &group.UpdatedAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Group":                group,
		"PointsTableAvailable": false,
		"RankingAvailable":     false,
	}

	pointsTable, err := serializers.ParsePointsTable(group.RawPointsTable)

	if err != nil {
		tmpl.Execute(w, r, data)
		return
	}

	group.PointsTable = pointsTable
	data["PointsTableAvailable"] = true

	rows, err := infra.Db.Query(GET_GROUP_PARTICIPANTS, id)
	defer rows.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group.Users = make([]*models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Username, &user.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		group.Users = append(group.Users, &user)
	}

	group.NumberOfMembers = len(group.Users)

	orderedRanking, err := serializers.ParseGroupRanking(group.RawRanking)

	if err != nil {
		tmpl.Execute(w, r, data)
		return
	}

	ranking := make([]*RankingEntry, 0)

	for _, entry := range orderedRanking {
		var currentUser *models.User
		for _, user := range group.Users {
			if user.Id == entry.UserId {
				currentUser = user
				break
			}
		}
		ranking = append(ranking, &RankingEntry{
			Points: entry.Points,
			User:   currentUser,
		})
	}

	data["Ranking"] = ranking
	data["RankingAvailable"] = true

	fmt.Println(data)
	fmt.Println("HEREHERHEHREHE")

	tmpl.Execute(w, r, data)
})
