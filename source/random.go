package source

import (
	"crypto/rand"
	"errors"
	"github.com/pronovic/go-apologies/model"
	"math/big"
)

type randomInputSource struct { }

// RandomInputSource source of input for a character which chooses randomly from among legal moves.
func RandomInputSource() CharacterInputSource {
	return &randomInputSource {	}
}

func (s *randomInputSource) Name() string {
	return "RandomInputSource"
}

func (s *randomInputSource) ChooseMove(_ model.GameMode, _ model.PlayerView, legalMoves []model.Move) (model.Move, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(legalMoves))))
	if err != nil {
		return nil, errors.New("failed to generate random int for move")
	}

	return legalMoves[int(index.Int64())], nil
}