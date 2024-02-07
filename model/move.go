package model

import (
	"github.com/pronovic/go-apologies/internal/equality"
	"github.com/pronovic/go-apologies/internal/identifier"
)

// ActionType defines all actions that a character can take
type ActionType struct{ value string }

// Value implements the enum.Enum interface for ActionType
func (e ActionType) Value() string { return e.value }

// MoveToStart move a pawn back to its start area
var MoveToStart = ActionType{"MoveToStart"}

// MoveToPosition move a pawn to a specific position on the board
var MoveToPosition = ActionType{"MoveToPosition"}

// Action is an action that can be taken as part of a move
type Action interface {

	equality.EqualsByValue[Action]  // This interface implements equality by value

	// Type The type of the action
	Type() ActionType

	// Pawn the pawn that the action operates on
	Pawn() Pawn

	// Position a position that the pawn should move to (optional)
	Position() Position // optional

	// SetPosition Set the position on the action (can be nil)
	SetPosition(position Position)
}

type action struct {
	actionType ActionType
	pawn Pawn
	position Position
}

// NewAction constructs a new Action
func NewAction(actionType ActionType, pawn Pawn, position Position) Action {
	return &action{
		actionType: actionType,
		pawn: pawn,
		position: position,
	}
}

func (a *action) Type() ActionType {
	return a.actionType
}

func (a *action) Pawn() Pawn {
	return a.pawn
}

func (a *action) Position() Position {
	return a.position
}

func (a *action) SetPosition(position Position) {
	a.position = position
}

func (a *action) Equals(other Action) bool {
	return a.actionType == other.Type() &&
		equality.ByValueEquals[Pawn](a.pawn, other.Pawn()) &&
		equality.ByValueEquals[Position](a.position, other.Position())
}

// Move is a player's move on the board, which consists of one or more actions
//
// Note that the actions associated with a move include both the immediate actions that the player
// chose (such as moving a pawn from start or swapping places with a different pawn), but also
// any side-effects (such as pawns that are bumped back to start because of a slide).  As a result,
// executing a move becomes very easy and no validation is required.  All of the work is done
// up-front.
type Move interface {
	Id() string
	Card() Card
	Actions() []Action
	SideEffects() []Action
	AddSideEffect(action Action)
	MergedActions() []Action
}

type move struct {
	id string
	card Card
	actions []Action
	sideEffects []Action
}

// NewMove constructs a new move
func NewMove(card Card, actions []Action, sideEffects []Action) Move {
	return newMove(card, actions, sideEffects, identifier.NewFactory())
}

// newMove constructs a new move while accepting an identifier factory (intended for unit testing)
func newMove(card Card, actions []Action, sideEffects []Action, factory identifier.Factory) Move {
	return &move{
		id: factory.RandomId(),
		card: card,
		actions: actions,
		sideEffects: sideEffects,
	}
}

func (m *move) Id() string {
	return m.id
}

func (m *move) Card() Card {
	return m.card
}

func (m *move) Actions() []Action {
	return m.actions
}

func (m *move) SideEffects() []Action {
	return m.sideEffects
}

func (m *move) AddSideEffect(action Action) {
	found := false
	for _, a := range m.actions {
		if a.Equals(action) {
			found = true
			break
		}
	}

	if !found {
		m.sideEffects = append(m.sideEffects, action)
	}
}

func (m *move) MergedActions() []Action {
	merged := make([]Action, 0, len(m.actions) + len(m.sideEffects))

	for _, action := range m.actions {
		merged = append(merged, action)
	}

	for _, action := range m.sideEffects {
		merged = append(merged, action)
	}

	return merged
}
