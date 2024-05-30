package models

import "database/sql"

type Group struct {
	Id              int
	Name            string
	Description     string
	NumberOfMembers int
}

type GroupWithParticipants struct {
	Id          int
	Name        string
	Description string
	Users       []*User
}

var groupWithParticipantsQuery = `
	SELECT Users.id AS user_id, Users.username, Users.name AS user_name, Users.email,
	       Groups.id AS group_id, Groups.name AS group_name, Groups.description
	FROM Users
	JOIN User_Groups ON Users.id = User_Groups.user_id
	JOIN Groups ON User_Groups.group_id = Groups.id
	WHERE Groups.id = ?;
	`

var groupsOfUserQuery = `
	SELECT g.id AS group_id, g.name AS group_name, g.description,
	       (SELECT COUNT(ug.user_id)
	        FROM User_Groups ug
	        WHERE ug.group_id = g.id) AS user_count
	FROM Groups g
	JOIN User_Groups ug ON g.id = ug.group_id
	WHERE ug.user_id = ?;
	`

func GetGroupWithParticipantsById(db *sql.DB, id int64) *GroupWithParticipants {
	var group GroupWithParticipants

	rows, err := db.Query(groupWithParticipantsQuery, id)

	if err != nil {
		return nil
	}
	defer rows.Close()

	group.Users = make([]*User, 0)

	for rows.Next() {
		var user User
		var description sql.NullString

		err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Email, &group.Id, &group.Name, &description)
		if err != nil {
			return nil
		}

		group.Description = description.String

		group.Users = append(group.Users, &user)
	}
	return &group
}

func GetGroupsAssociatedWithUserId(db *sql.DB, userId int64) []*Group {
	var groups []*Group
	rows, err := db.Query(groupsOfUserQuery, userId)

	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		var description sql.NullString

		err := rows.Scan(&group.Id, &group.Name, &description, &group.NumberOfMembers)
		if err != nil {
			return nil
		}

		group.Description = description.String

		groups = append(groups, &group)
	}
	return groups
}
