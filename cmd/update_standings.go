package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/infra"
)

type MatchesStatusGoals struct {
	For     int `json:"for"`
	Against int `json:"against"`
}

type MatchesStats struct {
	Played int                `json:"played"`
	Win    int                `json:"win"`
	Draw   int                `json:"draw"`
	Lose   int                `json:"lose"`
	Goals  MatchesStatusGoals `json:"goals"`
}

type StandingEntryTeam struct {
	ApiFootballId int `json:"id"`
}

type StandingEntry struct {
	Rank        int               `json:"rank"`
	Team        StandingEntryTeam `json:"team"`
	Points      int               `json:"points"`
	GoalsDiff   int               `json:"goalsDiff"`
	Group       string            `json:"group"`
	Form        string            `json:"form"`
	Status      string            `json:"status"`
	Description string            `json:"description"`
	All         MatchesStats      `json:"all"`
	Home        MatchesStats      `json:"home"`
	Away        MatchesStats      `json:"away"`
}

var apiFootballTeamsIdsMap map[int]int

func (s StandingEntry) MarshalJSON() ([]byte, error) {
	type Alias StandingEntry
	return json.Marshal(&struct {
		Alias
		TeamID int `json:"teamId"`
	}{
		Alias:  (Alias)(s),
		TeamID: apiFootballTeamsIdsMap[s.Team.ApiFootballId],
	})
}

type StandingResponseEntry struct {
	ApiFootballId int               `json:"id"`
	Season        int               `json:"season"`
	Standings     [][]StandingEntry `json:"standings"`
}

type StandingResponse struct {
	League StandingResponseEntry `json:"league"`
}

type StandingJson struct {
	Response []StandingResponse `json:"response"`
}

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

	jsonFiles := getJsons()

	for _, fileStruct := range jsonFiles {
		path := [2]string{"json", fileStruct.Filename}
		file := strings.Join(path[:], "/")
		updateStandingsFromFile(file)
	}
}

type fileWithTimestamp struct {
	Filename  string
	Timestamp int64
}

func getJsons() map[int64]fileWithTimestamp {
	regex, err := regexp.Compile(`^standings-\d+-(\d+)\.json$`)
	if err != nil {
		log.Fatalf("failed to compile regex: %v", err)
		panic(err)
	}

	files, err := os.ReadDir("json")
	if err != nil {
		log.Fatalf("failed to read json directory: %v", err)
		panic(err)
	}

	var jsonFiles []fileWithTimestamp

	for _, file := range files {
		if !file.IsDir() {
			matches := regex.FindStringSubmatch(file.Name())
			if len(matches) > 1 {
				timestamp, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					log.Printf("Failed to parse timestamp: %v", err)
					continue
				}
				jsonFiles = append(jsonFiles, fileWithTimestamp{
					Filename:  file.Name(),
					Timestamp: timestamp,
				})
			}
		}
	}

	var filteredFiles map[int64]fileWithTimestamp
	filteredFiles = make(map[int64]fileWithTimestamp)

	for _, file := range jsonFiles {
		filename := file.Filename
		filenameTokens := strings.Split(filename, "-")
		apiFootballLeagueId, err := strconv.ParseInt(filenameTokens[1], 10, 64)
		if err != nil {
			log.Printf("wrong filename structure: %s", filename)
			continue
		}
		timestamp, err := strconv.ParseInt(strings.Split(filenameTokens[2], ".")[0], 10, 64)
		if err != nil {
			log.Printf("invalid file: %s", filename)
			continue
		}

		entry, exists := filteredFiles[apiFootballLeagueId]

		if !exists {
			filteredFiles[apiFootballLeagueId] = fileWithTimestamp{
				Filename:  filename,
				Timestamp: timestamp,
			}
		} else if entry.Timestamp < timestamp {
			filteredFiles[apiFootballLeagueId] = fileWithTimestamp{
				Filename:  filename,
				Timestamp: timestamp,
			}
		}
	}

	return filteredFiles
}

func updateStandingsFromFile(file string) {
	db := infra.Db
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
		panic(err)
	}

	var jsonContent StandingJson
	err = json.Unmarshal(byteValue, &jsonContent)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %v", err)
		panic(err)
	}

	if len(jsonContent.Response) == 0 {
		log.Printf("file %s does not contain a response", file)
		return
	}

	data := jsonContent.Response[0]

	var league_season_id int
	row := db.QueryRow(`
        SELECT league_season.id
        FROM Leagues league
        JOIN Leagues_Seasons league_season ON league_season.league_id = league.id
        WHERE league.api_football_id = ? AND season = ?
    `, data.League.ApiFootballId, data.League.Season)
	err = row.Scan(&league_season_id)
	if err != nil {
		log.Fatalf("failed scanning league_season_id %d row: %v", league_season_id, err)
		panic(err)
	}

	standings := data.League.Standings[0]
	apiFootballTeamsIdsMap = map[int]int{}

	for _, entry := range standings {
		apiFootballTeamId := entry.Team.ApiFootballId
		row := db.QueryRow(`
            SELECT id FROM Teams WHERE api_football_id = ?
        `, apiFootballTeamId)
		var teamId int
		row.Scan(&teamId)
		apiFootballTeamsIdsMap[apiFootballTeamId] = teamId
	}

	// FIXME: this currently makes we save team.id (api_football_id) and team_id (internal id)
	// I have to figure out how to make an unmarshalled field private when marshalling
	marshalledData, err := json.Marshal(standings)

	if err != nil {
		log.Fatalf("failed to marshal data: %v", err)
		panic(err)
	}

	_, err = db.Exec("UPDATE Leagues_Seasons SET standings = ? WHERE id = ?", marshalledData, league_season_id)

	if err != nil {
		log.Fatalf("failed to update League_Seasons table: %v", err)
		panic(err)
	}

	log.Printf("updated league_season %d", league_season_id)

}
