package generator

import (
	"errors"

	"github.com/pronovic/go-apologies/internal/equality"
	"github.com/pronovic/go-apologies/model"
)

// splitPair defines a legal way to split up a move of 7
type splitPair struct {
	left  int
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

type MoveGenerator interface {
	LegalMoves(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move
	CalculatePosition(color model.PlayerColor, position model.Position, squares int) (model.Position, error)
}

type moveGenerator struct{}

// NewGenerator constructs a new move generator, optionally accepting an identifier factory
func NewGenerator() MoveGenerator {
	return &moveGenerator{}
}

// LegalMoves Generate the set of legal moves for a pawn using a card, possibly empty.
func (g *moveGenerator) LegalMoves(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	var moves []model.Move
	if pawn.Position().Home() {
		moves = make([]model.Move, 0)
	} else {
		switch card.Type() {
		case model.Card1:
			moves = g.legalMovesCard1(color, card, pawn, allPawns)
		case model.Card2:
			moves = g.legalMovesCard2(color, card, pawn, allPawns)
		case model.Card3:
			moves = g.legalMovesCard3(color, card, pawn, allPawns)
		case model.Card4:
			moves = g.legalMovesCard4(color, card, pawn, allPawns)
		case model.Card5:
			moves = g.legalMovesCard5(color, card, pawn, allPawns)
		case model.Card7:
			moves = g.legalMovesCard7(color, card, pawn, allPawns)
		case model.Card8:
			moves = g.legalMovesCard8(color, card, pawn, allPawns)
		case model.Card10:
			moves = g.legalMovesCard10(color, card, pawn, allPawns)
		case model.Card11:
			moves = g.legalMovesCard11(color, card, pawn, allPawns)
		case model.Card12:
			moves = g.legalMovesCard12(color, card, pawn, allPawns)
		case model.CardApologies:
			moves = g.legalMovesApologies(color, card, pawn, allPawns)
		}
	}
	g.augmentWithSlides(allPawns, moves)
	return moves
}

// Return the set of legal moves for a pawn using Card1, possibly empty.
func (g *moveGenerator) legalMovesCard1(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveCircle(&moves, color, card, pawn, allPawns)
	g.moveSimple(&moves, color, card, pawn, allPawns, 1)
	return moves
}

// Return the set of legal moves for a pawn using Card2, possibly empty.
func (g *moveGenerator) legalMovesCard2(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveCircle(&moves, color, card, pawn, allPawns)
	g.moveSimple(&moves, color, card, pawn, allPawns, 2)
	return moves
}

// Return the set of legal moves for a pawn using Card3, possibly empty.
func (g *moveGenerator) legalMovesCard3(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSimple(&moves, color, card, pawn, allPawns, 3)
	return moves
}

// Return the set of legal moves for a pawn using Card4, possibly empty.
func (g *moveGenerator) legalMovesCard4(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSimple(&moves, color, card, pawn, allPawns, -4)
	return moves
}

// Return the set of legal moves for a pawn using Card5, possibly empty.
func (g *moveGenerator) legalMovesCard5(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSimple(&moves, color, card, pawn, allPawns, 5)
	return moves
}

// Return the set of legal moves for a pawn using Card7, possibly empty.
func (g *moveGenerator) legalMovesCard7(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSimple(&moves, color, card, pawn, allPawns, 7)
	g.moveSplit(&moves, color, card, pawn, allPawns)
	return moves
}

// Return the set of legal moves for a pawn using Card8, possibly empty.
func (g *moveGenerator) legalMovesCard8(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSimple(&moves, color, card, pawn, allPawns, 8)
	return moves
}

// Return the set of legal moves for a pawn using Card10, possibly empty.
func (g *moveGenerator) legalMovesCard10(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSimple(&moves, color, card, pawn, allPawns, 10)
	g.moveSimple(&moves, color, card, pawn, allPawns, -1)
	return moves
}

// Return the set of legal moves for a pawn using Card11, possibly empty.
func (g *moveGenerator) legalMovesCard11(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSwap(&moves, color, card, pawn, allPawns)
	g.moveSimple(&moves, color, card, pawn, allPawns, 11)
	return moves
}

// Return the set of legal moves for a pawn using Card12, possibly empty.
func (g *moveGenerator) legalMovesCard12(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveSimple(&moves, color, card, pawn, allPawns, 12)
	return moves
}

// Return the set of legal moves for a pawn using CardApologies, possibly empty.
func (g *moveGenerator) legalMovesApologies(color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) []model.Move {
	moves := make([]model.Move, 0)
	g.moveApologies(&moves, color, card, pawn, allPawns)
	return moves
}

// Return the first pawn at the indicated position, or None.
func (g *moveGenerator) findPawn(allPawns []model.Pawn, position model.Position) model.Pawn {
	for _, p := range allPawns {
		if equality.EqualByValue(p.Position(), position) {
			return p
		}
	}

	return nil
}

func (g *moveGenerator) moveCircle(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For start-related cards, a pawn in the start area can move to the associated
	// circle position if that position is not occupied by another pawn of the same color.
	if pawn.Position().Start() {
		conflict := g.findPawn(allPawns, model.StartCircles[color])
		if conflict == nil {
			actions := []model.Action{model.NewAction(model.MoveToPosition, pawn, model.StartCircles[color].Copy())}
			sideEffects := make([]model.Action, 0)
			move := model.NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		} else if conflict != nil && conflict.Color() != color {
			actions := []model.Action{model.NewAction(model.MoveToPosition, pawn, model.StartCircles[color].Copy())}
			sideEffects := []model.Action{model.NewAction(model.MoveToStart, conflict, nil)}
			move := model.NewMove(card, actions, sideEffects)
			*moves = append(*moves, move)
		}
	}
}

func (g *moveGenerator) moveSimple(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn, squares int) {
	// For most cards, a pawn on the board can move forward or backward if the
	// resulting position is not occupied by another pawn of the same color.
	if pawn.Position().Square() != nil || pawn.Position().Safe() != nil {
		target, err := g.CalculatePosition(color, pawn.Position(), squares)
		if err == nil { // if the requested position is not legal, then just ignore it
			if target.Home() || target.Start() { // by definition, there can't be a conflict going to home or start
				actions := []model.Action{model.NewAction(model.MoveToPosition, pawn, target)}
				sideEffects := make([]model.Action, 0)
				move := model.NewMove(card, actions, sideEffects)
				*moves = append(*moves, move)
			} else {
				conflict := g.findPawn(allPawns, target)
				if conflict == nil {
					actions := []model.Action{model.NewAction(model.MoveToPosition, pawn, target)}
					sideEffects := make([]model.Action, 0)
					move := model.NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				} else if conflict != nil && conflict.Color() != color {
					actions := []model.Action{model.NewAction(model.MoveToPosition, pawn, target)}
					sideEffects := []model.Action{model.NewAction(model.MoveToStart, conflict, nil)}
					move := model.NewMove(card, actions, sideEffects)
					*moves = append(*moves, move)
				}
			}
		}
	}
}

func (g *moveGenerator) moveSplit(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For the 7 card, we can split up the move between two different pawns.
	// Any combination of 7 forward moves is legal, as long as the resulting position
	// is not occupied by another pawn of the same color.

	for _, other := range allPawns {
		if !equality.EqualByValue(other, pawn) && other.Color() == color && !other.Position().Home() && !other.Position().Start() {

			// any pawn except other
			filtered := make([]model.Pawn, 0)
			for _, p := range allPawns {
				if !equality.EqualByValue(p, other) {
					filtered = append(filtered, p)
				}
			}

			for _, legal := range legalSplits {
				left := make([]model.Move, 0)
				g.moveSimple(&left, color, card, pawn, filtered, legal.left)

				right := make([]model.Move, 0)
				g.moveSimple(&right, color, card, other, filtered, legal.right)

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

func (g *moveGenerator) moveSwap(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For the 11 card, a pawn on the board can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	if pawn.Position().Square() != nil { // pawn is on the board
		for _, swap := range allPawns {
			if swap.Color() != color && !swap.Position().Home() && !swap.Position().Start() && swap.Position().Safe() == nil {
				actions := []model.Action{
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

func (g *moveGenerator) moveApologies(moves *[]model.Move, color model.PlayerColor, card model.Card, pawn model.Pawn, allPawns []model.Pawn) {
	// For the Apologies card, a pawn in start can swap with another pawn of a different
	// color, as long as that pawn is outside of the start area, safe area, or home area.
	if pawn.Position().Start() {
		for _, swap := range allPawns {
			if swap.Color() != color && !swap.Position().Home() && !swap.Position().Start() && swap.Position().Safe() == nil {
				actions := []model.Action{
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
func (g *moveGenerator) augmentWithSlides(allPawns []model.Pawn, moves []model.Move) {
	for _, move := range moves {
		for _, action := range move.Actions() {
			if action.Type() == model.MoveToPosition { // look at any move to a position on the board
				for _, color := range model.PlayerColors.Members() {
					if color != action.Pawn().Color() { // any color other than the pawn's
						for _, slide := range model.Slides[color] { // # look at all model.Slides with this color
							if action.Position() != nil && action.Position().Square() != nil && *action.Position().Square() == slide.Start() {
								_ = action.Position().MoveToSquare(slide.End()) // if the pawn landed on the start of the slide, move the pawn to the end of the slide
								for square := slide.Start() + 1; square <= slide.End(); square++ {
									// Note: in this one case, a pawn can bump another pawn of the same color
									tmp := model.NewPosition(false, false, nil, &square)
									pawn := g.findPawn(allPawns, tmp)
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

// CalculatePosition Calculate the new position for a forward or backwards move, taking into account safe zone turns but disregarding Slides.
func (g *moveGenerator) CalculatePosition(color model.PlayerColor, position model.Position, squares int) (model.Position, error) {
	if position.Home() || position.Start() {
		return (model.Position)(nil), errors.New("pawn in home or start may not move")
	} else if position.Safe() != nil {
		if squares == 0 {
			return position.Copy(), nil
		} else if squares > 0 {
			if *position.Safe()+squares < model.SafeSquares {
				copied := position.Copy()
				err := copied.MoveToSafe(*position.Safe() + squares)
				if err != nil {
					return (model.Position)(nil), err
				}
				return copied, nil
			} else if *position.Safe()+squares == model.SafeSquares {
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
			if *position.Safe()+squares >= 0 {
				copied := position.Copy()
				err := copied.MoveToSafe(*position.Safe() + squares)
				if err != nil {
					return (model.Position)(nil), err
				}
				return copied, nil
			} else { // handle moving back out of the safe area
				copied := position.Copy()
				err := copied.MoveToSquare(*model.TurnSquares[color].Square())
				if err != nil {
					return (model.Position)(nil), err
				}
				return g.CalculatePosition(color, copied, squares+*position.Safe()+1)
			}
		}
	} else if position.Square() != nil {
		if squares == 0 {
			return position.Copy(), nil
		} else if squares > 0 {
			if *position.Square()+squares < model.BoardSquares {
				if *position.Square() <= *model.TurnSquares[color].Square() && *position.Square()+squares > *model.TurnSquares[color].Square() {
					copied := position.Copy()
					err := copied.MoveToSafe(0)
					if err != nil {
						return (model.Position)(nil), err
					}
					return g.CalculatePosition(color, copied, squares-(*model.TurnSquares[color].Square()-*position.Square())-1)
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
				return g.CalculatePosition(color, copied, squares-(model.BoardSquares-*position.Square()))
			}
		} else { // squares < 0
			if *position.Square()+squares >= 0 {
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
				return g.CalculatePosition(color, copied, squares+*position.Square()+1)
			}
		}
	} else {
		return (model.Position)(nil), errors.New("position is in an illegal state")
	}
}
