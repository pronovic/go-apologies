package model

import (
	"bytes"
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayer(t *testing.T) {
	obj := NewPlayer(Red)
	assert.Equal(t, Red, obj.Color())
	assert.Equal(t, 0, obj.Turns())
	assert.Equal(t, Pawns, len(obj.Pawns()))
}

func TestNewPlayerFromJSON(t *testing.T) {
	var obj Player
	var err error
	var marshalled []byte
	var unmarshalled Player

	card1 := NewCard("0", CardApologies)
	obj = NewPlayer(Red)
	obj.AppendToHand(card1)
	_ = obj.Pawns()[0].Position().MoveToHome()
	_ = obj.Pawns()[1].Position().MoveToSafe(2)
	_ = obj.Pawns()[2].Position().MoveToSquare(32)

	marshalled, err = json.Marshal(obj)
	assert.Nil(t, err)
	unmarshalled, err = NewPlayerFromJSON(bytes.NewReader(marshalled))
	assert.Nil(t, err)
	assert.Equal(t, obj, unmarshalled)
}

func TestPlayerCopy(t *testing.T) {
	var err error

	card1 := NewCard("0", CardApologies)

	obj := NewPlayer(Red)
	obj.AppendToHand(card1)
	err = obj.Pawns()[0].Position().MoveToHome()
	assert.Nil(t, err)
	err = obj.Pawns()[1].Position().MoveToSafe(2)
	assert.Nil(t, err)
	err = obj.Pawns()[2].Position().MoveToSquare(32)
	assert.Nil(t, err)

	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPlayerPublicData(t *testing.T) {
	var err error

	card1 := NewCard("0", CardApologies)

	obj := NewPlayer(Red)
	obj.AppendToHand(card1)
	err = obj.Pawns()[0].Position().MoveToHome()
	assert.Nil(t, err)
	err = obj.Pawns()[1].Position().MoveToSafe(2)
	assert.Nil(t, err)
	err = obj.Pawns()[2].Position().MoveToSquare(32)
	assert.Nil(t, err)

	expected := obj.Copy()
	expected.RemoveFromHand(card1) // the hand is cleared

	public := obj.PublicData()
	assert.Equal(t, expected, public)
	assert.NotSame(t, obj, public)
}

func TestPlayerAppendAndRemoveHand(t *testing.T) {
	obj := NewPlayer(Red)

	card1 := NewCard("1", Card1)
	card2 := NewCard("2", Card2)
	card3 := NewCard("3", Card3)

	// remove is idempotent; if it's not there, then that's ok
	obj.RemoveFromHand(card1)
	obj.RemoveFromHand(card2)
	obj.RemoveFromHand(card3)

	obj.AppendToHand(card1)
	assert.Equal(t, []Card{card1}, obj.Hand())

	obj.AppendToHand(card2)
	assert.Equal(t, []Card{card1, card2}, obj.Hand())

	obj.AppendToHand(card3)
	assert.Equal(t, []Card{card1, card2, card3}, obj.Hand())

	obj.RemoveFromHand(card2)
	assert.Equal(t, []Card{card1, card3}, obj.Hand())

	obj.RemoveFromHand(card3)
	assert.Equal(t, []Card{card1}, obj.Hand())

	obj.RemoveFromHand(card1)
	assert.Equal(t, []Card{}, obj.Hand())
}

func TestPlayerFindFirstPawnInStart(t *testing.T) {
	obj := NewPlayer(Red)

	for i := 0; i < Pawns; i++ {
		assert.Same(t, obj.Pawns()[i], *obj.FindFirstPawnInStart())
		err := obj.Pawns()[i].Position().MoveToHome()
		assert.Nil(t, err)
	}

	assert.Nil(t, obj.FindFirstPawnInStart())
}

func TestPlayerAllPawnsInHome(t *testing.T) {
	obj := NewPlayer(Red)

	for i := 0; i < Pawns; i++ {
		assert.False(t, obj.AllPawnsInHome())
		err := obj.Pawns()[i].Position().MoveToHome()
		assert.Nil(t, err)
	}

	assert.True(t, obj.AllPawnsInHome())
}

func TestPlayerIncrementTurns(t *testing.T) {
	obj := NewPlayer(Red)
	assert.Equal(t, 0, obj.Turns())
	obj.IncrementTurns()
	assert.Equal(t, 1, obj.Turns())
	obj.IncrementTurns()
	obj.IncrementTurns()
	assert.Equal(t, 3, obj.Turns())
}

func TestNewPlayerView(t *testing.T) {
	player1 := NewPlayer(Blue)
	player2 := NewPlayer(Red)

	opponents := make(map[PlayerColor]Player, 1)
	opponents[Red] = player2

	obj := NewPlayerView(player1, opponents)
	assert.Equal(t, player1, obj.Player())
	assert.Equal(t, opponents, obj.Opponents())
	assert.Equal(t, player2, obj.Opponents()[Red])
}

func TestNewPlayerViewFromJSON(t *testing.T) {
	var obj PlayerView
	var err error
	var marshalled []byte
	var unmarshalled PlayerView

	player1 := NewPlayer(Blue)
	player2 := NewPlayer(Red)
	opponents := make(map[PlayerColor]Player, 1)
	opponents[Red] = player2
	obj = NewPlayerView(player1, opponents)

	marshalled, err = json.Marshal(obj)
	assert.Nil(t, err)
	unmarshalled, err = NewPlayerViewFromJSON(bytes.NewReader(marshalled))
	assert.Nil(t, err)
	assert.Equal(t, obj, unmarshalled)
}

func TestPlayerViewCopy(t *testing.T) {
	player1 := NewPlayer(Blue)
	player2 := NewPlayer(Red)

	opponents := make(map[PlayerColor]Player, 1)
	opponents[Red] = player2

	obj := NewPlayerView(player1, opponents)
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPlayerViewGetPawn(t *testing.T) {
	player1 := NewPlayer(Red)
	opponents := map[PlayerColor]Player{Green: NewPlayer(Green)}
	view := NewPlayerView(player1, opponents)
	assert.Equal(t, view.Player().Pawns()[3], view.GetPawn(NewPawn(Red, 3)))
	assert.Equal(t, view.Opponents()[Green].Pawns()[1], view.GetPawn(NewPawn(Green, 1)))
	assert.Nil(t, view.GetPawn(NewPawn(Yellow, 0)))
}

func TestPlayerViewAllPawns(t *testing.T) {
	player1 := NewPlayer(Red)
	opponents := map[PlayerColor]Player{Green: NewPlayer(Green)}
	view := NewPlayerView(player1, opponents)
	pawns := view.AllPawns()
	assert.Equal(t, 2*Pawns, len(pawns))
	for i := 0; i < Pawns; i++ {
		assert.True(t, slices.Contains(pawns, player1.Pawns()[i]))
		assert.True(t, slices.Contains(pawns, opponents[Green].Pawns()[i]))
	}
}
