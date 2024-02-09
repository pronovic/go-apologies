package model

import (
	"bytes"
	"encoding/json"
	"github.com/pronovic/go-apologies/internal/enum"
	"github.com/pronovic/go-apologies/internal/equality"
	"github.com/pronovic/go-apologies/internal/identifier"
	"io"
)

// ActionType defines all actions that a character can take
type ActionType struct{ value string }
func (e ActionType) Value() string { return e.value }
func (e ActionType) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
func (e *ActionType) UnmarshalText(text []byte) error { return enum.Unmarshal(e, text, ActionTypes) }
var ActionTypes = enum.NewValues[ActionType](MoveToStart, MoveToPosition)
var MoveToStart = ActionType{"MoveToStart"}
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
	XactionType ActionType	`json:"type"`
	Xpawn     Pawn `json:"pawn"`
	Xposition Position `json:"position"`
}

// NewAction constructs a new Action
func NewAction(actionType ActionType, pawn Pawn, position Position) Action {
	return &action{
		XactionType: actionType,
		Xpawn:       pawn,
		Xposition:   position,
	}
}

// NewActionFromJSON constructs a new object from JSON in an io.Reader
func NewActionFromJSON(reader io.Reader) (Action, error) {
	type raw struct {
		XactionType ActionType	`json:"type"`
		Xpawn     json.RawMessage `json:"pawn"`
		Xposition json.RawMessage `json:"position"`
	}

	var temp raw
	err := json.NewDecoder(reader).Decode(&temp)
	if err != nil {
		return nil, err
	}

	var Xpawn Pawn
	if temp.Xpawn != nil {
		Xpawn, err = NewPawnFromJSON(bytes.NewReader(temp.Xpawn))
		if err != nil {
			return nil, err
		}
	}

	var Xposition Position
	if temp.Xposition != nil {
		Xposition, err = NewPositionFromJSON(bytes.NewReader(temp.Xposition))
		if err != nil {
			return nil, err
		}
	}

	obj := action {
		XactionType: temp.XactionType,
		Xpawn: Xpawn,
		Xposition: Xposition,
	}

	return &obj, nil
}

func (a *action) Type() ActionType {
	return a.XactionType
}

func (a *action) Pawn() Pawn {
	return a.Xpawn
}

func (a *action) Position() Position {
	return a.Xposition
}

func (a *action) SetPosition(position Position) {
	a.Xposition = position
}

func (a *action) Equals(other Action) bool {
	return other != nil &&
		a.XactionType == other.Type() &&
		equality.ByValueEquals[Pawn](a.Xpawn, other.Pawn()) &&
		equality.ByValueEquals[Position](a.Xposition, other.Position())
}

// Move is a player's move on the board, which consists of one or more actions
//
// Note that the actions associated with a move include both the immediate actions that the player
// chose (such as moving a pawn from start or swapping places with a different pawn), but also
// any side-effects (such as pawns that are bumped back to start because of a slide).  As a result,
// executing a move becomes very easy and no validation is required.  All of the work is done
// up-front.
type Move interface {
	equality.EqualsByValue[Move]  // This interface implements equality by value
	Id() string
	Card() Card
	Actions() []Action
	SideEffects() []Action
	AddSideEffect(action Action)
	MergedActions() []Action
}

type move struct {
	Xid     string
	Xcard       Card
	Xactions     []Action
	XsideEffects []Action
}

// NewMove constructs a new move, optionally accepting an identify factory
// If there are no actions or side effects, you may pass nil, which is equivalent to a newly-allocated empty slice
func NewMove(card Card, actions []Action, sideEffects []Action, factory identifier.Factory) Move {
	if factory == nil {
		factory = identifier.NewFactory()
	}

	if actions == nil {
		actions = make([]Action, 0)
	}

	if sideEffects == nil {
		sideEffects = make([]Action, 0)
	}

	return &move{
		Xid:          factory.RandomId(),
		Xcard:        card,
		Xactions:     actions,
		XsideEffects: sideEffects,
	}
}

// NewMoveFromJSON constructs a new object from JSON in an io.Reader
func NewMoveFromJSON(reader io.Reader) (Move, error) {
	return nil, nil // TODO: implement NewMoveFromJSON
}

func (m *move) Equals(other Move) bool {
	// note that identifier is not included in eqa
	return other != nil &&
		equality.ByValueEquals[Card](m.Xcard, other.Card()) &&
		equality.SliceByValueEquals(m.Xactions, other.Actions()) &&
		equality.SliceByValueEquals(m.XsideEffects, other.SideEffects())
}

func (m *move) Id() string {
	return m.Xid
}

func (m *move) Card() Card {
	return m.Xcard
}

func (m *move) Actions() []Action {
	return m.Xactions
}

func (m *move) SideEffects() []Action {
	return m.XsideEffects
}

func (m *move) AddSideEffect(action Action) {
	found := false
	for _, a := range m.Xactions {
		if a.Equals(action) {
			found = true
			break
		}
	}

	if !found {
		m.XsideEffects = append(m.XsideEffects, action)
	}
}

func (m *move) MergedActions() []Action {
	merged := make([]Action, 0, len(m.Xactions) + len(m.XsideEffects))

	for _, action := range m.Xactions {
		merged = append(merged, action)
	}

	for _, action := range m.XsideEffects {
		merged = append(merged, action)
	}

	return merged
}
