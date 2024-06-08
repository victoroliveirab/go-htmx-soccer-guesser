package lib

import (
	"fmt"
	"testing"

	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

var possibilities = [...][2]int{
	{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5},
	{1, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5},
	{2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5},
	{3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4}, {3, 5},
	{4, 0}, {4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 5},
	{5, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5},
}

func TestDefineOutcomeBothTeamsScored(t *testing.T) {
	fixture := models.Fixture{
		Id: 1,
		League: models.League{
			Id: 1,
		},
		Season: 2024,
		HomeTeam: models.Team{
			Id: 1,
		},
		AwayTeam: models.Team{
			Id: 2,
		},
		TimestampNumb: 1,
		Status:        "FIN",
		Referee:       "Victor Oliveira",
		HomeScore:     4,
		AwayScore:     1,
		Winner:        "Home",
		Round:         "Regular Season - 7",
	}

	outcomes := [...]Outcome{
		None, None, None, DiffPlusOpposite, None, None,
		Winner, None, None, None, Opposite, None,
		Winner, WinnerPlusLoserGoals, None, None, None, DiffPlusOpposite,
		DiffPlusWinner, WinnerPlusLoserGoals, Winner, None, None, None,
		WinnerPlusWinnerGoals, Perfect, WinnerPlusWinnerGoals, WinnerPlusWinnerGoals, None, None,
		Winner, WinnerPlusLoserGoals, DiffPlusWinner, Winner, Winner, None,
	}

	for index, possibility := range possibilities {
		guess := models.Guess{
			HomeGoals: possibility[0],
			AwayGoals: possibility[1],
		}
		t.Run(fmt.Sprintf("Both Teams Scored - 4x1 - Guess %dx%d", guess.HomeGoals, guess.AwayGoals), func(t *testing.T) {
			expectedOutcome := outcomes[index]
			actualOutcome := DefineOutcome(&guess, &fixture)
			if expectedOutcome != actualOutcome {
				t.Errorf("Expected outcome %s, actually got %s", expectedOutcome, actualOutcome)
			}
		})
	}
}

func TestDefineOutcomeOnlyWinnerScored(t *testing.T) {
	fixture := models.Fixture{
		Id: 1,
		League: models.League{
			Id: 1,
		},
		Season: 2024,
		HomeTeam: models.Team{
			Id: 1,
		},
		AwayTeam: models.Team{
			Id: 2,
		},
		TimestampNumb: 1,
		Status:        "FIN",
		Referee:       "Victor Oliveira",
		HomeScore:     0,
		AwayScore:     2,
		Winner:        "Home",
		Round:         "Regular Season - 7",
	}

	outcomes := [...]Outcome{
		None, WinnerPlusLoserGoals, Perfect, WinnerPlusLoserGoals, WinnerPlusLoserGoals, WinnerPlusLoserGoals,
		None, None, WinnerPlusWinnerGoals, DiffPlusWinner, Winner, Winner,
		Opposite, None, None, Winner, DiffPlusWinner, Winner,
		None, DiffPlusOpposite, None, None, Winner, DiffPlusWinner,
		None, None, DiffPlusOpposite, None, None, Winner,
		None, None, None, DiffPlusOpposite, None, None,
	}

	for index, possibility := range possibilities {
		guess := models.Guess{
			HomeGoals: possibility[0],
			AwayGoals: possibility[1],
		}
		t.Run(fmt.Sprintf("Only Winner Scored - 0x2 - Guess %dx%d", guess.HomeGoals, guess.AwayGoals), func(t *testing.T) {
			expectedOutcome := outcomes[index]
			actualOutcome := DefineOutcome(&guess, &fixture)
			if expectedOutcome != actualOutcome {
				t.Errorf("Expected outcome %s, actually got %s", expectedOutcome, actualOutcome)
			}
		})
	}
}

func TestDefineOutcomeDraw(t *testing.T) {
	fixture := models.Fixture{
		Id: 1,
		League: models.League{
			Id: 1,
		},
		Season: 2024,
		HomeTeam: models.Team{
			Id: 1,
		},
		AwayTeam: models.Team{
			Id: 2,
		},
		TimestampNumb: 1,
		Status:        "FIN",
		Referee:       "Victor Oliveira",
		HomeScore:     3,
		AwayScore:     3,
		Winner:        "Home",
		Round:         "Regular Season - 7",
	}

	outcomes := [...]Outcome{
		Draw, None, None, OneGoalButDraw, None, None,
		None, Draw, None, OneGoalButDraw, None, None,
		None, None, Draw, OneGoalButDraw, None, None,
		OneGoalButDraw, OneGoalButDraw, OneGoalButDraw, Perfect, OneGoalButDraw, OneGoalButDraw,
		None, None, None, OneGoalButDraw, Draw, None,
		None, None, None, OneGoalButDraw, None, Draw,
	}

	for index, possibility := range possibilities {
		guess := models.Guess{
			HomeGoals: possibility[0],
			AwayGoals: possibility[1],
		}
		t.Run(fmt.Sprintf("Draw - 3x3 - Guess %dx%d", guess.HomeGoals, guess.AwayGoals), func(t *testing.T) {
			expectedOutcome := outcomes[index]
			actualOutcome := DefineOutcome(&guess, &fixture)
			if expectedOutcome != actualOutcome {
				t.Errorf("Expected outcome %s, actually got %s", expectedOutcome, actualOutcome)
			}
		})
	}
}
