package pkg

import (
	"fmt"
	"github.com/pronovic/go-apologies/pkg/util/timestamp"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func init() {
	timestamp.UseStubbedTime()  // once this has been called, it takes effect permanently for all unit tests
}

func TestNewSlide(t *testing.T) {
	obj := newSlide(1, 2)
	assert.Equal(t, 1, obj.Start())
	assert.Equal(t, 2, obj.End())
}

func TestNewCard(t *testing.T) {
	obj := NewCard("id", CardApologies)
	assert.Equal(t, "id", obj.Id())
	assert.Equal(t, CardApologies, obj.Type())
}

func TestCardCopy(t *testing.T) {
	obj := NewCard("id", CardApologies)
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestNewDeck(t *testing.T) {
	obj := NewDeck()
	underlying := obj.(*deck)

	assert.Equal(t, DeckSize, len(underlying.drawPile))
	assert.Equal(t, 0, len(underlying.discardPile))

	var counts = make(map[CardType]int, len(CardTypes.Members()))
	for i := range CardTypes.Members() {
		cardtype := CardTypes.Members()[i]
		counts[cardtype] = 0
	}

	for _, value := range underlying.drawPile {
		cardtype := value.Type()
		counts[cardtype] += 1
	}

	for i := range CardTypes.Members() {
		cardtype := CardTypes.Members()[i]
		assert.Equal(t, counts[cardtype], DeckCounts[cardtype])
	}
}

func TestDeckCopy(t *testing.T) {
	obj := NewDeck()
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestDeckDrawAndDiscard(t *testing.T) {
	var card1 Card
	var card2 Card
	var card3 Card
	var err error

	obj := NewDeck()
	underlying := obj.(*deck)

	// Check that we can draw the entire deck
	var drawn = make([]Card, 0, DeckSize)
	for i := 0; i < DeckSize; i++ {
		card1, err = obj.Draw()
		assert.Nil(t, err)
		drawn = append(drawn, card1)
	}
	assert.Equal(t, DeckSize, len(drawn))
	assert.Equal(t, 0, len(underlying.drawPile))
	_, err = obj.Draw()
	assert.EqualError(t, err, "no cards available in deck")

	// Discard one card and prove that we can draw it
	card1 = drawn[0]
	drawn = slices.Delete(drawn, 0, 1)
	assert.Equal(t, len(underlying.discardPile), 0)
	err = obj.Discard(card1)
	assert.Nil(t, err)
	assert.Equal(t, len(underlying.discardPile), 1)
	card2, err = obj.Draw()
	assert.Same(t, card1, card2)
	assert.Equal(t, len(underlying.discardPile), 0)
	assert.Equal(t, len(underlying.drawPile), 0)

	// Confirm that we're not allowed to discard the same card twice
	card1 = drawn[0]
	drawn = slices.Delete(drawn, 0, 1)
	err = obj.Discard(card1)
	assert.Nil(t, err)
	err = obj.Discard(card1)
	assert.EqualError(t, err, "card already exists in deck")

	// Discard a few others and can prove that they can also be drawn
	card1 = drawn[0]
	drawn = slices.Delete(drawn, 0, 1)
	card2 = drawn[0]
	drawn = slices.Delete(drawn, 0, 1)
	card3 = drawn[0]
	drawn = slices.Delete(drawn, 0, 1)
	err = obj.Discard(card1)
	assert.Nil(t, err)
	err = obj.Discard(card2)
	assert.Nil(t, err)
	err = obj.Discard(card3)
	assert.Nil(t, err)
	assert.Equal(t, len(underlying.discardPile), 4)
	assert.Equal(t, len(underlying.drawPile), 0)
	_, err = obj.Draw()
	assert.Nil(t, err)
	assert.Equal(t, len(underlying.discardPile), 0)
	assert.Equal(t, len(underlying.drawPile), 3)
	_, err = obj.Draw()
	assert.Nil(t, err)
	assert.Equal(t, len(underlying.discardPile), 0)
	assert.Equal(t, len(underlying.drawPile), 2)
	_, err = obj.Draw()
	assert.Nil(t, err)
	assert.Equal(t, len(underlying.discardPile), 0)
	assert.Equal(t, len(underlying.drawPile), 1)
	_, err = obj.Draw()
	assert.Nil(t, err)
	assert.Equal(t, len(underlying.discardPile), 0)
	assert.Equal(t, len(underlying.drawPile), 0)

	// Make sure that the deck still gives an error when empty
	_, err = obj.Draw()
	assert.EqualError(t, err, "no cards available in deck")
}

func TestNewPosition(t *testing.T) {
	var obj Position

	obj = NewPosition(true, false, nil, nil)
	assert.Equal(t, true, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "start", fmt.Sprintf("%s", obj))

	obj = NewPosition(false, true, nil, nil)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, true, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "home", fmt.Sprintf("%s", obj))

	obj = NewPosition(false, false, nil, nil)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "uninitialized", fmt.Sprintf("%s", obj))

	square := 5
	obj = NewPosition(false, false, nil, &square)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Equal(t, &square, obj.Square())
	assert.Equal(t, "square 5", fmt.Sprintf("%s", obj))

	safe := 10
	obj = NewPosition(false, false, &safe, nil)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Equal(t, &safe, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "safe 10", fmt.Sprintf("%s", obj))
}

func TestPositionCopy(t *testing.T) {
	var obj Position
	var copied Position

	obj = NewPosition(true, false, nil, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	obj = NewPosition(false, true, nil, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	obj = NewPosition(false, false, nil, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	square := 5
	obj = NewPosition(false, false, nil, &square)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	safe := 10
	obj = NewPosition(false, false, &safe, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPositionMoveToPositionValidStart(t *testing.T) {
	target := NewPosition(true, false, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "start", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionValidHome(t *testing.T) {
	target := NewPosition(false, true, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "home", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionValidSafe(t *testing.T) {
	safe := 3
	target := NewPosition(false, false, &safe, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "safe 3", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionValidSquare(t *testing.T) {
	square := 10
	target := NewPosition(false, false, nil, &square)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "square 10", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionInvalidMultiple(t *testing.T) {
	one := 1
	for _, target := range []Position {
		NewPosition(true, true, nil, nil),
		NewPosition(true, false, &one, nil),
		NewPosition(true, false, nil, &one),
		NewPosition(false, true, &one, nil),
		NewPosition(false, true, nil, &one),
		NewPosition(false, false, &one, &one),
	} {
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToPosition(target)
		assert.EqualError(t, err, "invalid position")
	}
}

func TestPositionMoveToPositionInvalidNone(t *testing.T) {
	target := NewPosition(false, false, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.EqualError(t, err, "invalid position")
}

func TestPositionMoveToPositionInvalidSafe(t *testing.T) {
	for _, safe := range []int {-1000, -2, -1, 5, 6, 1000, } {
		target := NewPosition(false, false, &safe, nil)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToPosition(target)
		assert.EqualError(t, err, "invalid safe square")
	}
}

func TestPositionMoveToPositionInvalidSquare(t *testing.T) {
	for _, square := range []int { -1000, -2, -1, 60, 61, 1000, } {
		target := NewPosition(false, false, nil, &square)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToPosition(target)
		assert.EqualError(t, err, "invalid square")
	}
}

func TestPositionMoveToStart(t *testing.T) {
	expected := NewPosition(true, false, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToStart()
	assert.Nil(t, err)
	assert.Equal(t, expected, position)
}

func TestPositionMoveToHome(t *testing.T) {
	expected := NewPosition(false, true, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToHome()
	assert.Nil(t, err)
	assert.Equal(t, expected, position)
}

func TestPositionMoveToSafeValid(t *testing.T) {
	for safe := 0; safe < SafeSquares; safe++ {
		expected := NewPosition(false, false, &safe, nil)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSafe(safe)
		assert.Nil(t, err)
		assert.Equal(t, expected, position)
	}
}

func TestPositionMoveToSafeInvalid(t *testing.T) {
	for _, safe := range []int {-1000, -2, -1, 5, 6, 1000, } {
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSafe(safe)
		assert.EqualError(t, err, "invalid safe square")
	}
}

func TestPositionMoveToSquareValid(t *testing.T) {
	for square := 0; square < BoardSquares; square++ {
		expected := NewPosition(false, false, nil, &square)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSquare(square)
		assert.Nil(t, err)
		assert.Equal(t, expected, position)
	}
}

func TestPositionMoveToSquareInvalid(t *testing.T) {
	for _, square := range []int { -1000, -2, -1, 60, 61, 1000, } {
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSquare(square)
		assert.EqualError(t, err, "invalid square")
	}
}

func TestNewPawn(t *testing.T) {
	obj := NewPawn(Red, 13)
	assert.Equal(t, Red, obj.Color())
	assert.Equal(t, 13, obj.Index())
	assert.Equal(t, "Red13", obj.Name())
	assert.Equal(t, NewPosition(true, false, nil, nil), obj.Position())
	assert.Equal(t, "Red13->start", fmt.Sprintf("%s", obj))
}

func TestPawnCopy(t *testing.T) {
	obj := NewPawn(Red, 13)
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPawnSetPosition(t *testing.T) {
	obj := NewPawn(Red, 13)
	target := NewPosition(false, true, nil, nil)
	obj.SetPosition(target)
	assert.Equal(t, target, obj.Position())
	assert.Equal(t, "Red13->home", fmt.Sprintf("%s", obj))
}

func TestNewPlayer(t *testing.T) {
	obj := NewPlayer(Red)
	assert.Equal(t, Red, obj.Color())
	assert.Equal(t, 0, obj.Turns())
	assert.Equal(t, Pawns, len(obj.Pawns()))
}

func TestPlayerCopy(t *testing.T) {
	var err error

	card1 := NewCard("0", CardApologies)

	obj := NewPlayer(Red)
	obj.AppendToHand(card1)
	err = obj.Pawns()[0].Position().MoveToHome()
	assert.Nil(t, err)
	err = obj.Pawns()[1].Position().MoveToSafe(2)
	assert.Nil(t, err)
	err = obj.Pawns()[2].Position().MoveToSquare(32)
	assert.Nil(t, err)

	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPlayerPublicData(t *testing.T) {
	var err error

	card1 := NewCard("0", CardApologies)

	obj := NewPlayer(Red)
	obj.AppendToHand(card1)
	err = obj.Pawns()[0].Position().MoveToHome()
	assert.Nil(t, err)
	err = obj.Pawns()[1].Position().MoveToSafe(2)
	assert.Nil(t, err)
	err = obj.Pawns()[2].Position().MoveToSquare(32)
	assert.Nil(t, err)

	expected := obj.Copy()
	expected.RemoveFromHand(card1) // the hand is cleared

	public := obj.PublicData()
	assert.Equal(t, expected, public)
	assert.NotSame(t, obj, public)
}

func TestPlayerAppendAndRemoveHand(t *testing.T) {
	obj := NewPlayer(Red)

	card1 := NewCard("1", Card1)
	card2 := NewCard("2", Card2)
	card3 := NewCard("3", Card3)

	// remove is idempotent; if it's not there, then that's ok
	obj.RemoveFromHand(card1)
	obj.RemoveFromHand(card2)
	obj.RemoveFromHand(card3)

	obj.AppendToHand(card1)
	assert.Equal(t, []Card{card1}, obj.Hand())

	obj.AppendToHand(card2)
	assert.Equal(t, []Card{card1, card2}, obj.Hand())

	obj.AppendToHand(card3)
	assert.Equal(t, []Card{card1, card2, card3}, obj.Hand())

	obj.RemoveFromHand(card2)
	assert.Equal(t, []Card{card1, card3}, obj.Hand())

	obj.RemoveFromHand(card3)
	assert.Equal(t, []Card{card1}, obj.Hand())

	obj.RemoveFromHand(card1)
	assert.Equal(t, []Card{}, obj.Hand())
}

func TestPlayerFindFirstPawnInStart(t *testing.T) {
	obj := NewPlayer(Red)

	for i := 0; i < Pawns; i++ {
		assert.Same(t, obj.Pawns()[i], *obj.FindFirstPawnInStart())
		err := obj.Pawns()[i].Position().MoveToHome()
		assert.Nil(t, err)
	}

	assert.Nil(t, obj.FindFirstPawnInStart())
}

func TestAllPawnsInHome(t *testing.T) {
	obj := NewPlayer(Red)

	for i := 0; i < Pawns; i++ {
		assert.False(t, obj.AllPawnsInHome())
		err := obj.Pawns()[i].Position().MoveToHome()
		assert.Nil(t, err)
	}

	assert.True(t, obj.AllPawnsInHome())
}

func TestIncrementTurns(t *testing.T) {
	obj := NewPlayer(Red)
	assert.Equal(t, 0, obj.Turns())
	obj.IncrementTurns()
	assert.Equal(t, 1, obj.Turns())
	obj.IncrementTurns()
	obj.IncrementTurns()
	assert.Equal(t, 3, obj.Turns())
}

func TestNewHistory(t *testing.T) {
	var obj History

	obj = NewHistory("action", nil, nil)
	assert.Equal(t, "action", obj.Action())
	assert.Nil(t, obj.Color())
	assert.Nil(t, obj.Card())
	assert.Equal(t, timestamp.GetStubbedTime(), obj.Timestamp())
	assert.Equal(t, fmt.Sprintf("[%s] General - action", timestamp.StubbedTime), fmt.Sprintf("%s", obj))

	color := Blue
	obj = NewHistory("action", &color, nil)
	assert.Equal(t, &color, obj.Color())
	assert.Nil(t, obj.Card())
	assert.Equal(t, timestamp.GetStubbedTime(), obj.Timestamp())
	assert.Equal(t, fmt.Sprintf("[%s] Blue - action", timestamp.StubbedTime), fmt.Sprintf("%s", obj))

	card1 := Card12
	obj = NewHistory("action", nil, &card1)
	assert.Nil(t, obj.Color())
	assert.Equal(t, &card1, obj.Card())
	assert.Equal(t, timestamp.GetStubbedTime(), obj.Timestamp())
	assert.Equal(t, fmt.Sprintf("[%s] General - action", timestamp.StubbedTime), fmt.Sprintf("%s", obj))
}

func TestHistoryCopy(t *testing.T) {
	color := Blue
	card1 := Card12
	obj := NewHistory("action", &color, &card1)
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestNewPlayerView(t *testing.T) {
	player1 := NewPlayer(Blue)
	player2 := NewPlayer(Red)

	opponents := make(map[PlayerColor]Player, 1)
	opponents[Red] = player2

	obj := NewPlayerView(player1, opponents)
	assert.Equal(t, player1, obj.Player())
	assert.Equal(t, opponents, obj.Opponents())
	assert.Equal(t, player2, obj.Opponents()[Red])
}

func TestPlayerViewCopy(t *testing.T) {
	player1 := NewPlayer(Blue)
	player2 := NewPlayer(Red)

	opponents := make(map[PlayerColor]Player, 1)
	opponents[Red] = player2

	obj := NewPlayerView(player1, opponents)
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPlayerViewGetPawn(t *testing.T) {
	player := NewPlayer(Red)
	opponents := map[PlayerColor]Player{Green: NewPlayer(Green)}
	view := NewPlayerView(player, opponents)
	assert.Equal(t, &view.Player().Pawns()[3], view.GetPawn(NewPawn(Red, 3)))
	assert.Equal(t, &view.Opponents()[Green].Pawns()[1], view.GetPawn(NewPawn(Green, 1)))
	assert.Nil(t, view.GetPawn(NewPawn(Yellow, 0)))
}

func TestPlayerViewAllPawns(t *testing.T) {
	player := NewPlayer(Red)
	opponents := map[PlayerColor]Player{Green: NewPlayer(Green)}
	view := NewPlayerView(player, opponents)
	pawns := view.AllPawns()
	assert.Equal(t, 2 * Pawns, len(pawns))
	for i := 0; i < Pawns; i++ {
		assert.True(t, slices.Contains(pawns, player.Pawns()[i]))
		assert.True(t, slices.Contains(pawns, opponents[Green].Pawns()[i]))
	}
}