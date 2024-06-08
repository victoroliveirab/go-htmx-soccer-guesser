package constants

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

func (outcome Outcome) String() string {
	switch outcome {
	case Perfect:
		return "Perfect"
	case Opposite:
		return "Opposite"
	case DiffPlusOpposite:
		return "Difference + Opposite"
	case DiffPlusWinner:
		return "Difference + Winner"
	case WinnerPlusWinnerGoals:
		return "Winner + Winner Goals"
	case WinnerPlusLoserGoals:
		return "Winner + Loser Goals"
	case Winner:
		return "Winner"
	case Draw:
		return "Draw"
	case OneGoalButDraw:
		return "One-side Goal, but Draw"
	default:
		return "None"
	}
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
