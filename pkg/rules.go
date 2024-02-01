package pkg

import (
	"errors"
	"github.com/pronovic/go-apologies/pkg/util/enum"
	"github.com/google/uuid"
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
type Action interface {

	// Type The type of the action
	Type() ActionType

	// Pawn the pawn that the action operates on
	Pawn() Pawn

	// Position a position that the pawn should move to (optional)
	Position() Position // optional

	// SetPosition Set the position on the action (can be nil)
	SetPosition(position Position)

	// Equals Checks for value equality on the interface
	Equals(other Action) bool
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
		a.pawn.Equals(other.Pawn()) &&
		((a.position == nil && other.Position() == nil) || a.position.Equals(other.Position()))
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
}

type move struct {
	id string
	card Card
	actions []Action
	sideEffects []Action
}

func NewMove(card Card, actions []Action, sideEffects []Action) Move {
	return &move{
		id: uuid.New().String(),
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

// constructLegalMoves Return the set of legal moves for a pawn using a card, possibly empty.
func constructLegalMoves(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	var moves []Move
	if pawn.Position().Home() {
		moves = make([]Move, 0)
	} else {
		switch card.Type() {
		case Card1:
			moves = constructLegalMoves1(color, card, pawn, allPawns)
		case Card2:
			moves = constructLegalMoves2(color, card, pawn, allPawns)
		case Card3:
			moves = constructLegalMoves3(color, card, pawn, allPawns)
		case Card4:
			moves = constructLegalMoves4(color, card, pawn, allPawns)
		case Card5:
			moves = constructLegalMoves5(color, card, pawn, allPawns)
		case Card7:
			moves = constructLegalMoves7(color, card, pawn, allPawns)
		case Card8:
			moves = constructLegalMoves8(color, card, pawn, allPawns)
		case Card10:
			moves = constructLegalMoves10(color, card, pawn, allPawns)
		case Card11:
			moves = constructLegalMoves11(color, card, pawn, allPawns)
		case Card12:
			moves = constructLegalMoves12(color, card, pawn, allPawns)
		case CardApologies:
			moves = constructLegalMovesApologies(color, card, pawn, allPawns)
		}
	}
	augmentWithSlides(allPawns, moves)
	return moves
}

// distanceToHome Return the distance to home for this pawn, a number of squares when moving forward.
func distanceToHome(pawn Pawn) int {
	if pawn.Position().Home() {
		return 0
	} else if pawn.Position().Start() {
		return 65
	} else if pawn.Position().Safe() != nil {
		return SafeSquares - *pawn.Position().Safe()
	} else {
		circle := *StartCircles[pawn.Color()].Square()
		turn := *TurnSquares[pawn.Color()].Square()
		square := *pawn.Position().Square()
		squareToCorner := BoardSquares - square
		cornerToTurn := turn
		turnToHome := SafeSquares + 1
		total := squareToCorner + cornerToTurn + turnToHome
		if turn < square && square < circle {
			return total
		} else {
			if total < 65 {
				return total
			} else {
				return total - 60
			}
		}
	}
}

// Calculate the new position for a forward or backwards move, taking into account safe zone turns but disregarding slides.
func calculatePosition(color PlayerColor, position Position, squares int) (Position, error) {
	if position.Home() || position.Start() {
		return (Position)(nil), errors.New("pawn in home or start may not move")
	} else if position.Safe() != nil {
		if squares == 0 {
			return position.Copy(), nil
		} else if squares > 0 {
			if *position.Safe() + squares < SafeSquares {
				copied := position.Copy()
				err := copied.MoveToSafe(*position.Safe() + squares)
				if err != nil {
					return (Position)(nil), err
				}
				return copied, nil
			} else if *position.Safe() + squares == SafeSquares {
				copied := position.Copy()
				err := copied.MoveToHome()
				if err != nil {
					return (Position)(nil), err
				}
				return copied, nil
			} else {
				return (Position)(nil), errors.New("pawn cannot move past home")
			}
		} else { // squares < 0
			if *position.Safe() + squares >= 0 {
				copied := position.Copy()
				err := copied.MoveToSafe(*position.Safe() + squares)
				if err != nil {
					return (Position)(nil), err
				}
				return copied, nil
			} else {  // handle moving back out of the safe area
				copied := position.Copy()
				err := copied.MoveToSquare(*TurnSquares[color].Square())
				if err != nil {
					return (Position)(nil), err
				}
				return calculatePosition(color, copied, squares + *position.Safe() + 1)
			}
		}
	} else if position.Square() != nil {
		if squares == 0 {
			return position.Copy(), nil
		} else if squares > 0 {
			if *position.Square() + squares < BoardSquares {
				if *position.Square() <= *TurnSquares[color].Square() && *position.Square() + squares > *TurnSquares[color].Square() {
					copied := position.Copy()
					err := copied.MoveToSafe(0)
					if err != nil {
						return (Position)(nil), err
					}
					return calculatePosition(color, copied, squares - (*TurnSquares[color].Square() - *position.Square()) - 1)
				} else {
					copied := position.Copy()
					err := copied.MoveToSquare(*position.Square() + squares)
					if err != nil {
						return (Position)(nil), err
					}
					return copied, nil
				}
			} else { // handle turning the corner
				copied := position.Copy()
				err := copied.MoveToSquare(0)
				if err != nil {
					return (Position)(nil), err
				}
				return calculatePosition(color, copied, squares - (BoardSquares - *position.Square()))
			}
		} else { // squares < 0
			if *position.Square() + squares >= 0 {
				copied := position.Copy()
				err := copied.MoveToSquare(*position.Square() + squares)
				if err != nil {
					return (Position)(nil), err
				}
				return copied, nil
			} else { // handle turning the corner
				copied := position.Copy()
				err := copied.MoveToSquare(BoardSquares - 1)
				if err != nil {
					return (Position)(nil), err
				}
				return calculatePosition(color, copied, squares + *position.Square() + 1)
			}
		}
	} else {
		return (Position)(nil), errors.New("position is in an illegal state")
	}
}

// Return the set of legal moves for a pawn using Card1, possibly empty.
func constructLegalMoves1(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveCircle(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 1)
	return moves
}

// Return the set of legal moves for a pawn using Card2, possibly empty.
func constructLegalMoves2(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveCircle(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 1)
	return moves
}

// Return the set of legal moves for a pawn using Card3, possibly empty.
func constructLegalMoves3(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 3)
	return moves
}

// Return the set of legal moves for a pawn using Card4, possibly empty.
func constructLegalMoves4(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, -4)
	return moves
}

// Return the set of legal moves for a pawn using Card5, possibly empty.
func constructLegalMoves5(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 5)
	return moves
}

// Return the set of legal moves for a pawn using Card7, possibly empty.
func constructLegalMoves7(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 7)
	moveSplit(&moves, color, card, pawn, allPawns)
	return moves
}

// Return the set of legal moves for a pawn using Card8, possibly empty.
func constructLegalMoves8(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 8)
	return moves
}

// Return the set of legal moves for a pawn using Card10, possibly empty.
func constructLegalMoves10(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 10)
	moveSimple(&moves, color, card, pawn, allPawns, -1)
	return moves
}

// Return the set of legal moves for a pawn using Card11, possibly empty.
func constructLegalMoves11(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSwap(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 11)
	return moves
}

// Return the set of legal moves for a pawn using Card12, possibly empty.
func constructLegalMoves12(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 12)
	return moves
}

// Return the set of legal moves for a pawn using CardApologies, possibly empty.
func constructLegalMovesApologies(color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) []Move {
	moves := make([]Move, 0)
	moveApologies(&moves, color, card, pawn, allPawns)
	return moves
}

// Return the first pawn at the indicated position, or None.
func findPawn(allPawns []Pawn, position Position) Pawn {
	for _, p := range allPawns {
		if p.Position().Equals(position) {
			return p
		}
	}

	return nil
}

func moveCircle(moves *[]Move, color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) {
	// For start-related cards, a pawn in the start area can move to the associated
	// circle position if that position is not occupied by another pawn of the same color.
	if pawn.Position().Start() {
		conflict := findPawn(allPawns, StartCircles[color])
		if conflict == nil {
			actions := []Action { NewAction(MoveToPosition, pawn, StartCircles[color].Copy()) }
			sideEffects := make([]Action, 0)
			move := NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		} else if conflict != nil && conflict.Color() != color {
			actions := []Action { NewAction(MoveToPosition, pawn, StartCircles[color].Copy())}
			sideEffects := []Action { NewAction(MoveToStart, conflict, nil) }
			move := NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		}
	}
}

func moveSimple(moves *[]Move, color PlayerColor, card Card, pawn Pawn, allPawns []Pawn, squares int) {
	// For most cards, a pawn on the board can move forward or backward if the
	// resulting position is not occupied by another pawn of the same color.
	if pawn.Position().Square() != nil || pawn.Position().Safe() != nil {
		target, err := calculatePosition(color, pawn.Position(), squares)
		if err == nil { // if the requested position is not legal, then just ignore it
			if target.Home() || target.Start() { // by definition, there can't be a conflict going to home or start
				actions := []Action { NewAction(MoveToPosition, pawn, target) }
				sideEffects := make([]Action, 0)
				move := NewMove(card, actions, sideEffects)
				*moves = append(*moves, move)
			} else {
				conflict := findPawn(allPawns, target)
				if conflict == nil {
					actions := []Action { NewAction(MoveToPosition, pawn, target) }
					sideEffects := make([]Action, 0)
					move := NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				} else if conflict != nil && conflict.Color() != color {
					actions := []Action { NewAction(MoveToPosition, pawn, target)}
					sideEffects := []Action { NewAction(MoveToStart, conflict, nil) }
					move := NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				}
			}
		}
	}
}

func moveSplit(moves *[]Move, color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) {
	// For the 7 card, we can split up the move between two different pawns.
	// Any combination of 7 forward moves is legal, as long as the resulting position
	// is not occupied by another pawn of the same color.

	for _, other := range allPawns {
		if !other.Equals(pawn) && other.Color() == color && !other.Position().Home() && !other.Position().Start() {

			// any pawn except other
			filtered := make([]Pawn, 0)
			for _, p := range allPawns {
				if !p.Equals(other) {
					filtered = append(filtered, p)
				}
			}

			for _, legal := range LegalSplits {
				left := make([]Move, 0)
				moveSimple(&left, color, card, pawn, filtered, legal.Left())

				right := make([]Move, 0)
				moveSimple(&right, color, card, pawn, filtered, legal.Right())

				if len(left) > 0 && len(right) > 0 {
					actions := make([]Action, 0)
					sideEffects := make([]Action, 0)

					for _, l := range left[0].Actions() {
						actions = append(actions, l)
					}

					for _, l := range left[0].SideEffects() {
						sideEffects = append(sideEffects, l)
					}

					for _, r := range right[0].Actions() {
						actions = append(actions, r)
					}

					for _, r := range right[0].SideEffects() {
						sideEffects = append(sideEffects, r)
					}

					move := NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				}
			}
		}
	}
}

func moveSwap(moves *[]Move, color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) {
	// For the 11 card, a pawn on the board can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	if pawn.Position().Square() != nil { // pawn is on the board
		 for _, swap := range allPawns {
			 if swap.Color() != color && !swap.Position().Home() && !swap.Position().Start() && swap.Position().Safe() == nil {
				 actions := []Action {
					 NewAction(MoveToPosition, pawn, swap.Position().Copy()),
					 NewAction(MoveToPosition, swap, pawn.Position().Copy()),
				 }
				 sideEffects := make([]Action, 0)
				 move := NewMove(card, actions, sideEffects)
				 *moves = append(*moves, move)
			 }
		 }
	}
}

func moveApologies(moves *[]Move, color PlayerColor, card Card, pawn Pawn, allPawns []Pawn) {
	// For the Apologies card, a pawn in start can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	if pawn.Position().Start() {
		for _, swap := range allPawns {
			if swap.Color() != color && !swap.Position().Home() && !swap.Position().Start() && swap.Position().Safe() == nil {
				actions := []Action {
					NewAction(MoveToPosition, pawn, swap.Position().Copy()),
					NewAction(MoveToStart, swap, nil),
				}
				sideEffects := make([]Action, 0)
				move := NewMove(card, actions, sideEffects)
				*moves = append(*moves, move)
			}
		}
	}
}

// Augment any legal moves with additional side-effects that occur as a result of slides.
func augmentWithSlides(allPawns []Pawn, moves []Move) {
	for _, move := range moves {
		for _, action := range move.Actions() {
			if action.Type() == MoveToPosition { // look at any move to a position on the board
				for _, color := range PlayerColors.Members() {
					if color != action.Pawn().Color() { // any color other than the pawn's
						for _, slide := range Slides[color] { // # look at all slides with this color
							if action.Position() != nil && action.Position().Square() != nil && *action.Position().Square() == slide.Start() {
								_ = action.Position().MoveToSquare(slide.End()) // if the pawn landed on the start of the slide, move the pawn to the end of the slide
								for square := slide.Start()+1; square <= slide.End(); square++ {
									// Note: in this one case, a pawn can bump another pawn of the same color
									tmp := NewPosition(false, false, nil, &square)
									pawn := findPawn(allPawns, tmp)
									if pawn != nil {
										bump := NewAction(MoveToStart, pawn, nil)
										move.AddSideEffect(bump)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}