package engine

import (
	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/rules"
	"github.com/pronovic/go-apologies/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEngine(t *testing.T) {
	input := source.MockCharacterInputSource{}

	character1 := NewCharacter("character1", &input)
	character2 := NewCharacter("character2", &input)
	character3 := NewCharacter("character3", &input)
	character4 := NewCharacter("character4", &input)
	characters := []Character { character1, character2, character3, character4 }

	e, _ := NewEngine(model.StandardMode, characters, nil)
	assert.Equal(t, model.StandardMode, e.Mode())
	assert.Equal(t, characters, e.Characters())
	assert.NotNil(t, e.First())
	assert.Equal(t, 4, e.Players())
	assert.Equal(t, "Game waiting to start", e.State())
	assert.NotNil(t, e.Game())
	assert.Equal(t, 4, len(e.Game().Players()))
	assert.Equal(t, false, e.Started())
	assert.Equal(t, false, e.Completed())
	assert.Equal(t, character1, e.ColorMap()[model.Red])
	assert.Equal(t, character2, e.ColorMap()[model.Yellow])
	assert.Equal(t, character3, e.ColorMap()[model.Green])
	assert.Equal(t, character4, e.ColorMap()[model.Blue])
}

func TestEngineFirst(t *testing.T) {
	e := createEngine(model.AdultMode, nil)

	first := e.First()
	assert.NotNil(t, first)

	for _, color := range model.PlayerColors.Members() {
		e.SetFirst(color)
		assert.Equal(t, color, e.First())
	}
}

func TestEngineStarted(t *testing.T) {
	e := createEngine(model.AdultMode, nil)
	assert.False(t, e.Started())
	assert.Equal(t, "Game waiting to start", e.State())
	_, _ = e.StartGame()
	assert.True(t, e.Started())
	assert.Equal(t, "Game in progress", e.State())
}

func TestEngineCompleted(t *testing.T) {
	e := createEngine(model.AdultMode, nil)
	assert.False(t, e.Completed())

	// move all of red player's pawns to home, which makes them the winner
	for _, pawn := range e.Game().Players()[model.Red].Pawns() {
		_ = pawn.Position().MoveToHome()
	}

	assert.True(t, e.Completed())
}

func TestEngineWinner(t *testing.T) {
	e := createEngine(model.AdultMode, nil)
	assert.Nil(t, e.Winner())

	// move all of red player's pawns to home, which makes them the winner
	for _, pawn := range e.Game().Players()[model.Red].Pawns() {
		_ = pawn.Position().MoveToHome()
	}

	assert.Same(t, e.ColorMap()[model.Red], e.Winner())
}

func TestEngineReset(t *testing.T) {
	e := createEngine(model.AdultMode, nil)
	saved := e.Game()
	game, err := e.Reset()
	assert.Nil(t, err)
	assert.Same(t, e.Game(), game)
	assert.NotSame(t, saved, game)
	assert.False(t, e.Game().Started())
}

func TestEngineStartGame(t *testing.T) {
	e := createEngine(model.AdultMode, nil)
	assert.False(t, e.Started())
	assert.False(t, e.Game().Started())
	game, err := e.StartGame()
	assert.Nil(t, err)
	assert.Same(t, e.Game(), game)
	assert.True(t, e.Started())
	assert.True(t, e.Game().Started())
}

func TestEngineDrawAndDiscard(t *testing.T) {
	e := createEngine(model.AdultMode, nil)

	// draw all of the cards from the deck
	var drawn = make([]model.Card, 0, model.DeckSize)
	for i := 0; i < model.DeckSize; i++ {
		c, err := e.Draw()
		assert.Nil(t, err)
		drawn = append(drawn, c)
	}

	// confirm that we get an error, because the deck is empty
	_, err := e.Draw()
	assert.NotNil(t, err)

	// put back one card
	last := drawn[model.DeckSize - 1]
	err = e.Discard(last)
	assert.Nil(t, err)

	// now draw it again and confirm we get it back
	c, err := e.Draw()
	assert.Nil(t, err)
	assert.Same(t, last, c)
}

func TestEngineConstructLegalMovesStandardNoCard(t *testing.T) {
	evaluator := rules.MockRules{}
	e := createEngine(model.StandardMode, &evaluator)

	view := model.MockPlayerView{}
	drawcard := model.NewCard("1", model.Card1)
	movecard := model.NewCard("2", model.Card2)
	move := model.NewMove(movecard, []model.Action{}, []model.Action{})
	legalMoves := []model.Move{ move }

	configureDrawCard(e, drawcard) // so we know exactly which card will be drawn
	evaluator.On("ConstructLegalMoves", &view, drawcard).Return(legalMoves, nil)

	c, m, err := e.ConstructLegalMoves(&view, nil)
	assert.Nil(t, err)
	assert.Same(t, drawcard, c)
	assert.Equal(t, legalMoves, m)

	_, err = e.Draw()
	assert.NotNil(t, err) // confirm that the deck is empty and we did draw the drawcard
}

func TestEngineConstructLegalMovesStandardCard(t *testing.T) {
	evaluator := rules.MockRules{}
	e := createEngine(model.StandardMode, &evaluator)

	view := model.MockPlayerView{}
	drawcard := model.NewCard("1", model.Card1)
	movecard := model.NewCard("2", model.Card2)
	providedcard := model.NewCard("3", model.Card3)
	move := model.NewMove(movecard, []model.Action{}, []model.Action{})
	legalMoves := []model.Move{ move }

	configureDrawCard(e, drawcard) // so we know exactly which card will be drawn
	evaluator.On("ConstructLegalMoves", &view, providedcard).Return(legalMoves, nil)

	c, m, err := e.ConstructLegalMoves(&view, providedcard)
	assert.Nil(t, err)
	assert.Same(t, providedcard, c)
	assert.Equal(t, legalMoves, m)

	c, _ = e.Draw()
	assert.Same(t, drawcard, c)  // confirm that the card was not drawn from the deck
}

func TestEngineConstructLegalMovesAdultNoCard(t *testing.T) {
	evaluator := rules.MockRules{}
	e := createEngine(model.AdultMode, &evaluator)

	view := model.MockPlayerView{}
	drawcard := model.NewCard("1", model.Card1)
	movecard := model.NewCard("2", model.Card2)
	move := model.NewMove(movecard, []model.Action{}, []model.Action{})
	legalMoves := []model.Move{ move }

	configureDrawCard(e, drawcard) // so we know exactly which card will be drawn
	evaluator.On("ConstructLegalMoves", &view, nil).Return(legalMoves, nil)

	c, m, err := e.ConstructLegalMoves(&view, nil)
	assert.Nil(t, err)
	assert.Nil(t, c)
	assert.Equal(t, legalMoves, m)

	c, _ = e.Draw()
	assert.Same(t, drawcard, c)  // confirm that the card was not drawn from the deck
}

func TestEngineConstructLegalMovesAdultCard(t *testing.T) {
	evaluator := rules.MockRules{}
	e := createEngine(model.AdultMode, &evaluator)

	view := model.MockPlayerView{}
	drawcard := model.NewCard("1", model.Card1)
	movecard := model.NewCard("2", model.Card2)
	providedcard := model.NewCard("3", model.Card3)
	move := model.NewMove(movecard, []model.Action{}, []model.Action{})
	legalMoves := []model.Move{ move }

	configureDrawCard(e, drawcard) // so we know exactly which card will be drawn
	evaluator.On("ConstructLegalMoves", &view, providedcard).Return(legalMoves, nil)

	c, m, err := e.ConstructLegalMoves(&view, providedcard)
	assert.Nil(t, err)
	assert.Same(t, providedcard, c)
	assert.Equal(t, legalMoves, m)

	c, _ = e.Draw()
	assert.Same(t, drawcard, c)  // confirm that the card was not drawn from the deck
}

func TestEnginePlayNextCompleted(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextFailed(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextStandardForfeit(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextStandardIllegal(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextStandardLegal(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextStandardDrawAgain(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextStandardComplete(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextAdultForfeit(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextAdultIllegal(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextAdultLegal(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextAdultDrawAgain(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

func TestEnginePlayNextAdultComplete(t *testing.T) {
	t.Fail()  // TODO: implement test case
}

// createEngine creates an engine for testing, to avoid boilerplate in other methods
// a nil evaluator gets you a real rule.Rules implementation, otherwise pass in a rules.MockRules
func createEngine(mode model.GameMode, evaluator rules.Rules) Engine {
	input := source.MockCharacterInputSource{}

	character1 := NewCharacter("character1", &input)
	character2 := NewCharacter("character2", &input)
	characters := []Character { character1, character2 }

	e, _ := NewEngine(mode, characters, evaluator)
	e.SetFirst(model.Red)

	return e
}

// configureDrawCard configures the deck with a single card in it to be drawn
func configureDrawCard(e Engine, drawcard model.Card) {
	for i := 0; i < model.DeckSize; i++ {
		_, _ = e.Draw()
	}

	_ = e.Discard(drawcard)
}