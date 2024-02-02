package rules

import (
	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

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