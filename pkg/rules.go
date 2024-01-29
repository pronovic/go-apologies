package pkg

import (
	"github.com/pronovic/go-apologies/pkg/util/enum"
)

// ActionType defines all actions that a character can take
type ActionType struct{ value string }

// Value implements the enum.Enum interface for ActionType
func (e ActionType) Value() string { return e.value }

// MoveToStart move a pawn back to its start area
var MoveToStart = ActionType{"MoveToStart"}

// MoveToPosition move a pawn to a specific position on the board
var MoveToPosition = ActionType{"MoveToPosition"}

// ActionTypes is the list of all legal ActionType enumerations
var ActionTypes = enum.NewValues[ActionType](MoveToStart, MoveToPosition)

// Action is an action that can be taken as part of a move
// TODO: finish implementing Action
type Action interface {
	Type() ActionType
	Pawn() Pawn
	Position() *Position // optional
	SetPosition(*Position)
}

// Move is a player's move on the board, which consists of one or more actions
//
// Note that the actions associated with a move include both the immediate actions that the player
// chose (such as moving a pawn from start or swapping places with a different pawn), but also
// any side-effects (such as pawns that are bumped back to start because of a slide).  As a result,
// executing a move becomes very easy and no validation is required.  All of the work is done
// up-front.
// TODO: finish implementing Move
type Move struct {
}

// ConstructLegalMoves Return the set of legal moves for a pawn using a card, possibly empty.
func ConstructLegalMoves(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement ConstructLegalMoves
}

// DistanceToHome Return the distance to home for this pawn, a number of squares when moving forward.
func DistanceToHome(pawn Pawn) int {
	return 0  // TODO: implement DistanceToHome
}

// Calculate the new position for a forward or backwards move, taking into account safe zone turns but disregarding slides.
func newPosition(color PlayerColor, position Position, squares int) Position {
	return *new(Position) // TODO: implement position
}

// Return the set of legal moves for a pawn using Card1, possibly empty.
func constructLegalMoves1(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves1
}

// Return the set of legal moves for a pawn using Card2, possibly empty.
func constructLegalMoves2(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves2
}

// Return the set of legal moves for a pawn using Card3, possibly empty.
func constructLegalMoves3(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves3
}

// Return the set of legal moves for a pawn using Card4, possibly empty.
func constructLegalMoves4(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves4
}

// Return the set of legal moves for a pawn using Card5, possibly empty.
func constructLegalMoves5(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves5
}

// Return the set of legal moves for a pawn using Card7, possibly empty.
func constructLegalMoves7(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves7
}

// Return the set of legal moves for a pawn using Card8, possibly empty.
func constructLegalMoves8(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves8
}

// Return the set of legal moves for a pawn using Card10, possibly empty.
func constructLegalMoves10(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves10
}

// Return the set of legal moves for a pawn using Card11, possibly empty.
func constructLegalMoves11(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves11
}

// Return the set of legal moves for a pawn using Card12, possibly empty.
func constructLegalMoves12(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMoves12
}

// Return the set of legal moves for a pawn using CardApologies, possibly empty.
func constructLegalMovesApologies(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	return []Move{} // TODO: implement constructLegalMovesApologies
}

// Return the first pawn at the indicated position, or None.
func findPawn(allPawns []Pawn, position Position) *Pawn {
	return nil  // TODO: implement findPawn
}

func moveCircle(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	// For start-related cards, a pawn in the start area can move to the associated
	// circle position if that position is not occupied by another pawn of the same color.
	return []Move{} // TODO: implement moveCircle
}

func moveSimple(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	// For most cards, a pawn on the board can move forward or backward if the
	// resulting position is not occupied by another pawn of the same color.
	return []Move{} // TODO: implement moveSimple
}

func moveSplit(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	// For the 7 card, we can split up the move between two different pawns.
	// Any combination of 7 forward moves is legal, as long as the resulting position
	// is not occupied by another pawn of the same color.
	return []Move{} // TODO: implement moveSplit
}

func moveSwap(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	// For the 11 card, a pawn on the board can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	return []Move{} // TODO: implement moveSwap
}

func moveApologies(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	// For the Apologies card, a pawn in start can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	return []Move{} // TODO: implement moveApologies
}

// Augment any legal moves with additional side-effects that occur as a result of slides.
func augmentWithSlides(allPawns []Pawn, moves []Move) {
	// TODO: implement augmentWithSlides
}