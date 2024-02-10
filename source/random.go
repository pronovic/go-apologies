package source

import (
	"github.com/pronovic/go-apologies/internal/randomutil"
	"github.com/pronovic/go-apologies/model"
)

type randomInputSource struct{}

// RandomInputSource source of input for a character which chooses randomly from among legal moves.
func RandomInputSource() CharacterInputSource {
	return &randomInputSource{}
}

func (s *randomInputSource) Name() string {
	return "RandomInputSource"
}

func (s *randomInputSource) ChooseMove(_ model.GameMode, _ model.PlayerView, legalMoves []model.Move) (model.Move, error) {
	return randomutil.RandomChoice(legalMoves)
}
