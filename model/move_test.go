package model

import (
	"github.com/pronovic/go-apologies/internal/identifier"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	identifier.UseStubbedId()   // once this has been called, it takes effect permanently for all unit tests
}

func TestNewAction(t *testing.T) {
	pawn := NewPawn(Red, 0)
	position := NewPosition(false, false, nil, nil)
	obj := NewAction(MoveToPosition, pawn, position)
	assert.Equal(t, MoveToPosition, obj.Type())
	assert.Same(t, pawn, obj.Pawn())
	assert.Same(t, position, obj.Position())
}

func TestActionSetPosition(t *testing.T) {
	pawn := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	position2 := NewPosition(true, false, nil, nil)
	obj := NewAction(MoveToPosition, pawn, position1)
	obj.SetPosition(position2)
	assert.Same(t, position2, obj.Position())
}

func TestActionEquals(t *testing.T) {
	pawn1 := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	obj1 := NewAction(MoveToPosition, pawn1, position1)

	pawn2 := NewPawn(Red, 0)
	position2 := NewPosition(false, false, nil, nil)
	obj2 := NewAction(MoveToStart, pawn2, position2)

	position3 := NewPosition(false, false, nil, nil)
	obj3 := NewAction(MoveToStart, nil, position3)

	pawn4 := NewPawn(Red, 0)
	obj4 := NewAction(MoveToStart, pawn4, nil)

	obj5 := NewAction(MoveToPosition, nil, nil)

	assert.Equal(t, obj1, obj1)
	assert.Equal(t, obj2, obj2)
	assert.Equal(t, obj3, obj3)
	assert.Equal(t, obj4, obj4)
	assert.Equal(t, obj5, obj5)

	assert.NotEqual(t, obj1, nil)
	assert.NotEqual(t, obj1, nil)
	assert.NotEqual(t, obj1, nil)
	assert.NotEqual(t, obj1, nil)

	assert.NotEqual(t, obj1, obj2)
	assert.NotEqual(t, obj1, obj3)
	assert.NotEqual(t, obj1, obj4)
	assert.NotEqual(t, obj1, obj5)

	assert.NotEqual(t, obj5, obj1)
	assert.NotEqual(t, obj5, obj2)
	assert.NotEqual(t, obj5, obj3)
	assert.NotEqual(t, obj5, obj4)
}

func TestNewMove(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 1, 2)
	sideEffects := make([]Action, 2, 3)
	obj := NewMove(card, actions, sideEffects)
	assert.NotEmptyf(t, identifier.StubbedId, obj.Id()) // filled in with a UUID
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())
	assert.Equal(t, sideEffects, obj.SideEffects())
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
	assert.Equal(t, []Action {sideEffect}, obj.SideEffects())
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
	actions := []Action { action1, action2 }
	sideEffects := []Action { action3, action4 }
	expected := []Action { action1, action2, action3, action4 }
	obj := NewMove(card, actions, sideEffects)
	assert.Equal(t, expected, obj.MergedActions())
}