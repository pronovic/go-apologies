package pkg

import (
	"github.com/pronovic/go-apologies/internal/identifier"
	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	identifier.UseStubbedId()  // once this has been called, it takes effect permanently for all unit tests
}

func TestStartGameStandardMode(t *testing.T) {
	game, _ := model.NewGame(2)
	err := StartGame(game, model.StandardMode)

	assert.Nil(t, err)
	assert.True(t, game.Started())

	assert.Equal(t, model.Red, game.Players()[model.Red].Color())
	assert.Equal(t, 0, len(game.Players()[model.Red].Hand()))

	assert.Equal(t, model.Yellow, game.Players()[model.Yellow].Color())
	assert.Equal(t, 0, len(game.Players()[model.Yellow].Hand()))

	err = StartGame(game, model.StandardMode)
	assert.EqualError(t, err, "game is already started")
}

func TestStartGameAdultMode(t *testing.T) {
	game, _ := model.NewGame(4)
	err := StartGame(game, model.AdultMode)

	assert.Equal(t, model.Red, game.Players()[model.Red].Color())
	assert.Equal(t, model.AdultHand, len(game.Players()[model.Red].Hand()))
	assert.Equal(t, 4, *game.Players()[model.Red].Pawns()[0].Position().Square())

	assert.Equal(t, model.Yellow, game.Players()[model.Yellow].Color())
	assert.Equal(t, model.AdultHand, len(game.Players()[model.Yellow].Hand()))
	assert.Equal(t, 34, *game.Players()[model.Yellow].Pawns()[0].Position().Square())

	assert.Equal(t, model.Green, game.Players()[model.Green].Color())
	assert.Equal(t, model.AdultHand, len(game.Players()[model.Green].Hand()))
	assert.Equal(t, 49, *game.Players()[model.Green].Pawns()[0].Position().Square())

	assert.Equal(t, model.Blue, game.Players()[model.Blue].Color())
	assert.Equal(t, model.AdultHand, len(game.Players()[model.Blue].Hand()))
	assert.Equal(t, 19, *game.Players()[model.Blue].Pawns()[0].Position().Square())

	err = StartGame(game, model.AdultMode)
	assert.EqualError(t, err, "game is already started")
}

func TestExecuteMove(t *testing.T) {
	actions := []model.Action {
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Red, 1), positionSquare(10)),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Yellow, 3), positionSquare(11)),
	}

	sideEffects := []model.Action {
		model.NewAction(model.MoveToStart, model.NewPawn(model.Blue, 2), nil),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Green, 0), positionSquare(12)),
	}

	move := model.NewMove(model.NewCard("1", model.Card1), actions, sideEffects)

	game, _ := model.NewGame(4)
	player := game.Players()[model.Red]

	err := ExecuteMove(game, player, move)
	assert.Nil(t, err)

	assert.Equal(t, 10, *game.Players()[model.Red].Pawns()[1].Position().Square())
	assert.Equal(t, 11, *game.Players()[model.Yellow].Pawns()[3].Position().Square())
	assert.True(t, game.Players()[model.Blue].Pawns()[2].Position().Start())
	assert.Equal(t, 12, *game.Players()[model.Green].Pawns()[0].Position().Square())
}

func TestEvaluateMove(t *testing.T) {
	var err error
	var result model.PlayerView

	actions := []model.Action {
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Red, 1), positionSquare(10)),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Yellow, 3), positionSquare(11)),
	}

	sideEffects := []model.Action {
		model.NewAction(model.MoveToStart, model.NewPawn(model.Blue, 2), nil),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Green, 0), positionSquare(12)),
	}

	move := model.NewMove(model.NewCard("1", model.Card1), actions, sideEffects)

	game, _ := model.NewGame(4)
	view, err := game.CreatePlayerView(model.Red)
	assert.Nil(t, err)

	expected := view.Copy()

	err = expected.Player().Pawns()[1].Position().MoveToSquare(10)
	assert.Nil(t, err)

	err = expected.Opponents()[model.Yellow].Pawns()[3].Position().MoveToSquare(11)
	assert.Nil(t, err)

	err = expected.Opponents()[model.Yellow].Pawns()[2].Position().MoveToStart()
	assert.Nil(t, err)

	err = expected.Opponents()[model.Green].Pawns()[0].Position().MoveToSquare(12)
	assert.Nil(t, err)

	result, err = EvaluateMove(view, move)
	assert.Equal(t, expected, result)
}

func TestDistanceToHome(t *testing.T) {
	// distance from home is always 0
	for _, color := range []model.PlayerColor{ model.Red, model.Yellow, model.Green } {
		assert.Equal(t, 0, DistanceToHome(pawnHome(color)))
	}

	// distance from start is always 65
	for _, color := range []model.PlayerColor{ model.Red, model.Yellow, model.Green } {
		assert.Equal(t, 65, DistanceToHome(pawnStart(color)))
	}

	// distance from within safe is always <= 5
	assert.Equal(t, 5, DistanceToHome(pawnSafe(model.Red, 0)))
	assert.Equal(t, 4, DistanceToHome(pawnSafe(model.Red, 1)))
	assert.Equal(t, 3, DistanceToHome(pawnSafe(model.Red, 2)))
	assert.Equal(t, 2, DistanceToHome(pawnSafe(model.Red, 3)))
	assert.Equal(t, 1, DistanceToHome(pawnSafe(model.Red, 4)))

	// distance from circle is always 64
	assert.Equal(t, 64, DistanceToHome(pawnSquare(model.Red, 4)))
	assert.Equal(t, 64, DistanceToHome(pawnSquare(model.Blue, 19)))
	assert.Equal(t, 64, DistanceToHome(pawnSquare(model.Yellow, 34)))
	assert.Equal(t, 64, DistanceToHome(pawnSquare(model.Green, 49)))

	// distance from square between turn and circle is always 65
	assert.Equal(t, 65, DistanceToHome(pawnSquare(model.Red, 3)))
	assert.Equal(t, 65, DistanceToHome(pawnSquare(model.Blue, 18)))
	assert.Equal(t, 65, DistanceToHome(pawnSquare(model.Yellow, 33)))
	assert.Equal(t, 65, DistanceToHome(pawnSquare(model.Green, 48)))

	// distance from turn is always 6
	assert.Equal(t, 6, DistanceToHome(pawnSquare(model.Red, 2)))
	assert.Equal(t, 6, DistanceToHome(pawnSquare(model.Blue, 17)))
	assert.Equal(t, 6, DistanceToHome(pawnSquare(model.Yellow, 32)))
	assert.Equal(t, 6, DistanceToHome(pawnSquare(model.Green, 47)))

	// check some arbitrary squares
	assert.Equal(t, 7, DistanceToHome(pawnSquare(model.Red, 1)))
	assert.Equal(t, 8, DistanceToHome(pawnSquare(model.Red, 0)))
	assert.Equal(t, 9, DistanceToHome(pawnSquare(model.Red, 59)))
	assert.Equal(t, 59, DistanceToHome(pawnSquare(model.Red, 9)))
	assert.Equal(t, 23, DistanceToHome(pawnSquare(model.Blue, 0)))
	assert.Equal(t, 13, DistanceToHome(pawnSquare(model.Green, 40)))
}

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

func TestConstructLegalMovesNoMovesWithCard(t *testing.T) {
	// TODO: implement test
}

func TestConstructLegalMovesNoMovesNoCard(t *testing.T) {
	// TODO: implement test
}

func TestConstructLegalMovesWithMovesWithCard(t *testing.T) {
	// TODO: implement test
}

func TestConstructLegalMovesWithMovesNoCard(t *testing.T) {
	// TODO: implement test
}

func TestConstructLegalMovesCard1(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn in start, on the board, or in safe
	game = setupGame()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn from start with no conflicts
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(4)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move{}  // can't start because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSquare(4)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action { actionBump(view, model.Yellow, 0)}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 7)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(7)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move{}  // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(7)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 7)}, []model.Action { actionBump(view, model.Green, 1)}) }
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard2(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn in start, on the board, or in safe
	game = setupGame()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn from start with no conflicts
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(4)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move{}  // can't start because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn from start with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSquare(4)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action { actionBump(view, model.Yellow, 0)}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 8)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(8)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move{}  // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(8)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 8)}, []model.Action { actionBump(view, model.Green, 1)}) }
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard3(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card3)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card3)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 9)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(9)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card3)
	expected = []model.Move{}  // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(9)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card3)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 9)}, []model.Action { actionBump(view, model.Green, 1)}) }
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard4(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card4)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card4)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 2)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(2)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card4)
	expected = []model.Move{ }
	assert.Equal(t, expected, moves) // can't move because we have a pawn there already

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(2)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card4)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 2)}, []model.Action { actionBump(view, model.Green, 1)}) }
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard5(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card5)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card5)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 11)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(11)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card5)
	expected = []model.Move{}  // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(11)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card5)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 11)}, []model.Action { actionBump(view, model.Green, 1)}) }
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard7(t *testing.T) {
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
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card7)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// One move available if there is one pawn on the board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card7)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 13)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Multiple moves available if there is more than one pawn on the board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[2].Position().MoveToSquare(55)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card7)
	other = view.Player().Pawns()[2]
	expected = []model.Move {
		model.NewMove(card, []model.Action { actionSquare(pawn, 13)}, []model.Action {}), // move our pawn 7
		model.NewMove(card, []model.Action { actionSquare(pawn, 7), actionSquare(other, 1) }, []model.Action {}), // split (1, 6)
		model.NewMove(card, []model.Action { actionSquare(pawn, 8), actionSquare(other, 0) }, []model.Action {}), // split (2, 5)
		model.NewMove(card, []model.Action { actionSquare(pawn, 9), actionSquare(other, 59) }, []model.Action {}), // split (3, 4)
		model.NewMove(card, []model.Action { actionSquare(pawn, 10), actionSquare(other, 58) }, []model.Action {}), // split (4, 3)
		model.NewMove(card, []model.Action { actionSquare(pawn, 11), actionSquare(other, 57) }, []model.Action {}), // split (5, 2)
		model.NewMove(card, []model.Action { actionSquare(pawn, 12), actionSquare(other, 56) }, []model.Action {}), // split (6, 1)
	}
	assert.Equal(t, expected, moves)

	// Either half of a move might bump an opponent back to start
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[2].Position().MoveToSquare(55)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(10)
	_ = game.Players()[model.Blue].Pawns()[3].Position().MoveToSquare(56)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card7)
	other = view.Player().Pawns()[2]
	expected = []model.Move {
		model.NewMove(card, []model.Action { actionSquare(pawn, 13)}, []model.Action {}), // move our pawn 7
		model.NewMove(card, []model.Action { actionSquare(pawn, 7), actionSquare(other, 1) }, []model.Action {}), // split (1, 6)
		model.NewMove(card, []model.Action { actionSquare(pawn, 8), actionSquare(other, 0) }, []model.Action {}), // split (2, 5)
		model.NewMove(card, []model.Action { actionSquare(pawn, 9), actionSquare(other, 59) }, []model.Action {}), // split (3, 4)
		model.NewMove(card, []model.Action { actionSquare(pawn, 10), actionSquare(other, 58) }, []model.Action { actionBump(view, model.Green, 1)}), // split (4, 3)
		model.NewMove(card, []model.Action { actionSquare(pawn, 11), actionSquare(other, 57) }, []model.Action {}), // split (5, 2)
		model.NewMove(card, []model.Action { actionSquare(pawn, 12), actionSquare(other, 56) }, []model.Action { actionBump(view, model.Blue, 3)}), // split (6, 1)
	}
	assert.Equal(t, expected, moves)

	// If either half of the move has a conflict with another pawn of the same color, the entire move is invalidated
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(9)
	_ = game.Players()[model.Red].Pawns()[2].Position().MoveToSquare(55)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card7)
	other1 := view.Player().Pawns()[1]
	other2 := view.Player().Pawns()[2]
	expected = []model.Move {
		model.NewMove(card, []model.Action { actionSquare(pawn, 13)}, []model.Action {}), // move our pawn 7
		model.NewMove(card, []model.Action { actionSquare(pawn, 7), actionSquare(other1, 15) }, []model.Action {}), // split (1, 6)
		model.NewMove(card, []model.Action { actionSquare(pawn, 8), actionSquare(other1, 14) }, []model.Action {}), // split (2, 5)
		model.NewMove(card, []model.Action { actionSquare(pawn, 9), actionSquare(other1, 13) }, []model.Action {}), // split (3, 4)
		model.NewMove(card, []model.Action { actionSquare(pawn, 10), actionSquare(other1, 12) }, []model.Action { }), // split (4, 3)
		model.NewMove(card, []model.Action { actionSquare(pawn, 11), actionSquare(other1, 11) }, []model.Action {}), // split (5, 2)
		model.NewMove(card, []model.Action { actionSquare(pawn, 12), actionSquare(other1, 10) }, []model.Action { }), // split (6, 1)
		model.NewMove(card, []model.Action { actionSquare(pawn, 7), actionSquare(other2, 1) }, []model.Action {}), // split (1, 6)
		model.NewMove(card, []model.Action { actionSquare(pawn, 8), actionSquare(other2, 0) }, []model.Action {}), // split (2, 5)
		// the move for square 9 is disallowed because pawn[1] already lives there, and isn't part of this action
		model.NewMove(card, []model.Action { actionSquare(pawn, 10), actionSquare(other2, 58) }, []model.Action { }), // split (4, 3)
		model.NewMove(card, []model.Action { actionSquare(pawn, 11), actionSquare(other2, 57) }, []model.Action {}), // split (5, 2)
		model.NewMove(card, []model.Action { actionSquare(pawn, 12), actionSquare(other2, 56) }, []model.Action { }), // split (6, 1)
	}
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard8(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card8)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card8)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 14)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(14)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card8)
	expected = []model.Move{}  // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(14)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card8)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 14)}, []model.Action { actionBump(view, model.Green, 1)}) }
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard10(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card10)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card10)
	expected = []model.Move {
		model.NewMove(card, []model.Action { actionSquare(pawn, 15)}, []model.Action {}),
		model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action {}),
	}
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(15)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card10)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(4)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card10)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 15)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(15)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card10)
	expected = []model.Move {
		model.NewMove(card, []model.Action { actionSquare(pawn, 15)}, []model.Action { actionBump(view, model.Green, 1)}),
		model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action { }),
	}
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(5)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(4)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card10)
	expected = []model.Move {
		model.NewMove(card, []model.Action { actionSquare(pawn, 15)}, []model.Action { }),
		model.NewMove(card, []model.Action { actionSquare(pawn, 4)}, []model.Action { actionBump(view, model.Green, 1) }),
	}
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard11(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card11)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card11)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 26)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(26)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card11)
	expected = []model.Move{}  // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color), which also gets us a swap opportunity
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(26)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card11)
	expected = []model.Move {
		model.NewMove(card, actionSwap(view, pawn, model.Green, 1), []model.Action {}),
		model.NewMove(card, []model.Action { actionSquare(pawn, 26)}, []model.Action { actionBump(view, model.Green, 1)}),
	}
	assert.Equal(t, expected, moves)

	// Swap pawns elsewhere on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(32)  // can't be swapped, same color
	_ = game.Players()[model.Green].Pawns()[0].Position().MoveToStart()  // can't be swapped, in start area
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSafe(0)  // can't be swapped, in safe area
	_ = game.Players()[model.Yellow].Pawns()[3].Position().MoveToSquare(52)  // can be swapped, on board
	_ = game.Players()[model.Blue].Pawns()[1].Position().MoveToSquare(19)  // can be swapped, on board
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card11)
	expected = []model.Move {
		model.NewMove(card, actionSwap(view, pawn, model.Yellow, 3), []model.Action {}),
		model.NewMove(card, actionSwap(view, pawn, model.Blue, 1), []model.Action {}),
		model.NewMove(card, []model.Action { actionSquare(pawn, 26)}, []model.Action { }),
	}
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCard12(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn on the board, or in safe
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToHome()
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card12)
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Move pawn on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card12)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 18)}, []model.Action {}) }
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (same color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(18)
	_, _, _, moves = legalMoves(model.Red, game, 0, model.Card12)
	expected = []model.Move{}  // can't move because we have a pawn there already
	assert.Equal(t, expected, moves)

	// Move pawn on board with conflict (different color)
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(6)
	_ = game.Players()[model.Green].Pawns()[1].Position().MoveToSquare(18)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card12)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 18)}, []model.Action { actionBump(view, model.Green, 1)}) }
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCardApologies(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// No legal moves if no pawn in start
	game = setupGame()
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.CardApologies)
	_ = game.Players()[model.Yellow].Pawns()[3].Position().MoveToSquare(52) // can be swapped, on board
	_ = game.Players()[model.Blue].Pawns()[1].Position().MoveToSquare(19) // can be swapped, on board
	expected = []model.Move{}
	assert.Equal(t, expected, moves)

	// Swap pawns elsewhere on board
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToStart()
	_ = game.Players()[model.Green].Pawns()[0].Position().MoveToStart() // can't be swapped, in start area
	_ = game.Players()[model.Yellow].Pawns()[0].Position().MoveToSafe(0) // can't be swapped, in safe area
	_ = game.Players()[model.Yellow].Pawns()[3].Position().MoveToSquare(52) // can be swapped, on board
	_ = game.Players()[model.Blue].Pawns()[1].Position().MoveToSquare(19) // can be swapped, on board
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.CardApologies)
	expected = []model.Move {
		model.NewMove(card, []model.Action { actionSquare(pawn, 52), actionBump(view, model.Yellow, 3) }, []model.Action { }),
		model.NewMove(card, []model.Action { actionSquare(pawn, 19), actionBump(view, model.Blue, 1) }, []model.Action { }),
	}
	assert.Equal(t, expected, moves)
}

func TestConstructLegalMovesCardSpecial(t *testing.T) {
	var card model.Card
	var pawn model.Pawn
	var view model.PlayerView
	var moves []model.Move
	var expected []model.Move
	var game model.Game

	// Move pawn into safe zone
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(2)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSafe(pawn, 0) }, []model.Action {})}
	assert.Equal(t, expected, moves)

	// Move pawn to home
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSafe(4)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionHome(pawn) }, []model.Action {})}
	assert.Equal(t, expected, moves)

	// Move pawn past home
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSafe(4)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card2)
	expected = []model.Move{}
	assert.Equal(t, expected, moves) // no moves, because it isn't legal

	// Slide of the same color
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(8)
	card, pawn, _, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 9) }, []model.Action {})}
	assert.Equal(t, expected, moves)

	// Slide of a different color
	game = setupGame()
	_ = game.Players()[model.Red].Pawns()[0].Position().MoveToSquare(15)
	_ = game.Players()[model.Red].Pawns()[1].Position().MoveToSquare(17)
	_ = game.Players()[model.Yellow].Pawns()[2].Position().MoveToSquare(18)
	card, pawn, view, moves = legalMoves(model.Red, game, 0, model.Card1)
	expected = []model.Move { model.NewMove(card, []model.Action { actionSquare(pawn, 19) }, []model.Action { actionBump(view, model.Red, 1), actionBump(view, model.Yellow, 2) })}
	assert.Equal(t, expected, moves)
}

func setupGame() model.Game {
	game, _ := model.NewGame(4)

	for _, color := range model.PlayerColors.Members() {
		for pawn := 0; pawn < model.Pawns; pawn ++ {
			_ = game.Players()[color].Pawns()[pawn].Position().MoveToHome()
		}
	}

	return game
}

func legalMoves(color model.PlayerColor, game model.Game, index int, cardType model.CardType) (model.Card, model.Pawn, model.PlayerView, []model.Move) {
	card := model.NewCard("test", cardType)
	view, _ := game.CreatePlayerView(color)
	pawn := view.Player().Pawns()[index]
	moves := constructLegalMoves(view.Player().Color(), card, pawn, view.AllPawns())
	return card, pawn, view, moves
}

func calculatePositionSuccess(t *testing.T, color model.PlayerColor, start model.Position, squares int, expected model.Position) {
	result, err := calculatePosition(color, start, squares)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func calculatePositionFailure(t *testing.T, color model.PlayerColor, start model.Position, squares int, expected string) {
	_, err := calculatePosition(color, start, squares)
	assert.EqualError(t, err, expected)
}


func actionHome(pawn model.Pawn) model.Action {
	return model.NewAction(model.MoveToPosition, pawn, model.NewPosition(false, true, nil, nil))
}

func actionSquare(pawn model.Pawn, square int) model.Action {
	return model.NewAction(model.MoveToPosition, pawn, model.NewPosition(false, false, nil, &square))
}

func actionSafe(pawn model.Pawn, square int) model.Action {
	return model.NewAction(model.MoveToPosition, pawn, model.NewPosition(false, false, &square, nil))
}

func actionStart(pawn model.Pawn) model.Action {
	return model.NewAction(model.MoveToStart, pawn, nil)
}

func actionBump(view model.PlayerView, color model.PlayerColor, index int) model.Action {
	if view.Player().Color() == color {
		return actionStart(view.Player().Pawns()[index])
	} else {
		return actionStart(view.Opponents()[color].Pawns()[index])
	}
}

func actionSwap(view model.PlayerView, pawn model.Pawn, color model.PlayerColor, index int) []model.Action {
	other := view.Opponents()[color].Pawns()[index]
	return []model.Action {
		actionSquare(pawn, *other.Position().Square()),
		actionSquare(other, *pawn.Position().Square()),
	}
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

