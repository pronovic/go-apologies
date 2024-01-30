package pkg


import (
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

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
