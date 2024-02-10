package reward

// Version 1 of the reward algorithm was developed by hand following my own mental model for the
// strength of a position.
//
// Basically, a pawn is worth more the closer it is to home.  It's worth incrementally more in the
// safe zone (where it can't be hurt).  There's an additional bonus for winning, to incentivize the
// engine to pick a move that ends the game as fast as possible.  The score is calculated by
// comparing the player's position relative to the positions of all its opponents.  We want the
// engine to pick the move that both maximizes the player's position and also minimizes the
// positions of its opponents.
//
// In simulation runs in the original Python implementation, a reward-based character source vastly
// out performs a source that picks its moves randomly.  The worst-case scenario is a 4-player
// standard mode game between a single reward-based source and 3 random sources, where the
// reward-based source wins about 70% of the time.  This is probably because in a standard mode
// game, the possible moves in each turn are fairly limited, due to each player picking and playing
// the top card off the deck.  This evens the playing field, because it's quite likely that any
// player will have no good move on their turn.  In a 4-player adult mode game, where the engine has
// the opportunity to choose between more possible moves for each turn, a reward-based source wins
// more than 98% of the time against 3 random sources.

import (
	"github.com/pronovic/go-apologies/model"
)

type Calculator interface {

	// Calculate calculate the reward associated with a player view
	Calculate(view model.PlayerView) float32

	// Range Return the range of possible rewards for a game
	Range(players int) (float32, float32)
}

type calculator struct {
}

func NewCalculator() Calculator {
	return &calculator{}
}

func (c *calculator) Calculate(view model.PlayerView) float32 {
	return float32(calculateReward(view))
}

func (c *calculator) Range(players int) (float32, float32) {
	return 0.0, float32((players - 1) * 400) // reward is up to 400 points per opponent
}

func calculateReward(view model.PlayerView) int {
	// Reward measures this player's overall game position relative to their opponents
	playerScore := calculatePlayerScore(view.Player())
	opponentScore := 0
	for _, opponent := range view.Opponents() {
		opponentScore += calculatePlayerScore(opponent)
	}
	reward := (len(view.Opponents()) * playerScore) - opponentScore
	if reward < 0 {
		return 0
	} else {
		return reward
	}
}

func calculatePlayerScore(player model.Player) int {
	// There are 3 different incentives, designed to encourage the right behavior
	distanceIncentive := calculateDistanceIncentive(player)
	safeIncentive := calculateSafeIncentive(player)
	winnerIncentive := calculateWinnerIncentive(player)
	return distanceIncentive + safeIncentive + winnerIncentive
}

func calculateDistanceIncentive(player model.Player) int {
	// Incentive of 1 point for each square closer to home for each of the player's 4 pawns
	distance := 0
	for _, pawn := range player.Pawns() {
		distance += distanceToHome(pawn)
	}
	return 260 - distance // 260 = 4*65, max distance for 4 pawns
}

func calculateSafeIncentive(player model.Player) int {
	// Incentive of 10 points for each pawn in safe or home
	incentive := 0
	for _, pawn := range player.Pawns() {
		if pawn.Position().Home() || pawn.Position().Safe() != nil {
			incentive += 10
		}
	}
	return incentive
}

func calculateWinnerIncentive(player model.Player) int {
	// Incentive of 100 points for winning the game
	if player.AllPawnsInHome() {
		return 100
	} else {
		return 0
	}
}

// distanceToHome calculates the distance to home for this pawn, a number of squares when moving forward.
func distanceToHome(pawn model.Pawn) int {
	if pawn.Position().Home() {
		return 0
	} else if pawn.Position().Start() {
		return 65
	} else if pawn.Position().Safe() != nil {
		return model.SafeSquares - *pawn.Position().Safe()
	} else {
		circle := *model.StartCircles[pawn.Color()].Square()
		turn := *model.TurnSquares[pawn.Color()].Square()
		square := *pawn.Position().Square()
		squareToCorner := model.BoardSquares - square
		cornerToTurn := turn
		turnToHome := model.SafeSquares + 1
		total := squareToCorner + cornerToTurn + turnToHome
		if turn < square && square < circle {
			return total
		} else {
			if total < 65 {
				return total
			} else {
				return total - 60
			}
		}
	}
}
