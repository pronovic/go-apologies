package pkg

import (
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/model"
)

// splitPair defines a legal way to split up a move of 7
type splitPair struct {
	left int
	right int
}

// legalSplits defines legal ways to split up a move of 7
var legalSplits = []splitPair{
	{1, 6},
	{2, 5},
	{3, 4},
	{4, 3},
	{5, 2},
	{6, 1},
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
func ConstructLegalMoves(view model.PlayerView, card model.Card) ([]model.Move, error) {
	moves := make([]model.Move, 0)
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
			moves = append(moves, model.NewMove(played, []model.Action{}, []model.Action{}))
		}
	}

	if len(moves) == 0 {
		return []model.Move{}, errors.New("internal error: could not construct any legal moves")
	}

	return moves, nil
}

// ExecuteMove Execute a player's move, updating game state
func ExecuteMove(game model.Game, player model.Player, move model.Move) error {
	for _, action := range move.MergedActions() { // execute actions, then side effects, in order
		// keep in mind that the pawn on the action is a different object than the pawn in the game
		pawn := game.Players()[action.Pawn().Color()].Pawns()[action.Pawn().Index()]
		if action.Type() == model.MoveToStart {
			game.Track(fmt.Sprintf("Played card %s: [%s->start]", move.Card().Type().Value(), pawn.Name()), player, move.Card())
			err := pawn.Position().MoveToStart()
			if err != nil {
				return err
			}
		} else if action.Type() == model.MoveToPosition && action.Position() != nil {
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
func EvaluateMove(view model.PlayerView, move model.Move) (model.PlayerView, error) {
	result := view.Copy()

	for _, action := range move.MergedActions() { // execute actions, then side effects, in order
		// keep in mind that the pawn on the action is a different object than the pawn in the game
		pawn := result.GetPawn(action.Pawn())
		if pawn != nil {  // if the pawn isn't valid, just ignore it
			if action.Type() == model.MoveToStart {
				err := pawn.Position().MoveToStart()
				if err != nil {
					return nil, err
				}
			} else if action.Type() == model.MoveToPosition && action.Position() != nil {
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
func constructLegalMoves(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	var moves []model.Move
	if pawn.Position().Home() {
		moves = make([]model.Move, 0)
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

// DistanceToHome Return the distance to home for this pawn, a number of squares when moving forward.
func DistanceToHome(pawn model.Pawn) int {
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
func constructLegalMoves1(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveCircle(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 1)
	return moves
}

// Return the set of legal moves for a pawn using Card2, possibly empty.
func constructLegalMoves2(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveCircle(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 2)
	return moves
}

// Return the set of legal moves for a pawn using Card3, possibly empty.
func constructLegalMoves3(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 3)
	return moves
}

// Return the set of legal moves for a pawn using Card4, possibly empty.
func constructLegalMoves4(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, -4)
	return moves
}

// Return the set of legal moves for a pawn using Card5, possibly empty.
func constructLegalMoves5(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 5)
	return moves
}

// Return the set of legal moves for a pawn using Card7, possibly empty.
func constructLegalMoves7(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 7)
	moveSplit(&moves, color, card, pawn, allPawns)
	return moves
}

// Return the set of legal moves for a pawn using Card8, possibly empty.
func constructLegalMoves8(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 8)
	return moves
}

// Return the set of legal moves for a pawn using Card10, possibly empty.
func constructLegalMoves10(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 10)
	moveSimple(&moves, color, card, pawn, allPawns, -1)
	return moves
}

// Return the set of legal moves for a pawn using Card11, possibly empty.
func constructLegalMoves11(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSwap(&moves, color, card, pawn, allPawns)
	moveSimple(&moves, color, card, pawn, allPawns, 11)
	return moves
}

// Return the set of legal moves for a pawn using Card12, possibly empty.
func constructLegalMoves12(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	moveSimple(&moves, color, card, pawn, allPawns, 12)
	return moves
}

// Return the set of legal moves for a pawn using CardApologies, possibly empty.
func constructLegalMovesApologies(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
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

func moveCircle(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For start-related cards, a pawn in the start area can move to the associated
	// circle position if that position is not occupied by another pawn of the same color.
	if pawn.Position().Start() {
		conflict := findPawn(allPawns, model.StartCircles[color])
		if conflict == nil {
			actions := []model.Action { model.NewAction(model.MoveToPosition, pawn, model.StartCircles[color].Copy()) }
			sideEffects := make([]model.Action, 0)
			move := model.NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		} else if conflict != nil && conflict.Color() != color {
			actions := []model.Action { model.NewAction(model.MoveToPosition, pawn, model.StartCircles[color].Copy())}
			sideEffects := []model.Action { model.NewAction(model.MoveToStart, conflict, nil) }
			move := model.NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		}
	}
}

func moveSimple(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn, squares int) {
	// For most cards, a pawn on the board can move forward or backward if the
	// resulting position is not occupied by another pawn of the same color.
	if pawn.Position().Square() != nil || pawn.Position().Safe() != nil {
		target, err := calculatePosition(color, pawn.Position(), squares)
		if err == nil { // if the requested position is not legal, then just ignore it
			if target.Home() || target.Start() { // by definition, there can't be a conflict going to home or start
				actions := []model.Action { model.NewAction(model.MoveToPosition, pawn, target) }
				sideEffects := make([]model.Action, 0)
				move := model.NewMove(card, actions, sideEffects)
				*moves = append(*moves, move)
			} else {
				conflict := findPawn(allPawns, target)
				if conflict == nil {
					actions := []model.Action { model.NewAction(model.MoveToPosition, pawn, target) }
					sideEffects := make([]model.Action, 0)
					move := model.NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				} else if conflict != nil && conflict.Color() != color {
					actions := []model.Action { model.NewAction(model.MoveToPosition, pawn, target)}
					sideEffects := []model.Action { model.NewAction(model.MoveToStart, conflict, nil) }
					move := model.NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				}
			}
		}
	}
}

func moveSplit(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
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
				left := make([]model.Move, 0)
				moveSimple(&left, color, card, pawn, filtered, legal.left)

				right := make([]model.Move, 0)
				moveSimple(&right, color, card, other, filtered, legal.right)

				if len(left) > 0 && len(right) > 0 {
					actions := make([]model.Action, 0)
					sideEffects := make([]model.Action, 0)

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

					move := model.NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				}
			}
		}
	}
}

func moveSwap(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For the 11 card, a pawn on the board can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	if pawn.Position().Square() != nil { // pawn is on the board
		 for _, swap := range allPawns {
			 if swap.Color() != color && !swap.Position().Home() && !swap.Position().Start() && swap.Position().Safe() == nil {
				 actions := []model.Action {
					 model.NewAction(model.MoveToPosition, pawn, swap.Position().Copy()),
					 model.NewAction(model.MoveToPosition, swap, pawn.Position().Copy()),
				 }
				 sideEffects := make([]model.Action, 0)
				 move := model.NewMove(card, actions, sideEffects)
				 *moves = append(*moves, move)
			 }
		 }
	}
}

func moveApologies(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For the Apologies card, a pawn in start can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	if pawn.Position().Start() {
		for _, swap := range allPawns {
			if swap.Color() != color && !swap.Position().Home() && !swap.Position().Start() && swap.Position().Safe() == nil {
				actions := []model.Action {
					model.NewAction(model.MoveToPosition, pawn, swap.Position().Copy()),
					model.NewAction(model.MoveToStart, swap, nil),
				}
				sideEffects := make([]model.Action, 0)
				move := model.NewMove(card, actions, sideEffects)
				*moves = append(*moves, move)
			}
		}
	}
}

// Augment any legal moves with additional side-effects that occur as a result of model.Slides.
func augmentWithSlides(allPawns []model.Pawn, moves []model.Move) {
	for _, move := range moves {
		for _, action := range move.Actions() {
			if action.Type() == model.MoveToPosition { // look at any move to a position on the board
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
										bump := model.NewAction(model.MoveToStart, pawn, nil)
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