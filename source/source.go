package source

import (
	"github.com/pronovic/go-apologies/model"
)

type CharacterInputSource interface {
	// Name Get the name of the character input source.
	Name() string

	// ChooseMove Choose the next move for a character.
	//
	// There is always at least one legal move: a forfeit.  Nothing else is legal, so the
	// character must choose to discard one card.  In standard mode, there is effectively no
	// choice (since there is only one card in play), but in adult mode the character can choose
	// which to discard.  If a move has an empty list of actions, then this is a forfeit.
	//
	// The source must return a move from among the passed-in set of legal moves.  If a source
	// returns an illegal move, then a legal move will be chosen at random and executed.  This way,
	// a misbehaving source (or a source attempting to cheat) does not get an advantage.  The game
	// rules require a player to make a legal move if one is available, even if that move is
	// disadvantageous.
	ChooseMove(mode model.GameMode, view model.PlayerView, legalMoves []model.Move) (model.Move, error)
}
