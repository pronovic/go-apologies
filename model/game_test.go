package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pronovic/go-apologies/internal/timestamp"
	"github.com/stretchr/testify/assert"
)

var (
	stubbedString       = "2024-01-31T08:15:03.221Z"
	stubbedTimestamp, _ = timestamp.Parse(stubbedString)
	factory             timestamp.MockFactory
)

func init() {
	factory.On("CurrentTime").Return(stubbedTimestamp)
}

func TestNewHistory(t *testing.T) {
	var obj History

	obj = NewHistory("action", nil, nil, &factory)
	assert.Equal(t, "action", obj.Action())
	assert.Nil(t, obj.Color())
	assert.Nil(t, obj.Card())
	assert.Equal(t, stubbedTimestamp, obj.Timestamp())
	assert.Equal(t, fmt.Sprintf("[%s] General - action", stubbedString), fmt.Sprintf("%s", obj))

	color := Blue
	obj = NewHistory("action", &color, nil, &factory)
	assert.Equal(t, &color, obj.Color())
	assert.Nil(t, obj.Card())
	assert.Equal(t, stubbedTimestamp, obj.Timestamp())
	assert.Equal(t, fmt.Sprintf("[%s] Blue - action", stubbedString), fmt.Sprintf("%s", obj))

	card1 := Card12
	obj = NewHistory("action", nil, &card1, &factory)
	assert.Nil(t, obj.Color())
	assert.Equal(t, &card1, obj.Card())
	assert.Equal(t, stubbedTimestamp, obj.Timestamp())
	assert.Equal(t, fmt.Sprintf("[%s] General - action", stubbedString), fmt.Sprintf("%s", obj))
}

func TestNewHistoryFromJSON(t *testing.T) {
	var obj History
	var err error
	var marshalled []byte
	var unmarshalled History

	obj = NewHistory("something", nil, nil, nil)
	marshalled, err = json.Marshal(obj)
	assert.NoError(t, err)
	unmarshalled, err = NewHistoryFromJSON(bytes.NewReader(marshalled))
	assert.NoError(t, err)
	assert.Equal(t, obj, unmarshalled)

	color := Blue
	obj = NewHistory("something", &color, nil, nil)
	marshalled, err = json.Marshal(obj)
	assert.NoError(t, err)
	unmarshalled, err = NewHistoryFromJSON(bytes.NewReader(marshalled))
	assert.NoError(t, err)
	assert.Equal(t, obj, unmarshalled)

	card1 := Card12
	obj = NewHistory("something", nil, &card1, nil)
	marshalled, err = json.Marshal(obj)
	assert.NoError(t, err)
	unmarshalled, err = NewHistoryFromJSON(bytes.NewReader(marshalled))
	assert.NoError(t, err)
	assert.Equal(t, obj, unmarshalled)
}

func TestHistoryCopy(t *testing.T) {
	color := Blue
	card1 := Card12
	obj := NewHistory("action", &color, &card1, nil)
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestNewGameFromJSON(t *testing.T) {
	var obj Game
	var err error
	var marshalled []byte
	var unmarshalled Game

	obj = createRealisticGame()
	marshalled, err = json.Marshal(obj)
	assert.NoError(t, err)
	unmarshalled, err = NewGameFromJSON(bytes.NewReader(marshalled))
	assert.NoError(t, err)
	assert.Equal(t, obj, unmarshalled)
}

func TestNewGame2Players(t *testing.T) {
	game, err := NewGame(2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(game.Players()))
	assert.Equal(t, 0, len(game.History()))
	for _, color := range []PlayerColor{Red, Yellow} {
		assert.Equal(t, color, game.Players()[color].Color())
		assert.Equal(t, 0, len(game.Players()[color].Hand()))
	}
}

func TestNewGame3Players(t *testing.T) {
	game, err := NewGame(3, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(game.Players()))
	assert.Equal(t, 0, len(game.History()))
	for _, color := range []PlayerColor{Red, Yellow, Green} {
		assert.Equal(t, color, game.Players()[color].Color())
		assert.Equal(t, 0, len(game.Players()[color].Hand()))
	}
}

func TestNewGame4Players(t *testing.T) {
	game, err := NewGame(4, nil)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(game.Players()))
	assert.Equal(t, 0, len(game.History()))
	for _, color := range []PlayerColor{Red, Yellow, Green, Blue} {
		assert.Equal(t, color, game.Players()[color].Color())
		assert.Equal(t, 0, len(game.Players()[color].Hand()))
	}
}

func TestNewGameInvalidPlayers(t *testing.T) {
	for _, playerCount := range []int{-2, -1, 0, 1, 5, 6} {
		_, err := NewGame(playerCount, nil)
		assert.EqualError(t, err, "invalid number of players")
	}
}

func TestGameCopy(t *testing.T) {
	game := createRealisticGame()
	copied := game.Copy()
	assert.Equal(t, game, copied)
	assert.NotSame(t, game, copied)
}

func TestGameStarted(t *testing.T) {
	game, _ := NewGame(4, nil)
	assert.False(t, game.Started())
	game.Track("whatever", nil, nil)
	assert.True(t, game.Started())
}

func TestGameCompletedAndWinner(t *testing.T) {
	game, _ := NewGame(4, nil)

	// move all but last pawn into home for all of the players; the game is not complete
	for _, value := range game.Players() {
		for i := 0; i < Pawns-1; i++ {
			assert.False(t, game.Completed())
			_ = value.Pawns()[i].Position().MoveToHome()
		}
	}

	// move the final pawn to home for one player; now the game is complete
	_ = game.Players()[Red].Pawns()[Pawns-1].Position().MoveToHome()
	assert.True(t, game.Completed())
	expected := game.Players()[Red]
	assert.Equal(t, &expected, game.Winner())
}

func TestGameTrackNoPlayer(t *testing.T) {
	game, _ := NewGame(4, &factory)
	game.Track("action", nil, nil)
	assert.Equal(t, NewHistory("action", nil, nil, &factory), game.History()[0])
	assert.Equal(t, 0, game.Players()[Red].Turns())
	assert.Equal(t, 0, game.Players()[Yellow].Turns())
	assert.Equal(t, 0, game.Players()[Blue].Turns())
	assert.Equal(t, 0, game.Players()[Green].Turns())
}

func TestGameTrackWithColor(t *testing.T) {
	game, _ := NewGame(4, &factory)
	player := NewPlayer(Red)
	card := NewCard("x", Card12)
	game.Track("action", player, card)
	assert.Equal(t, NewHistory("action", &Red, &Card12, &factory), game.History()[0])
	assert.Equal(t, 1, game.Players()[Red].Turns())
	assert.Equal(t, 0, game.Players()[Yellow].Turns())
	assert.Equal(t, 0, game.Players()[Blue].Turns())
	assert.Equal(t, 0, game.Players()[Green].Turns())
}

func TestGameCreatePlayerViewInvalid(t *testing.T) {
	game, _ := NewGame(2, nil)
	_, err := game.CreatePlayerView(Blue) // no blue player in 2-player game
	assert.EqualError(t, err, "invalid color")
}

func TestGameCreatePlayerView(t *testing.T) {
	var card Card
	var err error

	game, _ := NewGame(4, nil)

	card, err = game.Deck().Draw()
	assert.NoError(t, err)
	game.Players()[Red].AppendToHand(card)

	card, err = game.Deck().Draw()
	assert.NoError(t, err)
	game.Players()[Yellow].AppendToHand(card)

	card, err = game.Deck().Draw()
	assert.NoError(t, err)
	game.Players()[Green].AppendToHand(card)

	card, err = game.Deck().Draw()
	assert.NoError(t, err)
	game.Players()[Blue].AppendToHand(card)

	view, err := game.CreatePlayerView(Red)
	assert.NoError(t, err)

	assert.NotSame(t, game.Players()[Red], view.Player())
	assert.NotSame(t, game.Players()[Yellow], view.Opponents()[Yellow])

	assert.Equal(t, game.Players()[Red], view.Player())

	for _, color := range []PlayerColor{Yellow, Green, Blue} {
		assert.Equal(t, color, view.Opponents()[color].Color())
		assert.Equal(t, 0, len(view.Opponents()[color].Hand()))
		assert.Equal(t, game.Players()[color].Pawns(), view.Opponents()[color].Pawns())
	}
}

func createRealisticGame() Game {
	// creates a realistic game with changes to the defaults for all types of values
	game, _ := NewGame(4, nil)
	game.Track("this happened", nil, nil)
	game.Track("another thing", game.Players()[Red], nil)
	card1, _ := game.Deck().Draw()
	card2, _ := game.Deck().Draw()
	_, _ = game.Deck().Draw() // just throw it away
	_ = game.Deck().Discard(card1)
	_ = game.Deck().Discard(card2)
	_ = game.Players()[Red].Pawns()[0].Position().MoveToSquare(32)
	_ = game.Players()[Blue].Pawns()[2].Position().MoveToHome()
	game.Players()[Blue].AppendToHand(card1)
	_ = game.Players()[Yellow].Pawns()[3].Position().MoveToSafe(1)
	_ = game.Players()[Green].Pawns()[1].Position().MoveToSquare(19)
	game.Players()[Green].AppendToHand(card2)
	return game
}
