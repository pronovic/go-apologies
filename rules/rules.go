package rules

import (
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/generator"
	"github.com/pronovic/go-apologies/internal/equality"
	"github.com/pronovic/go-apologies/model"
)

// Rules provides high-level game rules
type Rules interface {

	// StartGame starts a game using the passed-in mode
	StartGame(game model.Game, mode model.GameMode) error

	// ExecuteMove Execute a player's move, updating game state
	ExecuteMove(game model.Game, player model.Player, move model.Move) error

	// EvaluateMove constructs a new player view that results from executing the passed-in move.
	// This is equivalent to execute_move() but has no permanent effect on the game.  It's intended for
	// use by a character, to evaluate the results of each legal move.
	EvaluateMove(view model.PlayerView, move model.Move) (model.PlayerView, error)

	// ConstructLegalMoves returns the set of all legal moves for a player and its opponents
	// Pass the card to play, or nil if the move should come from the player's hand
	ConstructLegalMoves(view model.PlayerView, card model.Card) ([]model.Move, error)

	// DrawAgain Whether the player gets to draw again based on the passed-in card
	DrawAgain(card model.Card) bool

}

type rules struct {
	moveGenerator generator.MoveGenerator
}

// NewRules creates a new rules interface, optionally accepting a move generator
func NewRules(moveGenerator generator.MoveGenerator) Rules {
	if moveGenerator == nil {
		moveGenerator = generator.NewGenerator()
	}

	return &rules {
		moveGenerator: moveGenerator,
	}
}

func (r *rules) StartGame(game model.Game, mode model.GameMode) error {
	if game == nil {
		return errors.New("game is nil")
	}

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

func (r *rules) ExecuteMove(game model.Game, player model.Player, move model.Move) error {
	if game == nil {
		return errors.New("game is nil")
	}

	if player == nil {
		return errors.New("player is nil")
	}

	if move == nil {
		return errors.New("move is nil")
	}

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

func (r *rules) EvaluateMove(view model.PlayerView, move model.Move) (model.PlayerView, error) {
	if view == nil {
		return nil, errors.New("view is nil")
	}

	if move == nil {
		return nil, errors.New("move is nil")
	}

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

func (r *rules) ConstructLegalMoves(view model.PlayerView, card model.Card) ([]model.Move, error) {
	if view == nil {
		return nil, errors.New("view is nil")
	}

	allPawns := view.AllPawns()  // pre-calculate this once up-front

	var cards []model.Card
	if card != nil {
		cards = []model.Card { card }
	} else {
		cards = view.Player().Hand()
	}

	moves := make([]model.Move, 0)
	for _, played := range cards {
		for _, pawn := range view.Player().Pawns() {
			for _, move := range r.moveGenerator.LegalMoves(view.Player().Color(), played, pawn, allPawns) {
				if ! contains(moves, move) {
					moves = append(moves, move)  // eliminate duplicates
				}
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

func (r *rules) DrawAgain(card model.Card) bool {
	if card == nil {
		return false
	} else {
		return model.DrawAgain[card.Type()]
	}
}

func contains(moves []model.Move, move model.Move) bool {
	for _, element := range moves {
		if equality.EqualByValue(element, move) {
			return true
		}
	}

	return false
}

