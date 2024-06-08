package lib

import "github.com/victoroliveirab/go-htmx-soccer-guesser/models"

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
		return "DiffPlusOpposite"
	case DiffPlusWinner:
		return "DiffPlusWinner"
	case WinnerPlusWinnerGoals:
		return "WinnerPlusWinnerGoals"
	case WinnerPlusLoserGoals:
		return "WinnerPlusLoserGoals"
	case Winner:
		return "Winner"
	case Draw:
		return "Draw"
	case OneGoalButDraw:
		return "OneGoalButDraw"
	default:
		return "None"
	}
}

var DefaultPointsMap = map[Outcome]int{
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

func DefineOutcome(guess *models.Guess, fixture *models.Fixture) Outcome {
	guessedHg := guess.HomeGoals
	guessedAg := guess.AwayGoals
	actualHg := fixture.HomeScore
	actualAg := fixture.AwayScore

	if guessedHg == actualHg && guessedAg == actualAg {
		return Perfect
	}

	if guessedHg == actualAg && guessedAg == actualHg {
		return Opposite
	}

	hasHomeWon := actualHg > actualAg
	hasTeamsDrawn := actualHg == actualAg
	hasAwayWon := actualHg < actualAg

	guessedDiff := guessedHg - guessedAg
	actualDiff := actualHg - actualAg

	if guessedDiff == actualDiff && !hasTeamsDrawn {
		return DiffPlusWinner
	}

	if guessedDiff*-1 == actualDiff && !hasTeamsDrawn {
		return DiffPlusOpposite
	}

	hasHomeWonOnGuess := guessedHg > guessedAg
	hasDrawnOnGuess := guessedHg == guessedAg
	hasAwayWonOnGuess := guessedHg < guessedAg

	if hasHomeWon && hasHomeWonOnGuess && guessedHg == actualHg {
		return WinnerPlusWinnerGoals
	}

	if hasAwayWon && hasAwayWonOnGuess && guessedAg == actualAg {
		return WinnerPlusWinnerGoals
	}

	if hasHomeWon && hasHomeWonOnGuess && guessedAg == actualAg {
		return WinnerPlusLoserGoals
	}

	if hasAwayWon && hasAwayWonOnGuess && guessedHg == actualHg {
		return WinnerPlusLoserGoals
	}

	if hasHomeWon && hasHomeWonOnGuess {
		return Winner
	}

	if hasAwayWon && hasAwayWonOnGuess {
		return Winner
	}

	if hasTeamsDrawn && hasDrawnOnGuess {
		return Draw
	}

	if hasTeamsDrawn && guessedHg == actualHg {
		return OneGoalButDraw
	}

	if hasTeamsDrawn && guessedAg == actualAg {
		return OneGoalButDraw
	}

	return None
}
