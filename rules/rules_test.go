package rules

import (
	"testing"

	"github.com/pronovic/go-apologies/generator"
	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
)

func TestStartGameStandardMode(t *testing.T) {
	game, _ := model.NewGame(2, nil)
	err := NewRules(nil).StartGame(game, model.StandardMode)

	assert.Nil(t, err)
	assert.True(t, game.Started())

	assert.Equal(t, model.Red, game.Players()[model.Red].Color())
	assert.Equal(t, 0, len(game.Players()[model.Red].Hand()))

	assert.Equal(t, model.Yellow, game.Players()[model.Yellow].Color())
	assert.Equal(t, 0, len(game.Players()[model.Yellow].Hand()))

	err = NewRules(nil).StartGame(game, model.StandardMode)
	assert.EqualError(t, err, "game is already started")
}

func TestStartGameAdultMode(t *testing.T) {
	game, _ := model.NewGame(4, nil)
	err := NewRules(nil).StartGame(game, model.AdultMode)

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

	err = NewRules(nil).StartGame(game, model.AdultMode)
	assert.EqualError(t, err, "game is already started")
}

func TestExecuteMove(t *testing.T) {
	actions := []model.Action{
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Red, 1), positionSquare(10)),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Yellow, 3), positionSquare(11)),
	}

	sideEffects := []model.Action{
		model.NewAction(model.MoveToStart, model.NewPawn(model.Blue, 2), nil),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Green, 0), positionSquare(12)),
	}

	move := model.NewMove(model.NewCard("1", model.Card1), actions, sideEffects)

	game, _ := model.NewGame(4, nil)
	player := game.Players()[model.Red]

	err := NewRules(nil).ExecuteMove(game, player, move)
	assert.Nil(t, err)

	assert.Equal(t, 10, *game.Players()[model.Red].Pawns()[1].Position().Square())
	assert.Equal(t, 11, *game.Players()[model.Yellow].Pawns()[3].Position().Square())
	assert.True(t, game.Players()[model.Blue].Pawns()[2].Position().Start())
	assert.Equal(t, 12, *game.Players()[model.Green].Pawns()[0].Position().Square())
}

func TestEvaluateMove(t *testing.T) {
	var err error
	var result model.PlayerView

	actions := []model.Action{
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Red, 1), positionSquare(10)),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Yellow, 3), positionSquare(11)),
	}

	sideEffects := []model.Action{
		model.NewAction(model.MoveToStart, model.NewPawn(model.Blue, 2), nil),
		model.NewAction(model.MoveToPosition, model.NewPawn(model.Green, 0), positionSquare(12)),
	}

	move := model.NewMove(model.NewCard("1", model.Card1), actions, sideEffects)

	game, _ := model.NewGame(4, nil)
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

	result, err = NewRules(nil).EvaluateMove(view, move)
	assert.Equal(t, expected, result)
}

// For these ConstructLegalMoves() tests, I have replicated the test design from
// the Python implementation, although it's a bit more awkward here in Go than
// it was in Python.

func TestConstructLegalMovesNoMovesWithCard(t *testing.T) {
	card := model.NewCard("card", model.Card10)

	hand1 := model.NewCard("hand1", model.Card1)
	hand2 := model.NewCard("hand2", model.Card2)
	hand := []model.Card{hand1, hand2}

	pawn1 := model.NewPawn(model.Red, 0)
	pawn2 := model.NewPawn(model.Red, 1)
	playerPawns := []model.Pawn{pawn1, pawn2}

	dummy1 := model.MockPawn{}
	dummy2 := model.MockPawn{}
	allPawns := []model.Pawn{&dummy1, &dummy2}

	cardPawn1Moves := make([]model.Move, 0)
	cardPawn2Moves := make([]model.Move, 0)

	// result is a forfeit for the only card
	expectedMoves := []model.Move{move(card, nil, nil)}

	player := model.MockPlayer{}
	player.On("Color").Return(model.Red)
	player.On("Hand").Return(hand)
	player.On("Pawns").Return(playerPawns)

	view := model.MockPlayerView{}
	view.On("Player").Return(&player)
	view.On("AllPawns").Return(allPawns)

	var moveGenerator generator.MockMoveGenerator
	moveGenerator.On("LegalMoves", model.Red, card, pawn1, allPawns).Return(cardPawn1Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, card, pawn2, allPawns).Return(cardPawn2Moves).Once()

	rules := NewRules(&moveGenerator)
	result, err := rules.ConstructLegalMoves(&view, card)
	assert.Nil(t, err)
	assert.Equal(t, expectedMoves, result)
}

func TestConstructLegalMovesNoMovesNoCard(t *testing.T) {
	var card model.Card = nil

	hand1 := model.NewCard("hand1", model.Card1)
	hand2 := model.NewCard("hand2", model.Card2)
	hand := []model.Card{hand1, hand2}

	pawn1 := model.NewPawn(model.Red, 0)
	pawn2 := model.NewPawn(model.Red, 1)
	playerPawns := []model.Pawn{pawn1, pawn2}

	dummy1 := model.MockPawn{}
	dummy2 := model.MockPawn{}
	allPawns := []model.Pawn{&dummy1, &dummy2}

	hand1Pawn1Moves := make([]model.Move, 0)
	hand1Pawn2Moves := make([]model.Move, 0)
	hand2Pawn1Moves := make([]model.Move, 0)
	hand2Pawn2Moves := make([]model.Move, 0)

	// result is a forfeit for all cards in the hand
	expectedMoves := []model.Move{move(hand1, nil, nil), move(hand2, nil, nil)}

	player := model.MockPlayer{}
	player.On("Color").Return(model.Red)
	player.On("Hand").Return(hand)
	player.On("Pawns").Return(playerPawns)

	view := model.MockPlayerView{}
	view.On("Player").Return(&player)
	view.On("AllPawns").Return(allPawns)

	var moveGenerator generator.MockMoveGenerator
	moveGenerator.On("LegalMoves", model.Red, hand1, pawn1, allPawns).Return(hand1Pawn1Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, hand1, pawn2, allPawns).Return(hand1Pawn2Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, hand2, pawn1, allPawns).Return(hand2Pawn1Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, hand2, pawn2, allPawns).Return(hand2Pawn2Moves).Once()

	rules := NewRules(&moveGenerator)
	result, err := rules.ConstructLegalMoves(&view, card)
	assert.Nil(t, err)
	assert.Equal(t, expectedMoves, result)
}

func TestConstructLegalMovesWithMovesWithCard(t *testing.T) {
	card := model.NewCard("card", model.Card10)

	hand1 := model.NewCard("hand1", model.Card1)
	hand2 := model.NewCard("hand2", model.Card2)
	hand := []model.Card{hand1, hand2}

	pawn1 := model.NewPawn(model.Red, 0)
	pawn2 := model.NewPawn(model.Red, 1)
	playerPawns := []model.Pawn{pawn1, pawn2}

	dummy1 := model.MockPawn{}
	dummy2 := model.MockPawn{}
	allPawns := []model.Pawn{&dummy1, &dummy2}

	cardPawn1Moves := []model.Move{
		move(card, []model.Action{actionStart(pawn1)}, nil),
		move(card, []model.Action{actionStart(pawn1)}, nil),
	}

	cardPawn2Moves := []model.Move{
		move(card, []model.Action{actionPosition(pawn2), actionStart(pawn2)}, nil),
	}

	// result is a list of all returned moves, with duplicates removed
	expectedMoves := []model.Move{
		move(card, []model.Action{actionStart(pawn1)}, nil),
		move(card, []model.Action{actionPosition(pawn2), actionStart(pawn2)}, nil),
	}

	player := model.MockPlayer{}
	player.On("Color").Return(model.Red)
	player.On("Hand").Return(hand)
	player.On("Pawns").Return(playerPawns)

	view := model.MockPlayerView{}
	view.On("Player").Return(&player)
	view.On("AllPawns").Return(allPawns)

	var moveGenerator generator.MockMoveGenerator
	moveGenerator.On("LegalMoves", model.Red, card, pawn1, allPawns).Return(cardPawn1Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, card, pawn2, allPawns).Return(cardPawn2Moves).Once()

	rules := NewRules(&moveGenerator)
	result, err := rules.ConstructLegalMoves(&view, card)
	assert.Nil(t, err)
	assert.Equal(t, expectedMoves, result)
}

func TestConstructLegalMovesWithMovesNoCard(t *testing.T) {
	var card model.Card = nil

	hand1 := model.NewCard("hand1", model.Card1)
	hand2 := model.NewCard("hand2", model.Card2)
	hand := []model.Card{hand1, hand2}

	pawn1 := model.NewPawn(model.Red, 0)
	pawn2 := model.NewPawn(model.Red, 1)
	playerPawns := []model.Pawn{pawn1, pawn2}

	dummy1 := model.MockPawn{}
	dummy2 := model.MockPawn{}
	allPawns := []model.Pawn{&dummy1, &dummy2}

	hand1Pawn1Moves := []model.Move{
		move(card, []model.Action{actionStart(pawn1)}, nil),
		move(card, []model.Action{actionStart(pawn1)}, nil),
	}

	hand1Pawn2Moves := []model.Move{
		move(card, []model.Action{actionStart(pawn2), actionPosition(pawn2)}, nil),
	}

	hand2Pawn1Moves := []model.Move{
		move(card, []model.Action{actionPosition(pawn1)}, nil),
	}

	hand2Pawn2Moves := []model.Move{
		move(card, []model.Action{actionPosition(pawn2)}, nil),
	}

	// result is a list of all returned moves, with duplicates removed
	expectedMoves := []model.Move{
		move(card, []model.Action{actionStart(pawn1)}, nil),
		move(card, []model.Action{actionStart(pawn2), actionPosition(pawn2)}, nil),
		move(card, []model.Action{actionPosition(pawn1)}, nil),
		move(card, []model.Action{actionPosition(pawn2)}, nil),
	}

	player := model.MockPlayer{}
	player.On("Color").Return(model.Red)
	player.On("Hand").Return(hand)
	player.On("Pawns").Return(playerPawns)

	view := model.MockPlayerView{}
	view.On("Player").Return(&player)
	view.On("AllPawns").Return(allPawns)

	var moveGenerator generator.MockMoveGenerator
	moveGenerator.On("LegalMoves", model.Red, hand1, pawn1, allPawns).Return(hand1Pawn1Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, hand1, pawn2, allPawns).Return(hand1Pawn2Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, hand2, pawn1, allPawns).Return(hand2Pawn1Moves).Once()
	moveGenerator.On("LegalMoves", model.Red, hand2, pawn2, allPawns).Return(hand2Pawn2Moves).Once()

	rules := NewRules(&moveGenerator)
	result, err := rules.ConstructLegalMoves(&view, card)
	assert.Nil(t, err)
	assert.Equal(t, expectedMoves, result)
}

func TestDrawAgain(t *testing.T) {
	rules := NewRules(nil)
	for _, cardType := range model.CardTypes.Members() {
		assert.Equal(t, model.DrawAgain[cardType], rules.DrawAgain(model.NewCard("id", cardType)))
	}
}

func actionPosition(pawn model.Pawn) model.Action {
	return model.NewAction(model.MoveToPosition, pawn, nil)
}

func actionStart(pawn model.Pawn) model.Action {
	return model.NewAction(model.MoveToStart, pawn, nil)
}

func positionSquare(square int) model.Position {
	return model.NewPosition(false, false, nil, &square)
}

func move(card model.Card, actions []model.Action, sideEffects []model.Action) model.Move {
	return model.NewMove(card, actions, sideEffects)
}
