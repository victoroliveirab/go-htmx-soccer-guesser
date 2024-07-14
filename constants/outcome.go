package constants

import "database/sql"

type Outcome int

const (
	Perfect Outcome = iota
	Opposite
	DiffPlusWinner
	DiffPlusOpposite
	WinnerPlusWinnerGoals
	WinnerPlusLoserGoals
	Winner
	Draw
	OneGoalButDraw
	None
)

var OutcomesMap = map[string]Outcome{
	"Perfect":               Perfect,
	"Opposite":              Opposite,
	"DiffPlusWinner":        DiffPlusWinner,
	"DiffPlusOpposite":      DiffPlusOpposite,
	"WinnerPlusWinnerGoals": WinnerPlusWinnerGoals,
	"WinnerPlusLoserGoals":  WinnerPlusLoserGoals,
	"Winner":                Winner,
	"Draw":                  Draw,
	"OneGoalButDraw":        OneGoalButDraw,
	"None":                  None,
}

var OutcomesReverseMap = map[Outcome]string{
	Perfect:               "Perfect",
	Opposite:              "Opposite",
	DiffPlusWinner:        "DiffPlusWinner",
	DiffPlusOpposite:      "DiffPlusOpposite",
	WinnerPlusWinnerGoals: "WinnerPlusWinnerGoals",
	WinnerPlusLoserGoals:  "WinnerPlusLoserGoals",
	Winner:                "Winner",
	Draw:                  "Draw",
	OneGoalButDraw:        "OneGoalButDraw",
	None:                  "None",
}

func (outcome Outcome) String() string {
	str, exists := OutcomesReverseMap[outcome]
	if !exists {
		return None.String()
	}
	return str
}

func NormalizeOutcome(outcome sql.NullInt64) string {
	if !outcome.Valid {
		return "N/A"
	}
	return Outcome(outcome.Int64).String()
}

var DefaultOutcomePointsMap = map[Outcome]int{
	Perfect:               20,
	Opposite:              -10,
	DiffPlusWinner:        15,
	DiffPlusOpposite:      -5,
	WinnerPlusWinnerGoals: 12,
	WinnerPlusLoserGoals:  11,
	Winner:                10,
	Draw:                  15,
	OneGoalButDraw:        4,
	None:                  0,
}
