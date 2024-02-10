package reward

import (
	"testing"

	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
)

func TestRewardRange(t *testing.T) {
	var left, right float32
	calc := NewCalculator()

	left, right = calc.Range(2)
	assert.Equal(t, float32(0), left)
	assert.Equal(t, float32(400), right)

	left, right = calc.Range(3)
	assert.Equal(t, float32(0), left)
	assert.Equal(t, float32(800), right)

	left, right = calc.Range(4)
	assert.Equal(t, float32(0), left)
	assert.Equal(t, float32(1200), right)
}

func TestCalculateRewardEmptyGame(t *testing.T) {
	calc := NewCalculator()
	for _, count := range []int{2, 3, 4} {
		for _, color := range model.PlayerColors.Members()[0:count] {
			game, _ := model.NewGame(count, nil)
			view, _ := game.CreatePlayerView(color)
			assert.Equal(t, float32(0.0), calc.Calculate(view)) // score is always zero if all pawns are in start
		}
	}
}

func TestCalculateRewardEquivalentState(t *testing.T) {
	calc := NewCalculator()
	game, _ := model.NewGame(4, nil)
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(4)
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSquare(34)
	_ = game.Players()[model.Green].Pawns()[0].Position().MoveToSquare(49)
	_ = game.Players()[model.Blue].Pawns()[0].Position().MoveToSquare(19)
	for _, color := range model.PlayerColors.Members() {
		view, _ := game.CreatePlayerView(color)
		assert.Equal(t, float32(0.0), calc.Calculate(view)) // score is always zero if all players are equivalent
	}
}

func TestCalculateRewardSafeZone(t *testing.T) {
	calc := NewCalculator()
	game, _ := model.NewGame(4, nil)
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSafe(4) // last safe square before home
	view, _ := game.CreatePlayerView(model.Red)
	assert.Equal(t, float32(222), calc.Calculate(view))
	for _, color := range []model.PlayerColor{model.Yellow, model.Green, model.Blue} {
		view, _ = game.CreatePlayerView(color)
		assert.Equal(t, float32(0), calc.Calculate(view)) // score is always zero if all pawns are in start
	}
}

func TestCalculateRewardWinner(t *testing.T) {
	calc := NewCalculator()

	game2, _ := model.NewGame(2, nil)
	_ = game2.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_ = game2.Players()[model.Red].Pawns()[1].Position().MoveToHome()
	_ = game2.Players()[model.Red].Pawns()[2].Position().MoveToHome()
	_ = game2.Players()[model.Red].Pawns()[3].Position().MoveToHome()
	view2, _ := game2.CreatePlayerView(model.Red)
	assert.Equal(t, float32(400), calc.Calculate(view2))
	for _, color := range []model.PlayerColor{model.Yellow} {
		view2, _ = game2.CreatePlayerView(color)
		assert.Equal(t, float32(0), calc.Calculate(view2)) // score is always zero if all pawns are in start
	}

	game3, _ := model.NewGame(3, nil)
	_ = game3.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_ = game3.Players()[model.Red].Pawns()[1].Position().MoveToHome()
	_ = game3.Players()[model.Red].Pawns()[2].Position().MoveToHome()
	_ = game3.Players()[model.Red].Pawns()[3].Position().MoveToHome()
	view3, _ := game3.CreatePlayerView(model.Red)
	assert.Equal(t, float32(800), calc.Calculate(view3))
	for _, color := range []model.PlayerColor{model.Yellow, model.Green} {
		view3, _ = game3.CreatePlayerView(color)
		assert.Equal(t, float32(0), calc.Calculate(view3)) // score is always zero if all pawns are in start
	}

	game4, _ := model.NewGame(4, nil)
	_ = game4.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_ = game4.Players()[model.Red].Pawns()[1].Position().MoveToHome()
	_ = game4.Players()[model.Red].Pawns()[2].Position().MoveToHome()
	_ = game4.Players()[model.Red].Pawns()[3].Position().MoveToHome()
	view4, _ := game4.CreatePlayerView(model.Red)
	assert.Equal(t, float32(1200), calc.Calculate(view4))
	for _, color := range []model.PlayerColor{model.Yellow, model.Green, model.Blue} {
		view4, _ = game4.CreatePlayerView(color)
		assert.Equal(t, float32(0), calc.Calculate(view4)) // score is always zero if all pawns are in start
	}
}

func TestCalculateRewardArbitrary(t *testing.T) {
	calc := NewCalculator()
	game, _ := model.NewGame(4, nil)

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

	assert.Equal(t, float32(319), calc.Calculate(red))
	assert.Equal(t, float32(239), calc.Calculate(yellow))
	assert.Equal(t, float32(0), calc.Calculate(green))
	assert.Equal(t, float32(0), calc.Calculate(blue))
}

func TestDistanceToHome(t *testing.T) {
	// distance from home is always 0
	for _, color := range []model.PlayerColor{model.Red, model.Yellow, model.Green} {
		assert.Equal(t, 0, distanceToHome(pawnHome(color)))
	}

	// distance from start is always 65
	for _, color := range []model.PlayerColor{model.Red, model.Yellow, model.Green} {
		assert.Equal(t, 65, distanceToHome(pawnStart(color)))
	}

	// distance from within safe is always <= 5
	assert.Equal(t, 5, distanceToHome(pawnSafe(model.Red, 0)))
	assert.Equal(t, 4, distanceToHome(pawnSafe(model.Red, 1)))
	assert.Equal(t, 3, distanceToHome(pawnSafe(model.Red, 2)))
	assert.Equal(t, 2, distanceToHome(pawnSafe(model.Red, 3)))
	assert.Equal(t, 1, distanceToHome(pawnSafe(model.Red, 4)))

	// distance from circle is always 64
	assert.Equal(t, 64, distanceToHome(pawnSquare(model.Red, 4)))
	assert.Equal(t, 64, distanceToHome(pawnSquare(model.Blue, 19)))
	assert.Equal(t, 64, distanceToHome(pawnSquare(model.Yellow, 34)))
	assert.Equal(t, 64, distanceToHome(pawnSquare(model.Green, 49)))

	// distance from square between turn and circle is always 65
	assert.Equal(t, 65, distanceToHome(pawnSquare(model.Red, 3)))
	assert.Equal(t, 65, distanceToHome(pawnSquare(model.Blue, 18)))
	assert.Equal(t, 65, distanceToHome(pawnSquare(model.Yellow, 33)))
	assert.Equal(t, 65, distanceToHome(pawnSquare(model.Green, 48)))

	// distance from turn is always 6
	assert.Equal(t, 6, distanceToHome(pawnSquare(model.Red, 2)))
	assert.Equal(t, 6, distanceToHome(pawnSquare(model.Blue, 17)))
	assert.Equal(t, 6, distanceToHome(pawnSquare(model.Yellow, 32)))
	assert.Equal(t, 6, distanceToHome(pawnSquare(model.Green, 47)))

	// check some arbitrary squares
	assert.Equal(t, 7, distanceToHome(pawnSquare(model.Red, 1)))
	assert.Equal(t, 8, distanceToHome(pawnSquare(model.Red, 0)))
	assert.Equal(t, 9, distanceToHome(pawnSquare(model.Red, 59)))
	assert.Equal(t, 59, distanceToHome(pawnSquare(model.Red, 9)))
	assert.Equal(t, 23, distanceToHome(pawnSquare(model.Blue, 0)))
	assert.Equal(t, 13, distanceToHome(pawnSquare(model.Green, 40)))
}

func pawnHome(color model.PlayerColor) model.Pawn {
	pawn := model.NewPawn(color, 0)
	pawn.SetPosition(positionHome())
	return pawn
}

func pawnStart(color model.PlayerColor) model.Pawn {
	pawn := model.NewPawn(color, 0)
	pawn.SetPosition(positionStart())
	return pawn
}

func pawnSafe(color model.PlayerColor, safe int) model.Pawn {
	pawn := model.NewPawn(color, 0)
	pawn.SetPosition(positionSafe(safe))
	return pawn
}

func pawnSquare(color model.PlayerColor, square int) model.Pawn {
	pawn := model.NewPawn(color, 0)
	pawn.SetPosition(positionSquare(square))
	return pawn
}

func positionHome() model.Position {
	return model.NewPosition(false, true, nil, nil)
}

func positionStart() model.Position {
	return model.NewPosition(true, false, nil, nil)
}

func positionSafe(safe int) model.Position {
	return model.NewPosition(false, false, &safe, nil)
}

func positionSquare(square int) model.Position {
	return model.NewPosition(false, false, nil, &square)
}
