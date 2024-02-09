package model

import (
	"github.com/pronovic/go-apologies/internal/identifier"
	"github.com/stretchr/testify/assert"
	"testing"
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
	t.Fail() // TODO: implement TestNewActionFromJSON()
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

	// note: it is important to test with assert.True()/assert.False() and x.Equals(y)
	// because assert.Equals() and assert.NotEquals() are not aware of our equality by value concept

	assert.True(t, obj1.Equals(obj1))
	assert.True(t, obj2.Equals(obj2))
	assert.True(t, obj3.Equals(obj3))
	assert.True(t, obj4.Equals(obj4))
	assert.True(t, obj5.Equals(obj5))

	assert.False(t, obj1.Equals(nil))
	assert.False(t, obj1.Equals( nil))
	assert.False(t, obj1.Equals( nil))
	assert.False(t, obj1.Equals( nil))

	assert.False(t, obj1.Equals(obj2))
	assert.False(t, obj1.Equals(obj3))
	assert.False(t, obj1.Equals(obj4))
	assert.False(t, obj1.Equals(obj5))

	assert.False(t, obj5.Equals(obj1))
	assert.False(t, obj5.Equals(obj2))
	assert.False(t, obj5.Equals(obj3))
	assert.False(t, obj5.Equals(obj4))
}

func TestNewMove(t *testing.T) {
	var factory identifier.MockFactory
	factory.On("RandomId").Return("id")
	card := NewCard("1", Card1)
	actions := make([]Action, 1, 2)
	sideEffects := make([]Action, 2, 3)
	obj := NewMove(card, actions, sideEffects, &factory)
	assert.Equal(t, "id", obj.Id()) // filled in with the constant id from the mock factory
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())
	assert.Equal(t, sideEffects, obj.SideEffects())
}

func TestNewMoveFromJSON(t *testing.T) {
	t.Fail() // TODO: implement TestNewMoveFromJSON()
}

func TestMoveEquals(t *testing.T) {
	card1 := NewCard("1", Card1)
	card2 := NewCard("2", Card2)
	card3 := NewCard("3", Card3)
	card4 := NewCard("4", Card4)

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

	obj1 := NewMove(card1, nil, nil, nil)
	obj2 := NewMove(card1, nil, nil, nil)  // note: equivalent to obj1, but will have a different id
	obj3 := NewMove(card2, []Action { action1 }, nil, nil)
	obj4 := NewMove(card3, nil, []Action { action2 }, nil)
	obj5 := NewMove(card4, []Action { action3 }, []Action { action4 }, nil)
	obj6 := NewMove(card4, []Action { action3 }, []Action { action4 }, nil)  // note: equivalent to obj6, but will have different id

	// note: it is important to test with assert.True()/assert.False() and x.Equals(y)
	// because assert.Equals() and assert.NotEquals() are not aware of our equality by value concept

	assert.False(t, obj1.Equals(nil))
	assert.True(t, obj1.Equals(obj1))
	assert.True(t, obj1.Equals(obj2))
	assert.False(t, obj1.Equals(obj3))
	assert.False(t, obj1.Equals(obj4))
	assert.False(t, obj1.Equals(obj5))
	assert.False(t, obj1.Equals(obj6))

	assert.False(t, obj2.Equals(nil))
	assert.True(t, obj2.Equals(obj1))
	assert.True(t, obj2.Equals(obj2))
	assert.False(t, obj2.Equals(obj3))
	assert.False(t, obj2.Equals(obj4))
	assert.False(t, obj2.Equals(obj5))
	assert.False(t, obj2.Equals(obj6))

	assert.False(t, obj3.Equals(nil))
	assert.False(t, obj3.Equals(obj1))
	assert.False(t, obj3.Equals(obj2))
	assert.True(t, obj3.Equals(obj3))
	assert.False(t, obj3.Equals(obj4))
	assert.False(t, obj3.Equals(obj5))
	assert.False(t, obj3.Equals(obj6))

	assert.False(t, obj4.Equals(nil))
	assert.False(t, obj4.Equals(obj1))
	assert.False(t, obj4.Equals(obj2))
	assert.False(t, obj4.Equals(obj3))
	assert.True(t, obj4.Equals(obj4))
	assert.False(t, obj4.Equals(obj5))
	assert.False(t, obj4.Equals(obj6))

	assert.False(t, obj5.Equals(nil))
	assert.False(t, obj5.Equals(obj1))
	assert.False(t, obj5.Equals(obj2))
	assert.False(t, obj5.Equals(obj3))
	assert.False(t, obj5.Equals(obj4))
	assert.True(t, obj5.Equals(obj5))
	assert.True(t, obj5.Equals(obj6))

	assert.False(t, obj6.Equals(nil))
	assert.False(t, obj6.Equals(obj1))
	assert.False(t, obj6.Equals(obj2))
	assert.False(t, obj6.Equals(obj3))
	assert.False(t, obj6.Equals(obj4))
	assert.True(t, obj6.Equals(obj5))
	assert.True(t, obj6.Equals(obj6))
}

func TestNewMoveEmptySlice(t *testing.T) {
	card := NewCard("1", Card1)
	var actions = make([]Action, 0)
	var sideEffects = make([]Action, 0)
	obj := NewMove(card, actions, sideEffects, nil)
	assert.NotEmpty(t, obj.Id()) // filled in with an id from the standard identifier factory
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())
	assert.Equal(t, sideEffects, obj.SideEffects())
}

func TestNewMoveNilSlice(t *testing.T) {
	card := NewCard("1", Card1)
	var actions = make([]Action, 0)
	var sideEffects = make([]Action, 0)
	obj := NewMove(card, nil, nil, nil)
	assert.NotEmpty(t, obj.Id()) // filled in with an id from the standard identifier factory
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())  // nil is converted to a newly-allocated empty slice
	assert.Equal(t, sideEffects, obj.SideEffects()) // nil is converted to a newly-allocated empty slice
}

func TestMoveAddSideEffect(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 0)
	sideEffects := make([]Action, 0)
	obj := NewMove(card, actions, sideEffects, nil)

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
	obj := NewMove(card, actions, sideEffects, nil)
	assert.Equal(t, expected, obj.MergedActions())
}