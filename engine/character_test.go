package engine

import (
	"testing"

	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/source"
	"github.com/stretchr/testify/assert"
)

func TestNewCharacter(t *testing.T) {
	input := source.MockCharacterInputSource{}
	obj := NewCharacter("character", &input)
	assert.Equal(t, "character", obj.Name())
	assert.Same(t, &input, obj.Source())
}

func TestCharacterColor(t *testing.T) {
	input := source.MockCharacterInputSource{}
	obj := NewCharacter("character", &input)
	for _, color := range model.PlayerColors.Members() {
		obj.SetColor(color)
		assert.Equal(t, color, obj.Color())
	}
}

func TestCharacterChooseMove(t *testing.T) {
	input := source.MockCharacterInputSource{}
	obj := NewCharacter("character", &input)
	mode := model.StandardMode
	view := model.MockPlayerView{}
	move := model.MockMove{}
	legalMoves := make([]model.Move, 0)
	input.On("ChooseMove", mode, &view, legalMoves).Return(&move, nil)
	result, err := obj.ChooseMove(mode, &view, legalMoves)
	assert.NoError(t, err)
	assert.Same(t, &move, result)
}
