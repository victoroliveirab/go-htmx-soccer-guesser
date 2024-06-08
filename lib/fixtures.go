package lib

import (
	"github.com/victoroliveirab/go-htmx-soccer-guesser/constants"
	"github.com/victoroliveirab/go-htmx-soccer-guesser/models"
)

func DefineOutcome(guess *models.Guess, fixture *models.Fixture) constants.Outcome {
	guessedHg := guess.HomeGoals
	guessedAg := guess.AwayGoals
	actualHg := fixture.HomeScore
	actualAg := fixture.AwayScore

	if guessedHg == actualHg && guessedAg == actualAg {
		return constants.Perfect
	}

	if guessedHg == actualAg && guessedAg == actualHg {
		return constants.Opposite
	}

	hasHomeWon := actualHg > actualAg
	hasTeamsDrawn := actualHg == actualAg
	hasAwayWon := actualHg < actualAg

	guessedDiff := guessedHg - guessedAg
	actualDiff := actualHg - actualAg

	if guessedDiff == actualDiff && !hasTeamsDrawn {
		return constants.DiffPlusWinner
	}

	if guessedDiff*-1 == actualDiff && !hasTeamsDrawn {
		return constants.DiffPlusOpposite
	}

	hasHomeWonOnGuess := guessedHg > guessedAg
	hasDrawnOnGuess := guessedHg == guessedAg
	hasAwayWonOnGuess := guessedHg < guessedAg

	if hasHomeWon && hasHomeWonOnGuess && guessedHg == actualHg {
		return constants.WinnerPlusWinnerGoals
	}

	if hasAwayWon && hasAwayWonOnGuess && guessedAg == actualAg {
		return constants.WinnerPlusWinnerGoals
	}

	if hasHomeWon && hasHomeWonOnGuess && guessedAg == actualAg {
		return constants.WinnerPlusLoserGoals
	}

	if hasAwayWon && hasAwayWonOnGuess && guessedHg == actualHg {
		return constants.WinnerPlusLoserGoals
	}

	if hasHomeWon && hasHomeWonOnGuess {
		return constants.Winner
	}

	if hasAwayWon && hasAwayWonOnGuess {
		return constants.Winner
	}

	if hasTeamsDrawn && hasDrawnOnGuess {
		return constants.Draw
	}

	if hasTeamsDrawn && guessedHg == actualHg {
		return constants.OneGoalButDraw
	}

	if hasTeamsDrawn && guessedAg == actualAg {
		return constants.OneGoalButDraw
	}

	return constants.None
}
