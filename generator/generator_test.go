package generator

import (
	"testing"

	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
)

var emptyMoves = make([]model.Move, 0)

func TestCalculatePositionHome(t *testing.T) {
	for _, color := range model.PlayerColors.Members() {
		calculatePositionFailure(t, color, positionHome(), 1, "pawn in home or start may not move")
	}
}

func TestCalculatePositionStart(t *testing.T) {
	for _, color := range model.PlayerColors.Members() {
		calculatePositionFailure(t, color, positionStart(), 1, "pawn in home or start may not move")
	}
}

func TestCalculatePositionFromSafe(t *testing.T) {
	var color model.PlayerColor

	for _, color = range model.PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSafe(0), 0, positionSafe(0))
		calculatePositionSuccess(t, color, positionSafe(3), 0, positionSafe(3))
	}

	for _, color = range model.PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSafe(0), 1, positionSafe(1))
		calculatePositionSuccess(t, color, positionSafe(2), 2, positionSafe(4))
		calculatePositionSuccess(t, color, positionSafe(4), 1, positionHome())
	}

	for _, color = range model.PlayerColors.Members() {
		calculatePositionFailure(t, color, positionSafe(3), 3, "pawn cannot move past home")
		calculatePositionFailure(t, color, positionSafe(4), 2, "pawn cannot move past home")
	}

	for _, color = range model.PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSafe(4), -2, positionSafe(2))
		calculatePositionSuccess(t, color, positionSafe(1), -1, positionSafe(0))
	}

	calculatePositionSuccess(t, model.Red, positionSafe(0), -1, positionSquare(2))
	calculatePositionSuccess(t, model.Red, positionSafe(0), -2, positionSquare(1))
	calculatePositionSuccess(t, model.Red, positionSafe(0), -3, positionSquare(0))
	calculatePositionSuccess(t, model.Red, positionSafe(0), -4, positionSquare(59))
	calculatePositionSuccess(t, model.Red, positionSafe(0), -5, positionSquare(58))

	calculatePositionSuccess(t, model.Blue, positionSafe(0), -1, positionSquare(17))
	calculatePositionSuccess(t, model.Blue, positionSafe(0), -2, positionSquare(16))

	calculatePositionSuccess(t, model.Yellow, positionSafe(0), -1, positionSquare(32))
	calculatePositionSuccess(t, model.Yellow, positionSafe(0), -2, positionSquare(31))

	calculatePositionSuccess(t, model.Green, positionSafe(0), -1, positionSquare(47))
	calculatePositionSuccess(t, model.Green, positionSafe(0), -2, positionSquare(46))
}

func TestCalculatePositionFromSquare(t *testing.T) {
	var color model.PlayerColor

	calculatePositionSuccess(t, model.Red, positionSquare(58), 1, positionSquare(59))
	calculatePositionSuccess(t, model.Red, positionSquare(59), 1, positionSquare(0))
	calculatePositionSuccess(t, model.Red, positionSquare(54), 5, positionSquare(59))
	calculatePositionSuccess(t, model.Red, positionSquare(54), 6, positionSquare(0))
	calculatePositionSuccess(t, model.Red, positionSquare(54), 7, positionSquare(1))

	for _, color = range model.PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSquare(54), 5, positionSquare(59))
		calculatePositionSuccess(t, color, positionSquare(54), 6, positionSquare(0))
		calculatePositionSuccess(t, color, positionSquare(54), 7, positionSquare(1))
		calculatePositionSuccess(t, color, positionSquare(58), 1, positionSquare(59))
		calculatePositionSuccess(t, color, positionSquare(59), 1, positionSquare(0))
		calculatePositionSuccess(t, color, positionSquare(0), 1, positionSquare(1))
		calculatePositionSuccess(t, color, positionSquare(1), 1, positionSquare(2))
		calculatePositionSuccess(t, color, positionSquare(10), 5, positionSquare(15))
	}

	for _, color = range model.PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSquare(59), -5, positionSquare(54))
		calculatePositionSuccess(t, color, positionSquare(0), -6, positionSquare(54))
		calculatePositionSuccess(t, color, positionSquare(1), -7, positionSquare(54))
		calculatePositionSuccess(t, color, positionSquare(59), -1, positionSquare(58))
		calculatePositionSuccess(t, color, positionSquare(0), -1, positionSquare(59))
		calculatePositionSuccess(t, color, positionSquare(1), -1, positionSquare(0))
		calculatePositionSuccess(t, color, positionSquare(2), -1, positionSquare(1))
		calculatePositionSuccess(t, color, positionSquare(15), -5, positionSquare(10))
	}

	calculatePositionSuccess(t, model.Red, positionSquare(0), 3, positionSafe(0))
	calculatePositionSuccess(t, model.Red, positionSquare(1), 2, positionSafe(0))
	calculatePositionSuccess(t, model.Red, positionSquare(2), 1, positionSafe(0))
	calculatePositionSuccess(t, model.Red, positionSquare(1), 3, positionSafe(1))
	calculatePositionSuccess(t, model.Red, positionSquare(2), 2, positionSafe(1))
	calculatePositionSuccess(t, model.Red, positionSquare(2), 6, positionHome())
	calculatePositionSuccess(t, model.Red, positionSquare(51), 12, positionSafe(0))
	calculatePositionSuccess(t, model.Red, positionSquare(52), 12, positionSafe(1))
	calculatePositionSuccess(t, model.Red, positionSquare(58), 5, positionSafe(0))
	calculatePositionSuccess(t, model.Red, positionSquare(59), 4, positionSafe(0))
	calculatePositionFailure(t, model.Red, positionSquare(2), 7, "pawn cannot move past home")

	calculatePositionSuccess(t, model.Blue, positionSquare(16), 2, positionSafe(0))
	calculatePositionSuccess(t, model.Blue, positionSquare(17), 1, positionSafe(0))
	calculatePositionSuccess(t, model.Blue, positionSquare(16), 3, positionSafe(1))
	calculatePositionSuccess(t, model.Blue, positionSquare(17), 2, positionSafe(1))
	calculatePositionSuccess(t, model.Blue, positionSquare(17), 6, positionHome())
	calculatePositionFailure(t, model.Blue, positionSquare(17), 7, "pawn cannot move past home")

	calculatePositionSuccess(t, model.Yellow, positionSquare(31), 2, positionSafe(0))
	calculatePositionSuccess(t, model.Yellow, positionSquare(32), 1, positionSafe(0))
	calculatePositionSuccess(t, model.Yellow, positionSquare(31), 3, positionSafe(1))
	calculatePositionSuccess(t, model.Yellow, positionSquare(32), 2, positionSafe(1))
	calculatePositionSuccess(t, model.Yellow, positionSquare(32), 6, positionHome())
	calculatePositionFailure(t, model.Yellow, positionSquare(32), 7, "pawn cannot move past home")

	calculatePositionSuccess(t, model.Green, positionSquare(46), 2, positionSafe(0))
	calculatePositionSuccess(t, model.Green, positionSquare(47), 1, positionSafe(0))
	calculatePositionSuccess(t, model.Green, positionSquare(46), 3, positionSafe(1))
	calculatePositionSuccess(t, model.Green, positionSquare(47), 2, positionSafe(1))
	calculatePositionSuccess(t, model.Green, positionSquare(47), 6, positionHome())
	calculatePositionFailure(t, model.Green, positionSquare(47), 7, "pawn cannot move past home")
}

func TestLegalMovesCard1(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn in start, on the board, or in safe
	game = setupGame()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn from start with no conflicts
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(square(pawn, 4)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(4)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = emptyMoves // can't start because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSquare(4)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(square(pawn, 4)), actionSlice(bump(view, model.Yellow, 0))))
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(square(pawn, 7)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(7)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = emptyMoves // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(7)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(square(pawn, 7)), actionSlice(bump(view, model.Green, 1))))
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard2(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn in start, on the board, or in safe
	game = setupGame()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn from start with no conflicts
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = moveSlice(move(card, actionSlice(square(pawn, 4)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(4)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = emptyMoves // can't start because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSquare(4)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = moveSlice(move(card, actionSlice(square(pawn, 4)), actionSlice(bump(view, model.Yellow, 0))))
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = moveSlice(move(card, actionSlice(square(pawn, 8)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(8)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = emptyMoves // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(8)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = moveSlice(move(card, actionSlice(square(pawn, 8)), actionSlice(bump(view, model.Green, 1))))
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard3(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card3)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card3)
	expected = moveSlice(move(card, actionSlice(square(pawn, 9)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(9)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card3)
	expected = emptyMoves // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(9)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card3)
	expected = moveSlice(move(card, actionSlice(square(pawn, 9)), actionSlice(bump(view, model.Green, 1))))
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard4(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card4)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card4)
	expected = moveSlice(move(card, actionSlice(square(pawn, 2)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(2)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card4)
	expected = emptyMoves
	assert.Equal(t, expected, moves) // can't move because we have a pawn there already

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(2)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card4)
	expected = moveSlice(move(card, actionSlice(square(pawn, 2)), actionSlice(bump(view, model.Green, 1))))
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard5(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card5)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card5)
	expected = moveSlice(move(card, actionSlice(square(pawn, 11)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(11)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card5)
	expected = emptyMoves // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(11)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card5)
	expected = moveSlice(move(card, actionSlice(square(pawn, 11)), actionSlice(bump(view, model.Green, 1))))
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard7(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var other model.Pawn
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card7)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// One move available if there is one pawn on the board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card7)
	expected = moveSlice(move(card, actionSlice(square(pawn, 13)), nil))
	assert.Equal(t, expected, moves)

	// Multiple moves available if there is more than one pawn on the board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[2].Position().MoveToSquare(55)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card7)
	other = view.Player().Pawns()[2]
	expected = moveSlice(
		move(card, actionSlice(square(pawn, 13)), nil),                    // move our pawn 7
		move(card, actionSlice(square(pawn, 7), square(other, 1)), nil),   // split (1, 6)
		move(card, actionSlice(square(pawn, 8), square(other, 0)), nil),   // split (2, 5)
		move(card, actionSlice(square(pawn, 9), square(other, 59)), nil),  // split (3, 4)
		move(card, actionSlice(square(pawn, 10), square(other, 58)), nil), // split (4, 3)
		move(card, actionSlice(square(pawn, 11), square(other, 57)), nil), // split (5, 2)
		move(card, actionSlice(square(pawn, 12), square(other, 56)), nil), // split (6, 1)
	)
	assert.Equal(t, expected, moves)

	// Either half of a move might bump an opponent back to start
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[2].Position().MoveToSquare(55)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(10)
	_ = game.Players()[model.Blue].Pawns()[3].Position().MoveToSquare(56)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card7)
	other = view.Player().Pawns()[2]
	expected = moveSlice(
		move(card, actionSlice(square(pawn, 13)), nil),                                                        // move our pawn 7
		move(card, actionSlice(square(pawn, 7), square(other, 1)), nil),                                       // split (1, 6)
		move(card, actionSlice(square(pawn, 8), square(other, 0)), nil),                                       // split (2, 5)
		move(card, actionSlice(square(pawn, 9), square(other, 59)), nil),                                      // split (3, 4)
		move(card, actionSlice(square(pawn, 10), square(other, 58)), actionSlice(bump(view, model.Green, 1))), // split (4, 3)
		move(card, actionSlice(square(pawn, 11), square(other, 57)), nil),                                     // split (5, 2)
		move(card, actionSlice(square(pawn, 12), square(other, 56)), actionSlice(bump(view, model.Blue, 3))),  // split (6, 1)
	)
	assert.Equal(t, expected, moves)

	// If either half of the move has a conflict with another pawn of the same color, the entire move is invalidated
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(9)
	_ = game.Players()[model.Red].Pawns()[2].Position().MoveToSquare(55)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card7)
	other1 := view.Player().Pawns()[1]
	other2 := view.Player().Pawns()[2]
	expected = moveSlice(
		move(card, actionSlice(square(pawn, 13)), nil),                     // move our pawn 7
		move(card, actionSlice(square(pawn, 7), square(other1, 15)), nil),  // split (1, 6)
		move(card, actionSlice(square(pawn, 8), square(other1, 14)), nil),  // split (2, 5)
		move(card, actionSlice(square(pawn, 9), square(other1, 13)), nil),  // split (3, 4)
		move(card, actionSlice(square(pawn, 10), square(other1, 12)), nil), // split (4, 3)
		move(card, actionSlice(square(pawn, 11), square(other1, 11)), nil), // split (5, 2)
		move(card, actionSlice(square(pawn, 12), square(other1, 10)), nil), // split (6, 1)
		move(card, actionSlice(square(pawn, 7), square(other2, 1)), nil),   // split (1, 6)
		move(card, actionSlice(square(pawn, 8), square(other2, 0)), nil),   // split (2, 5)
		// the move for square 9 is disallowed because pawn[1] already lives there, and isn't part of this action
		move(card, actionSlice(square(pawn, 10), square(other2, 58)), nil), // split (4, 3)
		move(card, actionSlice(square(pawn, 11), square(other2, 57)), nil), // split (5, 2)
		move(card, actionSlice(square(pawn, 12), square(other2, 56)), nil), // split (6, 1)
	)
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard8(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card8)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card8)
	expected = moveSlice(move(card, actionSlice(square(pawn, 14)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(14)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card8)
	expected = emptyMoves // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(14)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card8)
	expected = moveSlice(move(card, actionSlice(square(pawn, 14)), actionSlice(bump(view, model.Green, 1))))
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard10(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card10)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card10)
	expected = moveSlice(
		move(card, actionSlice(square(pawn, 15)), nil),
		move(card, actionSlice(square(pawn, 4)), nil),
	)
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(15)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card10)
	expected = moveSlice(move(card, actionSlice(square(pawn, 4)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(4)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card10)
	expected = moveSlice(move(card, actionSlice(square(pawn, 15)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(15)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card10)
	expected = moveSlice(
		move(card, actionSlice(square(pawn, 15)), actionSlice(bump(view, model.Green, 1))),
		move(card, actionSlice(square(pawn, 4)), nil),
	)
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(4)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card10)
	expected = moveSlice(
		move(card, actionSlice(square(pawn, 15)), nil),
		move(card, actionSlice(square(pawn, 4)), actionSlice(bump(view, model.Green, 1))),
	)
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard11(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card11)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card11)
	expected = moveSlice(move(card, actionSlice(square(pawn, 26)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(26)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card11)
	expected = emptyMoves // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color), which also gets us a swap opportunity
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(26)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card11)
	expected = moveSlice(
		move(card, swap(view, pawn, model.Green, 1), nil),
		move(card, actionSlice(square(pawn, 26)), actionSlice(bump(view, model.Green, 1))),
	)
	assert.Equal(t, expected, moves)

	// Swap pawns elsewhere on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(32)    // can't be swapped, same color
	_ = game.Players()[model.Green].Pawns()[0].Position().MoveToStart()     // can't be swapped, in start area
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSafe(0)    // can't be swapped, in safe area
	_ = game.Players()[model.Yellow].Pawns()[3].Position().MoveToSquare(52) // can be swapped, on board
	_ = game.Players()[model.Blue].Pawns()[1].Position().MoveToSquare(19)   // can be swapped, on board
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card11)
	expected = moveSlice(
		move(card, swap(view, pawn, model.Yellow, 3), nil),
		move(card, swap(view, pawn, model.Blue, 1), nil),
		move(card, actionSlice(square(pawn, 26)), nil),
	)
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCard12(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card12)
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card12)
	expected = moveSlice(move(card, actionSlice(square(pawn, 18)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(18)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card12)
	expected = emptyMoves // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(18)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card12)
	expected = moveSlice(move(card, actionSlice(square(pawn, 18)), actionSlice(bump(view, model.Green, 1))))
	assert.Equal(t, expected, moves)
}

func TestLegalMovesCardApologies(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn in start
	game = setupGame()
	_, _, _, moves = buildMoves(model.Red, game, 0, model.CardApologies)
	_ = game.Players()[model.Yellow].Pawns()[3].Position().MoveToSquare(52) // can be swapped, on board
	_ = game.Players()[model.Blue].Pawns()[1].Position().MoveToSquare(19)   // can be swapped, on board
	expected = emptyMoves
	assert.Equal(t, expected, moves)

	// Swap pawns elsewhere on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Green].Pawns()[0].Position().MoveToStart()     // can't be swapped, in start area
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSafe(0)    // can't be swapped, in safe area
	_ = game.Players()[model.Yellow].Pawns()[3].Position().MoveToSquare(52) // can be swapped, on board
	_ = game.Players()[model.Blue].Pawns()[1].Position().MoveToSquare(19)   // can be swapped, on board
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.CardApologies)
	expected = moveSlice(
		move(card, actionSlice(square(pawn, 52), bump(view, model.Yellow, 3)), nil),
		move(card, actionSlice(square(pawn, 19), bump(view, model.Blue, 1)), nil),
	)
	assert.Equal(t, expected, moves)
}

func TestLegalMovesSpecial(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// Move pawn into safe zone
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(2)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(safe(pawn, 0)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn to home
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSafe(4)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(home(pawn)), nil))
	assert.Equal(t, expected, moves)

	// Move pawn past home
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSafe(4)
	_, _, _, moves = buildMoves(model.Red, game, 0, model.Card2)
	expected = emptyMoves
	assert.Equal(t, expected, moves) // no moves, because it isn't legal

	// Slide of the same color
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(8)
	card, pawn, _, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(square(pawn, 9)), nil))
	assert.Equal(t, expected, moves)

	// Slide of a different color
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(17)
	_ = game.Players()[model.Yellow].Pawns()[2].Position().MoveToSquare(18)
	card, pawn, view, moves = buildMoves(model.Red, game, 0, model.Card1)
	expected = moveSlice(move(card, actionSlice(square(pawn, 19)), actionSlice(bump(view, model.Red, 1), bump(view, model.Yellow, 2))))
	assert.Equal(t, expected, moves)
}

func setupGame() model.Game {
	game, _ := model.NewGame(4, nil)

	for _, color := range model.PlayerColors.Members() {
		for pawn := 0; pawn < model.Pawns; pawn++ {
			_ = game.Players()[color].Pawns()[pawn].Position().MoveToHome()
		}
	}

	return game
}

func buildMoves(color model.PlayerColor, game model.Game, index int, cardType model.CardType) (model.Card, model.Pawn, model.PlayerView, []model.Move) {
	card := model.NewCard("test", cardType)
	view, _ := game.CreatePlayerView(color)
	pawn := view.Player().Pawns()[index]
	moves := NewGenerator().LegalMoves(view.Player().Color(), card, pawn, view.AllPawns())
	return card, pawn, view, moves
}

func calculatePositionSuccess(t *testing.T, color model.PlayerColor, start model.Position, squares int, expected model.Position) {
	result, err := NewGenerator().CalculatePosition(color, start, squares)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func calculatePositionFailure(t *testing.T, color model.PlayerColor, start model.Position, squares int, expected string) {
	_, err := NewGenerator().CalculatePosition(color, start, squares)
	assert.EqualError(t, err, expected)
}

func home(pawn model.Pawn) model.Action {
	return model.NewAction(model.MoveToPosition, pawn, model.NewPosition(false, true, nil, nil))
}

func square(pawn model.Pawn, square int) model.Action {
	return model.NewAction(model.MoveToPosition, pawn, model.NewPosition(false, false, nil, &square))
}

func safe(pawn model.Pawn, square int) model.Action {
	return model.NewAction(model.MoveToPosition, pawn, model.NewPosition(false, false, &square, nil))
}

func start(pawn model.Pawn) model.Action {
	return model.NewAction(model.MoveToStart, pawn, nil)
}

func bump(view model.PlayerView, color model.PlayerColor, index int) model.Action {
	if view.Player().Color() == color {
		return start(view.Player().Pawns()[index])
	} else {
		return start(view.Opponents()[color].Pawns()[index])
	}
}

func swap(view model.PlayerView, pawn model.Pawn, color model.PlayerColor, index int) []model.Action {
	other := view.Opponents()[color].Pawns()[index]
	return actionSlice(
		square(pawn, *other.Position().Square()),
		square(other, *pawn.Position().Square()),
	)
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

func move(card model.Card, actions []model.Action, sideEffects []model.Action) model.Move {
	return model.NewMove(card, actions, sideEffects)
}

func actionSlice(actions ...model.Action) []model.Action {
	result := make([]model.Action, 0, len(actions))

	for i := 0; i < len(actions); i++ {
		result = append(result, actions[i])
	}

	return result
}

func moveSlice(moves ...model.Move) []model.Move {
	result := make([]model.Move, 0, len(moves))

	for i := 0; i < len(moves); i++ {
		result = append(result, moves[i])
	}

	return result
}
