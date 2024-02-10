package source

import (
	"testing"

	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/reward"
	"github.com/pronovic/go-apologies/rules"
	"github.com/stretchr/testify/assert"
)

func TestRewardInputSourceName(t *testing.T) {
	obj := RewardInputSource(nil, nil)
	assert.Equal(t, "RewardInputSource", obj.Name())
}

func TestRewardInputSourceChooseMove(t *testing.T) {
	view := model.MockPlayerView{}

	evaluated1 := model.MockPlayerView{}
	evaluated2 := model.MockPlayerView{}
	evaluated3 := model.MockPlayerView{}

	move1 := model.MockMove{}
	move2 := model.MockMove{}
	move3 := model.MockMove{}
	moves := []model.Move{&move1, &move2, &move3}

	evaluator := rules.MockRules{}
	calculator := reward.MockCalculator{}
	obj := RewardInputSource(&evaluator, &calculator)

	evaluator.On("EvaluateMove", &view, &move1).Return(&evaluated1, nil).Once()
	evaluator.On("EvaluateMove", &view, &move2).Return(&evaluated2, nil).Once()
	evaluator.On("EvaluateMove", &view, &move3).Return(&evaluated3, nil).Once()

	calculator.On("Calculate", &evaluated1).Return(float32(200.0)).Once()
	calculator.On("Calculate", &evaluated2).Return(float32(300.0)).Once()
	calculator.On("Calculate", &evaluated3).Return(float32(100.0)).Once()

	result, err := obj.ChooseMove(model.AdultMode, &view, moves)
	assert.Nil(t, err)
	assert.Same(t, &move2, result)
}
