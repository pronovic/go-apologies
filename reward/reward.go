package reward

import (
	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/pkg"
)

func Calculate(view model.PlayerView) float32 {
	return float32(calculateReward(view))
}

func Range(players int) (float32, float32) {
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
		distance += pkg.DistanceToHome(pawn)
	}
	return 260 - distance  // 260 = 4*65, max distance for 4 pawns
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