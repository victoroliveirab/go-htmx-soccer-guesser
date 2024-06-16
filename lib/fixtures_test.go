package lib_test

import (
	"fmt"
	"testing"

	c "github.com/victoroliveirab/go-htmx-soccer-guesser/constants"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/lib"
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

	outcomes := [...]c.Outcome{
		c.None, c.None, c.None, c.DiffPlusOpposite, c.None, c.None,
		c.Winner, c.None, c.None, c.None, c.Opposite, c.None,
		c.Winner, c.WinnerPlusLoserGoals, c.None, c.None, c.None, c.DiffPlusOpposite,
		c.DiffPlusWinner, c.WinnerPlusLoserGoals, c.Winner, c.None, c.None, c.None,
		c.WinnerPlusWinnerGoals, c.Perfect, c.WinnerPlusWinnerGoals, c.WinnerPlusWinnerGoals, c.None, c.None,
		c.Winner, c.WinnerPlusLoserGoals, c.DiffPlusWinner, c.Winner, c.Winner, c.None,
	}

	for index, possibility := range possibilities {
		guess := models.Guess{
			HomeGoals: possibility[0],
			AwayGoals: possibility[1],
		}
		t.Run(fmt.Sprintf("Both Teams Scored - 4x1 - Guess %dx%d", guess.HomeGoals, guess.AwayGoals), func(t *testing.T) {
			expectedOutcome := outcomes[index]
			actualOutcome := lib.DefineOutcome(&guess, &fixture)
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

	outcomes := [...]c.Outcome{
		c.None, c.WinnerPlusLoserGoals, c.Perfect, c.WinnerPlusLoserGoals, c.WinnerPlusLoserGoals, c.WinnerPlusLoserGoals,
		c.None, c.None, c.WinnerPlusWinnerGoals, c.DiffPlusWinner, c.Winner, c.Winner,
		c.Opposite, c.None, c.None, c.Winner, c.DiffPlusWinner, c.Winner,
		c.None, c.DiffPlusOpposite, c.None, c.None, c.Winner, c.DiffPlusWinner,
		c.None, c.None, c.DiffPlusOpposite, c.None, c.None, c.Winner,
		c.None, c.None, c.None, c.DiffPlusOpposite, c.None, c.None,
	}

	for index, possibility := range possibilities {
		guess := models.Guess{
			HomeGoals: possibility[0],
			AwayGoals: possibility[1],
		}
		t.Run(fmt.Sprintf("Only Winner Scored - 0x2 - Guess %dx%d", guess.HomeGoals, guess.AwayGoals), func(t *testing.T) {
			expectedOutcome := outcomes[index]
			actualOutcome := lib.DefineOutcome(&guess, &fixture)
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

	outcomes := [...]c.Outcome{
		c.Draw, c.None, c.None, c.OneGoalButDraw, c.None, c.None,
		c.None, c.Draw, c.None, c.OneGoalButDraw, c.None, c.None,
		c.None, c.None, c.Draw, c.OneGoalButDraw, c.None, c.None,
		c.OneGoalButDraw, c.OneGoalButDraw, c.OneGoalButDraw, c.Perfect, c.OneGoalButDraw, c.OneGoalButDraw,
		c.None, c.None, c.None, c.OneGoalButDraw, c.Draw, c.None,
		c.None, c.None, c.None, c.OneGoalButDraw, c.None, c.Draw,
	}

	for index, possibility := range possibilities {
		guess := models.Guess{
			HomeGoals: possibility[0],
			AwayGoals: possibility[1],
		}
		t.Run(fmt.Sprintf("Draw - 3x3 - Guess %dx%d", guess.HomeGoals, guess.AwayGoals), func(t *testing.T) {
			expectedOutcome := outcomes[index]
			actualOutcome := lib.DefineOutcome(&guess, &fixture)
			if expectedOutcome != actualOutcome {
				t.Errorf("Expected outcome %s, actually got %s", expectedOutcome, actualOutcome)
			}
		})
	}
}
