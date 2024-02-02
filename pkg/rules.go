package pkg

import (
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/internal/equality"
	"github.com/pronovic/go-apologies/internal/identifier"
	"github.com/pronovic/go-apologies/model"
)

// splitPair defines a legal way to split up a move of 7
type splitPair struct {
	left int
	right int
}

// legalSplits defines legal ways to split up a move of 7
var legalSplits = []splitPair{
	splitPair {1, 6},
	splitPair {2, 5},
	splitPair {3, 4},
	splitPair {4, 3},
	splitPair {5, 2},
	splitPair {6, 1},
}

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
	Pawn() model.Pawn

	// Position a position that the pawn should move to (optional)
	Position() model.Position // optional

	// SetPosition Set the position on the action (can be nil)
	SetPosition(position model.Position)
}

type action struct {
	actionType ActionType
	pawn model.Pawn
	position model.Position
}

// NewAction constructs a new Action
func NewAction(actionType ActionType, pawn model.Pawn, position model.Position) Action {
	return &action{
		actionType: actionType,
		pawn: pawn,
		position: position,
	}
}

func (a *action) Type() ActionType {
	return a.actionType
}

func (a *action) Pawn() model.Pawn {
	return a.pawn
}

func (a *action) Position() model.Position {
	return a.position
}

func (a *action) SetPosition(position model.Position) {
	a.position = position
}

func (a *action) Equals(other Action) bool {
	return a.actionType == other.Type() &&
		equality.ByValueEquals[model.Pawn](a.pawn, other.Pawn()) &&
		equality.ByValueEquals[model.Position](a.position, other.Position())
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
	Card() model.Card
	Actions() []Action
	SideEffects() []Action
	AddSideEffect(action Action)
	MergedActions() []Action
}

type move struct {
	id string
	card model.Card
	actions []Action
	sideEffects []Action
}

func NewMove(card model.Card, actions []Action, sideEffects []Action) Move {
	return &move{
		id: identifier.NewId(),
		card: card,
		actions: actions,
		sideEffects: sideEffects,
	}
}

func (m *move) Id() string {
	return m.id
}

func (m *move) Card() model.Card {
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

// StartGame starts a game using the passed-in mode
func StartGame(game model.Game, mode model.GameMode) error {
	if game.Started() {
		return errors.New("game is already started")
	}

	game.Track(fmt.Sprintf("Game started with mode: %s", mode), nil, nil)

	// the adult mode version of the game moves some pawns and deals some cards to each player
	if mode == model.AdultMode {
		for _, player := range game.Players() {
			err := player.Pawns()[0].Position().MoveToPosition(model.StartCircles[player.Color()])
			if err != nil {
				return err
			}
		}

		for i := 0; i < model.AdultHand; i++ {
			for _, player := range game.Players() {
				card, err := game.Deck().Draw()
				if err != nil {
					return err
				}
				player.AppendToHand(card)
			}
		}
	}

	return nil
}

// ConstructLegalMoves returns the set of all legal moves for a player and its opponents
// Pass the card to play, or nil if the move should come from the player's hand
func ConstructLegalMoves(view model.PlayerView, card model.Card) ([]Move, error) {
	moves := make([]Move, 0)
	allPawns := view.AllPawns()

	cards := view.Player().Hand()
	if card != nil {
		cards = make([]model.Card, 0)
		cards = append(cards, card)
	}

	for _, played := range cards {
		for _, pawn := range view.Player().Pawns() {
			for _, move := range constructLegalMoves(view.Player().Color(), played, pawn, allPawns) {
				moves = append(moves, move)  // TODO: filter out duplicates
			}
		}
	}

	// if there are no legal moves, then forfeit (discarding one card) becomes the only allowable move
	if len(moves) == 0 {
		for _, played := range cards {
			moves = append(moves, NewMove(played, []Action{}, []Action{}))
		}
	}

	if len(moves) == 0 {
		return []Move{}, errors.New("internal error: could not construct any legal moves")
	}

	return moves, nil
}

// ExecuteMove Execute a player's move, updating game state
func ExecuteMove(game model.Game, player model.Player, move Move) error {
	for _, action := range move.MergedActions() { // execute actions, then side effects, in order
		// keep in mind that the pawn on the action is a different object than the pawn in the game
		pawn := game.Players()[action.Pawn().Color()].Pawns()[action.Pawn().Index()]
		if action.Type() == MoveToStart {
			game.Track(fmt.Sprintf("Played card %s: [%s->start]", move.Card().Type().Value(), pawn.Name()), player, move.Card())
			err := pawn.Position().MoveToStart()
			if err != nil {
				return err
			}
		} else if action.Type() == MoveToPosition && action.Position() != nil {
			game.Track(fmt.Sprintf("Played card %s: [%s->position]", move.Card().Type().Value(), pawn.Name()), player, move.Card())
			err := pawn.Position().MoveToPosition(action.Position())
			if err != nil {
				return err
			}
		}
	}

	if game.Completed() {
		winner := *game.Winner()
		game.Track(fmt.Sprintf("Game completed: winner is %s after %d turns", winner.Color().Value(), winner.Turns()), nil, nil)
	}

	return nil
}

// EvaluateMove constructs a new player view that results from executing the passed-in move.
// This is equivalent to execute_move() but has no permanent effect on the game.  It's intended for
// use by a character, to evaluate the results of each legal move.
func EvaluateMove(view model.PlayerView, move Move) (model.PlayerView, error) {
	result := view.Copy()

	for _, action := range move.MergedActions() { // execute actions, then side effects, in order
		// keep in mind that the pawn on the action is a different object than the pawn in the game
		pawn := result.GetPawn(action.Pawn())
		if pawn != nil {  // if the pawn isn't valid, just ignore it
			if action.Type() == MoveToStart {
				err := pawn.Position().MoveToStart()
				if err != nil {
					return nil, err
				}
			} else if action.Type() == MoveToPosition && action.Position() != nil {
				err := pawn.Position().MoveToPosition(action.Position())
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return result, nil
}

// constructLegalMoves Return the set of legal moves for a pawn using a card, possibly empty.
func constructLegalMoves(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	var moves []Move
	if pawn.Position().Home() {
		moves = make([]Move, 0)
	} else {
		switch card.Type() {
		case model.Card1:
			moves = constructLegalMoves1(color, card, pawn, allPawns)
		case model.Card2:
			moves = constructLegalMoves2(color, card, pawn, allPawns)
		case model.Card3:
			moves = constructLegalMoves3(color, card, pawn, allPawns)
		case model.Card4:
			moves = constructLegalMoves4(color, card, pawn, allPawns)
		case model.Card5:
			moves = constructLegalMoves5(color, card, pawn, allPawns)
		case model.Card7:
			moves = constructLegalMoves7(color, card, pawn, allPawns)
		case model.Card8:
			moves = constructLegalMoves8(color, card, pawn, allPawns)
		case model.Card10:
			moves = constructLegalMoves10(color, card, pawn, allPawns)
		case model.Card11:
			moves = constructLegalMoves11(color, card, pawn, allPawns)
		case model.Card12:
			moves = constructLegalMoves12(color, card, pawn, allPawns)
		case model.CardApologies:
			moves = constructLegalMovesApologies(color, card, pawn, allPawns)
		}
	}
	augmentWithSlides(allPawns, moves)
	return moves
}

// distanceToHome Return the distance to home for this pawn, a number of squares when moving forward.
func distanceToHome(pawn model.Pawn) int {
	if pawn.Position().Home() {
		return 0
	} else if pawn.Position().Start() {
		return 65
	} else if pawn.Position().Safe() != nil {
		return model.SafeSquares - *pawn.Position().Safe()
	} else {
		circle := *model.StartCircles[pawn.Color()].Square()
		turn := *model.TurnSquares[pawn.Color()].Square()
		square := *pawn.Position().Square()
		squareToCorner := model.BoardSquares - square
		cornerToTurn := turn
		turnToHome := model.SafeSquares + 1
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

// Calculate the new position for a forward or backwards move, taking into account safe zone turns but disregarding model.Slides.
func calculatePosition(color model.PlayerColor, position model.Position, squares int) (model.Position, error) {
	if position.Home() || position.Start() {
		return (model.Position)(nil), errors.New("pawn in home or start may not move")
	} else if position.Safe() != nil {
		if squares == 0 {
			return position.Copy(), nil
		} else if squares > 0 {
			if *position.Safe() + squares < model.SafeSquares {
				copied := position.Copy()
				err := copied.MoveToSafe(*position.Safe() + squares)
				if err != nil {
					return (model.Position)(nil), err
				}
				return copied, nil
			} else if *position.Safe() + squares == model.SafeSquares {
				copied := position.Copy()
				err := copied.MoveToHome()
				if err != nil {
					return (model.Position)(nil), err
				}
				return copied, nil
			} else {
				return (model.Position)(nil), errors.New("pawn cannot move past home")
			}
		} else { // squares < 0
			if *position.Safe() + squares >= 0 {
				copied := position.Copy()
				err := copied.MoveToSafe(*position.Safe() + squares)
				if err != nil {
					return (model.Position)(nil), err
				}
				return copied, nil
			} else {  // handle moving back out of the safe area
				copied := position.Copy()
				err := copied.MoveToSquare(*model.TurnSquares[color].Square())
				if err != nil {
					return (model.Position)(nil), err
				}
				return calculatePosition(color, copied, squares + *position.Safe() + 1)
			}
		}
	} else if position.Square() != nil {
		if squares == 0 {
			return position.Copy(), nil
		} else if squares > 0 {
			if *position.Square() + squares < model.BoardSquares {
				if *position.Square() <= *model.TurnSquares[color].Square() && *position.Square() + squares > *model.TurnSquares[color].Square() {
					copied := position.Copy()
					err := copied.MoveToSafe(0)
					if err != nil {
						return (model.Position)(nil), err
					}
					return calculatePosition(color, copied, squares - (*model.TurnSquares[color].Square() - *position.Square()) - 1)
				} else {
					copied := position.Copy()
					err := copied.MoveToSquare(*position.Square() + squares)
					if err != nil {
						return (model.Position)(nil), err
					}
					return copied, nil
				}
			} else { // handle turning the corner
				copied := position.Copy()
				err := copied.MoveToSquare(0)
				if err != nil {
					return (model.Position)(nil), err
				}
				return calculatePosition(color, copied, squares - (model.BoardSquares - *position.Square()))
			}
		} else { // squares < 0
			if *position.Square() + squares >= 0 {
				copied := position.Copy()
				err := copied.MoveToSquare(*position.Square() + squares)
				if err != nil {
					return (model.Position)(nil), err
				}
				return copied, nil
			} else { // handle turning the corner
				copied := position.Copy()
				err := copied.MoveToSquare(model.BoardSquares - 1)
				if err != nil {
					return (model.Position)(nil), err
				}
				return calculatePosition(color, copied, squares + *position.Square() + 1)
			}
		}
	} else {
		return (model.Position)(nil), errors.New("position is in an illegal state")
	}
}

// Return the set of legal moves for a pawn using Card1, possibly empty.
func constructLegalMoves1(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveCircle(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 1)
	return moves
}

// Return the set of legal moves for a pawn using Card2, possibly empty.
func constructLegalMoves2(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveCircle(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 2)
	return moves
}

// Return the set of legal moves for a pawn using Card3, possibly empty.
func constructLegalMoves3(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 3)
	return moves
}

// Return the set of legal moves for a pawn using Card4, possibly empty.
func constructLegalMoves4(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, -4)
	return moves
}

// Return the set of legal moves for a pawn using Card5, possibly empty.
func constructLegalMoves5(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 5)
	return moves
}

// Return the set of legal moves for a pawn using Card7, possibly empty.
func constructLegalMoves7(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 7)
	moveSplit(&moves, color, card, pawn, allPawns)
	return moves
}

// Return the set of legal moves for a pawn using Card8, possibly empty.
func constructLegalMoves8(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 8)
	return moves
}

// Return the set of legal moves for a pawn using Card10, possibly empty.
func constructLegalMoves10(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 10)
	moveSimple(&moves, color, card, pawn, allPawns, -1)
	return moves
}

// Return the set of legal moves for a pawn using Card11, possibly empty.
func constructLegalMoves11(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSwap(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 11)
	return moves
}

// Return the set of legal moves for a pawn using Card12, possibly empty.
func constructLegalMoves12(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 12)
	return moves
}

// Return the set of legal moves for a pawn using CardApologies, possibly empty.
func constructLegalMovesApologies(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []Move {
	moves := make([]Move, 0)
	moveApologies(&moves, color, card, pawn, allPawns)
	return moves
}

// Return the first pawn at the indicated position, or None.
func findPawn(allPawns []model.Pawn, position model.Position) model.Pawn {
	for _, p := range allPawns {
		if p.Position().Equals(position) {
			return p
		}
	}

	return nil
}

func moveCircle(moves *[]Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For start-related cards, a pawn in the start area can move to the associated
	// circle position if that position is not occupied by another pawn of the same color.
	if pawn.Position().Start() {
		conflict := findPawn(allPawns, model.StartCircles[color])
		if conflict == nil {
			actions := []Action { NewAction(MoveToPosition, pawn, model.StartCircles[color].Copy()) }
			sideEffects := make([]Action, 0)
			move := NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		} else if conflict != nil && conflict.Color() != color {
			actions := []Action { NewAction(MoveToPosition, pawn, model.StartCircles[color].Copy())}
			sideEffects := []Action { NewAction(MoveToStart, conflict, nil) }
			move := NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		}
	}
}

func moveSimple(moves *[]Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn, squares int) {
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

func moveSplit(moves *[]Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For the 7 card, we can split up the move between two different pawns.
	// Any combination of 7 forward moves is legal, as long as the resulting position
	// is not occupied by another pawn of the same color.

	for _, other := range allPawns {
		if !other.Equals(pawn) && other.Color() == color && !other.Position().Home() && !other.Position().Start() {

			// any pawn except other
			filtered := make([]model.Pawn, 0)
			for _, p := range allPawns {
				if !p.Equals(other) {
					filtered = append(filtered, p)
				}
			}

			for _, legal := range legalSplits {
				left := make([]Move, 0)
				moveSimple(&left, color, card, pawn, filtered, legal.left)

				right := make([]Move, 0)
				moveSimple(&right, color, card, other, filtered, legal.right)

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

func moveSwap(moves *[]Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
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

func moveApologies(moves *[]Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
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

// Augment any legal moves with additional side-effects that occur as a result of model.Slides.
func augmentWithSlides(allPawns []model.Pawn, moves []Move) {
	for _, move := range moves {
		for _, action := range move.Actions() {
			if action.Type() == MoveToPosition { // look at any move to a position on the board
				for _, color := range model.PlayerColors.Members() {
					if color != action.Pawn().Color() { // any color other than the pawn's
						for _, slide := range model.Slides[color] { // # look at all model.Slides with this color
							if action.Position() != nil && action.Position().Square() != nil && *action.Position().Square() == slide.Start() {
								_ = action.Position().MoveToSquare(slide.End()) // if the pawn landed on the start of the slide, move the pawn to the end of the slide
								for square := slide.Start()+1; square <= slide.End(); square++ {
									// Note: in this one case, a pawn can bump another pawn of the same color
									tmp := model.NewPosition(false, false, nil, &square)
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