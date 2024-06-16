package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
)

type Group struct {
	Id      int
	Ranking map[string]int
}

type GuessToUpdate struct {
	Id     int
	UserId int
	Points int
}

type Ranking map[string]int

func init() {
	err := infra.DbConnect("file:local.db?_busy_timeout=5000")
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		panic(err)
	}
}

func main() {
	db := infra.Db
	defer db.Close()

	rows, err := db.Query(`
        SELECT id, ranking FROM Groups WHERE ranking_up_to_date = 0
    `)

	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		panic(err)
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		group, err := parseGroupRow(rows)
		if err != nil {
			log.Fatalf("failed parsing group row: %v", err)
			panic(err)
		}

		groupId := group.Id
		ranking := group.Ranking

		uncountedGuesses, err := getGuessesToUpdate(groupId)

		if err != nil {
			log.Fatalf("failed to get uncounted guesses: %v", err)
			panic(err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("failed to begin transaction: %v", err)
			panic(err)
		}

		for _, guess := range uncountedGuesses {
			guessId := guess.Id
			userId := strconv.FormatInt(int64(guess.UserId), 10)
			points := guess.Points
			ranking[userId] = ranking[userId] + points
			_, err = tx.Exec(`
                UPDATE guesses SET counted = 1 WHERE id = ?
            `, guessId)
			if err != nil {
				log.Fatalf("failed updating guess %d to be counted: %v", guessId, err)
				txErr := tx.Rollback()
				if txErr != nil {
					log.Fatalf("failed rolling back transaction: %v", txErr)
				}
				panic(err)
			}
		}

		log.Println("Ranking after:")
		log.Println(ranking)

		marshalledRanking, err := json.Marshal(ranking)
		if err != nil {
			log.Fatalf("failed marshalling ranking of group %d: %v", groupId, err)
			txErr := tx.Rollback()
			if txErr != nil {
				log.Fatalf("failed rolling back transaction: %v", txErr)
			}
			panic(err)
		}

		_, err = tx.Exec(`
            UPDATE Groups SET ranking = ?, ranking_up_to_date = 1 WHERE id = ?
        `, marshalledRanking, groupId)

		if err != nil {
			log.Fatalf("failed saving marshalled ranking of group %d to db: %v", groupId, err)
			txErr := tx.Rollback()
			if txErr != nil {
				log.Fatalf("failed rolling back transaction: %v", txErr)
			}
			panic(err)
		}

		err = tx.Commit()

		if err != nil {
			log.Fatalf("failed commiting transaction: %v", err)
			txErr := tx.Rollback()
			if txErr != nil {
				log.Fatalf("failed rolling back transaction: %v", txErr)
			}
			panic(err)
		}

		i++
	}

	log.Printf("finished updating %d groups", i+1)

}

func parseGroupRow(rows *sql.Rows) (*Group, error) {
	var id int
	var rawRanking string
	err := rows.Scan(&id, &rawRanking)
	if err != nil {
		return nil, err
	}

	var ranking map[string]int
	err = json.Unmarshal([]byte(rawRanking), &ranking)
	if err != nil {
		return nil, err
	}

	return &Group{
		Id:      id,
		Ranking: ranking,
	}, nil
}

func getGuessesToUpdate(groupId int) ([]*GuessToUpdate, error) {
	db := infra.Db
	rows, err := db.Query(`
        SELECT id, user_id, points FROM Guesses WHERE group_id = ? AND counted = 0
    `, groupId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var guesses []*GuessToUpdate
	guesses = make([]*GuessToUpdate, 0)

	for rows.Next() {
		var guess GuessToUpdate
		err = rows.Scan(&guess.Id, &guess.UserId, &guess.Points)
		if err != nil {
			return nil, err
		}

		log.Println("Guess:")
		log.Println(guess)

		guesses = append(guesses, &guess)
	}

	return guesses, nil
}
