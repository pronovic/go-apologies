package model

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAction(t *testing.T) {
	pawn := NewPawn(Red, 0)
	position := NewPosition(false, false, nil, nil)
	obj := NewAction(MoveToPosition, pawn, position)
	assert.Equal(t, MoveToPosition, obj.Type())
	assert.Same(t, pawn, obj.Pawn())
	assert.Same(t, position, obj.Position())
}

func TestNewActionFromJSON(t *testing.T) {
	var obj Action
	var err error
	var marshalled []byte
	var unmarshalled Action

	pawn1 := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	obj = NewAction(MoveToPosition, pawn1, position1)
	marshalled, err = json.Marshal(obj)
	assert.NoError(t, err)
	unmarshalled, err = NewActionFromJSON(bytes.NewReader(marshalled))
	assert.NoError(t, err)
	assert.Equal(t, obj, unmarshalled)
}

func TestActionSetPosition(t *testing.T) {
	pawn := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	position2 := NewPosition(true, false, nil, nil)
	obj := NewAction(MoveToPosition, pawn, position1)
	obj.SetPosition(position2)
	assert.Same(t, position2, obj.Position())
}

func TestNewMove(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 1, 2)
	sideEffects := make([]Action, 2, 3)
	obj := NewMove(card, actions, sideEffects)
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())
	assert.Equal(t, sideEffects, obj.SideEffects())
}

func TestNewMoveFromJSON(t *testing.T) {
	var obj Move
	var err error
	var marshalled []byte
	var unmarshalled Move

	card1 := NewCard("4", Card4)
	position1 := NewPosition(false, false, nil, nil)
	action1 := NewAction(MoveToStart, nil, position1)
	pawn2 := NewPawn(Red, 0)
	action2 := NewAction(MoveToStart, pawn2, nil)
	obj = NewMove(card1, []Action{action1}, []Action{action2})

	marshalled, err = json.Marshal(obj)
	assert.NoError(t, err)
	unmarshalled, err = NewMoveFromJSON(bytes.NewReader(marshalled))
	assert.NoError(t, err)
	assert.Equal(t, obj, unmarshalled)
}

func TestNewMoveEmptySlice(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 0)
	sideEffects := make([]Action, 0)
	obj := NewMove(card, actions, sideEffects)
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())
	assert.Equal(t, sideEffects, obj.SideEffects())
}

func TestNewMoveNilSlice(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 0)
	sideEffects := make([]Action, 0)
	obj := NewMove(card, nil, nil)
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())         // nil is converted to a newly-allocated empty slice
	assert.Equal(t, sideEffects, obj.SideEffects()) // nil is converted to a newly-allocated empty slice
}

func TestMoveAddSideEffect(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 0)
	sideEffects := make([]Action, 0)
	obj := NewMove(card, actions, sideEffects)

	pawn := NewPawn(Red, 0)
	position := NewPosition(false, false, nil, nil)
	sideEffect := NewAction(MoveToPosition, pawn, position)
	obj.AddSideEffect(sideEffect)
	assert.Equal(t, []Action{sideEffect}, obj.SideEffects())
}

func TestMoveMergedActions(t *testing.T) {
	pawn1 := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	action1 := NewAction(MoveToPosition, pawn1, position1)

	pawn2 := NewPawn(Red, 0)
	position2 := NewPosition(false, false, nil, nil)
	action2 := NewAction(MoveToStart, pawn2, position2)

	position3 := NewPosition(false, false, nil, nil)
	action3 := NewAction(MoveToStart, nil, position3)

	pawn4 := NewPawn(Red, 0)
	action4 := NewAction(MoveToStart, pawn4, nil)

	card := NewCard("1", Card1)
	actions := []Action{action1, action2}
	sideEffects := []Action{action3, action4}
	expected := []Action{action1, action2, action3, action4}
	obj := NewMove(card, actions, sideEffects)
	assert.Equal(t, expected, obj.MergedActions())
}
