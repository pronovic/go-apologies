package source

import (
	"cmp"
	"slices"

	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/reward"
	"github.com/pronovic/go-apologies/rules"
)

type rewardInputSource struct {
	evaluator  rules.Rules
	calculator reward.Calculator
}

type result struct {
	move  model.Move
	score float32
}

// RewardInputSource source of input for a character which chooses its next move based on a reward calculation.
func RewardInputSource(evaluator rules.Rules, calculator reward.Calculator) CharacterInputSource {
	if evaluator == nil {
		evaluator = rules.NewRules(nil)
	}

	if calculator == nil {
		calculator = reward.NewCalculator()
	}

	return &rewardInputSource{
		evaluator:  evaluator,
		calculator: calculator,
	}
}

func (s *rewardInputSource) Name() string {
	return "RewardInputSource"
}

func (s *rewardInputSource) ChooseMove(_ model.GameMode, view model.PlayerView, legalMoves []model.Move) (model.Move, error) {
	results := make([]result, 0, len(legalMoves))

	for _, move := range legalMoves {
		evaluated, err := s.evaluator.EvaluateMove(view, move)
		if err != nil {
			return nil, err
		}

		score := s.calculator.Calculate(evaluated)
		results = append(results, result{move, score})
	}

	// sort the highest-scoring move to the top
	slices.SortStableFunc(results, func(i, j result) int {
		return cmp.Compare(j.score, i.score) // j before i reverses the sort, so largest is at [0]
	})

	// return the highest-scoring move
	return results[0].move, nil
}
