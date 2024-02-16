package source

import (
	"testing"

	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
)

func TestRandomInputSourceName(t *testing.T) {
	obj := RandomInputSource()
	assert.Equal(t, "RandomInputSource", obj.Name())
}

func TestRandomInputSourceChooseMove(t *testing.T) {
	obj := RandomInputSource()

	move1 := model.MockMove{}
	move2 := model.MockMove{}
	move3 := model.MockMove{}
	moves := []model.Move{&move1, &move2, &move3}

	for i := 0; i < 100; i++ {
		result, err := obj.ChooseMove(model.AdultMode, nil, moves)
		assert.NoError(t, err)
		assert.True(t, result == &move1 || result == &move2 || result == &move3)
	}
}
