package reward

import (
	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRewardRange(t *testing.T) {
	var left, right float32

	left, right = Range(2)
	assert.Equal(t, float32(0), left)
	assert.Equal(t, float32(400), right)

	left, right = Range(3)
	assert.Equal(t, float32(0), left)
	assert.Equal(t, float32(800), right)

	left, right = Range(4)
	assert.Equal(t, float32(0), left)
	assert.Equal(t, float32(1200), right)
}

func TestCalculateRewardEmptyGame(t *testing.T) {
	for _, count := range []int{ 2, 3, 4 } {
		for _, color := range model.PlayerColors.Members()[0:count] {
			game, _ := model.NewGame(count)
			view, _ := game.CreatePlayerView(color)
			assert.Equal(t, float32(0.0), Calculate(view)) // score is always zero if all pawns are in start
		}
	}
}

func TestCalculateRewardEquivalentState(t *testing.T) {
	game, _ := model.NewGame(4)
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(4)
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSquare(34)
	_ = game.Players()[model.Green].Pawns()[0].Position().MoveToSquare(49)
	_ = game.Players()[model.Blue].Pawns()[0].Position().MoveToSquare(19)
	for _, color := range model.PlayerColors.Members() {
		view, _ := game.CreatePlayerView(color)
		assert.Equal(t, float32(0.0), Calculate(view)) // score is always zero if all players are equivalent
	}
}

func TestCalculateRewardSafeZone(t *testing.T) {
	game, _ := model.NewGame(4)
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSafe(4) // last safe square before home
	view, _ := game.CreatePlayerView(model.Red)
	assert.Equal(t, float32(222), Calculate(view))
	for _, color := range []model.PlayerColor { model.Yellow, model.Green, model.Blue } {
		view, _ = game.CreatePlayerView(color)
		assert.Equal(t, float32(0), Calculate(view)) // score is always zero if all pawns are in start
	}
}

func TestCalculateRewardWinner(t *testing.T) {
	game2, _ := model.NewGame(2)
	_ = game2.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_ = game2.Players()[model.Red].Pawns()[1].Position().MoveToHome()
	_ = game2.Players()[model.Red].Pawns()[2].Position().MoveToHome()
	_ = game2.Players()[model.Red].Pawns()[3].Position().MoveToHome()
	view2, _ := game2.CreatePlayerView(model.Red)
	assert.Equal(t, float32(400), Calculate(view2))
	for _, color := range []model.PlayerColor { model.Yellow } {
		view2, _ = game2.CreatePlayerView(color)
		assert.Equal(t, float32(0), Calculate(view2)) // score is always zero if all pawns are in start
	}

	game3, _ := model.NewGame(3)
	_ = game3.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_ = game3.Players()[model.Red].Pawns()[1].Position().MoveToHome()
	_ = game3.Players()[model.Red].Pawns()[2].Position().MoveToHome()
	_ = game3.Players()[model.Red].Pawns()[3].Position().MoveToHome()
	view3, _ := game3.CreatePlayerView(model.Red)
	assert.Equal(t, float32(800), Calculate(view3))
	for _, color := range []model.PlayerColor { model.Yellow, model.Green } {
		view3, _ = game3.CreatePlayerView(color)
		assert.Equal(t, float32(0), Calculate(view3)) // score is always zero if all pawns are in start
	}

	game4, _ := model.NewGame(4)
	_ = game4.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_ = game4.Players()[model.Red].Pawns()[1].Position().MoveToHome()
	_ = game4.Players()[model.Red].Pawns()[2].Position().MoveToHome()
	_ = game4.Players()[model.Red].Pawns()[3].Position().MoveToHome()
	view4, _ := game4.CreatePlayerView(model.Red)
	assert.Equal(t, float32(1200), Calculate(view4))
	for _, color := range []model.PlayerColor { model.Yellow, model.Green, model.Blue } {
		view4, _ = game4.CreatePlayerView(color)
		assert.Equal(t, float32(0), Calculate(view4)) // score is always zero if all pawns are in start
	}
}

func TestCalculateRewardArbitrary(t *testing.T) {
	game, _ := model.NewGame(4)

	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSafe(0)
	_ = game.Players()[model.Red].Pawns()[2].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[3].Position().MoveToSquare(10)

	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSquare(34)
	_ = game.Players()[model.Yellow].Pawns()[1].Position().MoveToSquare(32)
	_ = game.Players()[model.Yellow].Pawns()[2].Position().MoveToStart()
	_ = game.Players()[model.Yellow].Pawns()[3].Position().MoveToHome()

	_ = game.Players()[model.Green].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToStart()
	_ = game.Players()[model.Green].Pawns()[2].Position().MoveToSquare(59)
	_ = game.Players()[model.Green].Pawns()[3].Position().MoveToStart()

	_ = game.Players()[model.Blue].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Blue].Pawns()[1].Position().MoveToStart()
	_ = game.Players()[model.Blue].Pawns()[2].Position().MoveToStart()
	_ = game.Players()[model.Blue].Pawns()[3].Position().MoveToStart()

	red, _ := game.CreatePlayerView(model.Red)
	yellow, _ := game.CreatePlayerView(model.Yellow)
	green, _ := game.CreatePlayerView(model.Green)
	blue, _ := game.CreatePlayerView(model.Blue)

	assert.Equal(t, float32(319), Calculate(red))
	assert.Equal(t, float32(239), Calculate(yellow))
	assert.Equal(t, float32(0), Calculate(green))
	assert.Equal(t, float32(0), Calculate(blue))
}