package model

import (
	"encoding/json"
	"io"

	"github.com/pronovic/go-apologies/internal/enum"
	"github.com/pronovic/go-apologies/internal/equality"
	"github.com/pronovic/go-apologies/internal/jsonutil"
)

// ActionType defines all actions that a character can take
type ActionType struct{ value string }

func (e ActionType) Value() string                         { return e.value }
func (e ActionType) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
func (e *ActionType) UnmarshalText(text []byte) error      { return enum.Unmarshal(e, text, ActionTypes) }

var (
	ActionTypes    = enum.NewValues[ActionType](MoveToStart, MoveToPosition)
	MoveToStart    = ActionType{"MoveToStart"}
	MoveToPosition = ActionType{"MoveToPosition"}
)

// Action is an action that can be taken as part of a move
type Action interface {
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
	XactionType ActionType `json:"type"`
	Xpawn       Pawn       `json:"pawn"`
	Xposition   Position   `json:"position"`
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
		XactionType ActionType      `json:"type"`
		Xpawn       json.RawMessage `json:"pawn"`
		Xposition   json.RawMessage `json:"position"`
	}

	var temp raw
	err := json.NewDecoder(reader).Decode(&temp)
	if err != nil {
		return nil, err
	}

	var Xpawn Pawn
	Xpawn, err = jsonutil.DecodeInterfaceJSON(temp.Xpawn, NewPawnFromJSON)
	if err != nil {
		return nil, err
	}

	var Xposition Position
	Xposition, err = jsonutil.DecodeInterfaceJSON(temp.Xposition, NewPositionFromJSON)
	if err != nil {
		return nil, err
	}

	obj := action{
		XactionType: temp.XactionType,
		Xpawn:       Xpawn,
		Xposition:   Xposition,
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

// Move is a player's move on the board, which consists of one or more actions
//
// Note that the actions associated with a move include both the immediate actions that the player
// chose (such as moving a pawn from start or swapping places with a different pawn), but also
// any side-effects (such as pawns that are bumped back to start because of a slide).  As a result,
// executing a move becomes very easy and no validation is required.  All of the work is done
// up-front.
type Move interface {
	Card() Card
	Actions() []Action
	SideEffects() []Action
	AddSideEffect(action Action)
	MergedActions() []Action
}

type move struct {
	Xid          string   `json:"id"`
	Xcard        Card     `json:"card"`
	Xactions     []Action `json:"actions"`
	XsideEffects []Action `json:"sideeffects"`
}

// NewMove constructs a new move, optionally accepting an identify factory
// If there are no actions or side effects, you may pass nil, which is equivalent to a newly-allocated empty slice
func NewMove(card Card, actions []Action, sideEffects []Action) Move {
	if actions == nil {
		actions = make([]Action, 0)
	}

	if sideEffects == nil {
		sideEffects = make([]Action, 0)
	}

	return &move{
		Xcard:        card,
		Xactions:     actions,
		XsideEffects: sideEffects,
	}
}

// NewMoveFromJSON constructs a new object from JSON in an io.Reader
func NewMoveFromJSON(reader io.Reader) (Move, error) {
	type raw struct {
		Xid          string            `json:"id"`
		Xcard        json.RawMessage   `json:"card"`
		Xactions     []json.RawMessage `json:"actions"`
		XsideEffects []json.RawMessage `json:"sideeffects"`
	}

	var temp raw
	err := json.NewDecoder(reader).Decode(&temp)
	if err != nil {
		return nil, err
	}

	var Xcard Card
	Xcard, err = jsonutil.DecodeInterfaceJSON(temp.Xcard, NewCardFromJSON)
	if err != nil {
		return nil, err
	}

	var Xactions []Action
	Xactions, err = jsonutil.DecodeSliceJSON(temp.Xactions, NewActionFromJSON)
	if err != nil {
		return nil, err
	}

	var XsideEffects []Action
	XsideEffects, err = jsonutil.DecodeSliceJSON(temp.XsideEffects, NewActionFromJSON)
	if err != nil {
		return nil, err
	}

	obj := move{
		Xid:          temp.Xid,
		Xcard:        Xcard,
		Xactions:     Xactions,
		XsideEffects: XsideEffects,
	}

	return &obj, nil
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
		if equality.EqualByValue(a, action) {
			found = true
			break
		}
	}

	if !found {
		m.XsideEffects = append(m.XsideEffects, action)
	}
}

func (m *move) MergedActions() []Action {
	merged := make([]Action, 0, len(m.Xactions)+len(m.XsideEffects))

	for _, action := range m.Xactions {
		merged = append(merged, action)
	}

	for _, action := range m.XsideEffects {
		merged = append(merged, action)
	}

	return merged
}
