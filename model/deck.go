package model

import (
	"crypto/rand"
	"errors"
	"github.com/pronovic/go-apologies/internal/enum"
	"github.com/pronovic/go-apologies/internal/equality"
	"math/big"
	"strconv"
)

// AdultHand for an adult-mode game, we deal out 5 cards
const AdultHand = 5

// CardType defines all legal types of cards
// The "A" card (CardApologies) is like the "Sorry" card in the original game
type CardType struct{ value string }
func (e CardType) Value() string { return e.value }
func (e CardType) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
func (e *CardType) UnmarshalText(text []byte) error { return enum.Unmarshal(e, text, CardTypes) }
var CardTypes = enum.NewValues[CardType](Card1, Card2, Card3, Card4, Card5, Card7, Card8, Card10, Card11, Card12, CardApologies)
var Card1 = CardType{"1"}
var Card2 = CardType{"2"}
var Card3 = CardType{"3"}
var Card4 = CardType{"4"}
var Card5 = CardType{"5"}
var Card7 = CardType{"7"}
var Card8 = CardType{"8"}
var Card10 = CardType{"10"}
var Card11 = CardType{"11"}
var Card12 = CardType{"12"}
var CardApologies = CardType{"A"}

// DeckCounts defines the number of each type of card is in the deck
var DeckCounts = map[CardType]int{
	Card1:         5,
	Card2:         4,
	Card3:         4,
	Card4:         4,
	Card5:         4,
	Card7:         4,
	Card8:         4,
	Card10:        4,
	Card11:        4,
	Card12:        4,
	CardApologies: 4,
}

// DeckSize is the total size of the deck
var DeckSize = func(counts map[CardType]int) int {
	var total = 0
	for _, v := range counts {
		total += v
	}
	return total
}(DeckCounts)

// DrawAgain defines whether a given type of card draws again
var DrawAgain = map[CardType]bool{
	Card1:         false,
	Card2:         true,
	Card3:         false,
	Card4:         false,
	Card5:         false,
	Card7:         false,
	Card8:         false,
	Card10:        false,
	Card11:        false,
	Card12:        false,
	CardApologies: false,
}

// Card is a card in a deck or in a player's hand
type Card interface {

	equality.EqualsByValue[Card]  // This interface implements equality by value

	// Id Unique identifier for this card
	Id() string

	// Type The type of the card
	Type() CardType

	// Copy Return a fully-independent copy of the card.
	Copy() Card
}

type card struct {
	id       string
	cardType CardType
}

// NewCard constructs a new Card
func NewCard(id string, cardType CardType) Card {
	return &card{
		id: id,
		cardType: cardType,
	}
}

func (c *card) Id() string {
	return c.id
}

func (c *card) Type() CardType {
	return c.cardType
}

func (c *card) Copy() Card {
	return &card {
		id: c.id,
		cardType: c.cardType,
	}
}

func (c *card) Equals(other Card) bool {
	return other != nil &&
		c.id == other.Id() &&
		c.cardType == other.Type()
}

// Deck The deck of cards associated with a game.
type Deck interface {

	// Copy Return a fully-independent copy of the deck.
	Copy() Deck

	// Draw a card from the draw pile
	Draw() (Card, error)

	// Discard a card to the discard pile
	Discard(card Card) error

}

type deck struct {
	drawPile map[string]Card
	discardPile map[string]Card
}

// NewDeck constructs a new Deck
func NewDeck() Deck {
	var drawPile = make(map[string]Card, DeckSize)
	var discardPile = make(map[string]Card, DeckSize)

	var count = 0
	for _, c := range CardTypes.Members() {
		for i := 0; i < DeckCounts[c]; i++ {
			var id = strconv.Itoa(count)
			drawPile[id] = NewCard(id, c)
			count += 1
		}
	}

	return &deck{
		drawPile: drawPile,
		discardPile: discardPile,
	}
}

// Copy Return a fully-independent copy of the deck.
func (d *deck) Copy() Deck {
	var drawPileCopy = make(map[string]Card, DeckSize)
	for key := range d.drawPile {
		drawPileCopy[key] = d.drawPile[key].Copy()
	}

	var discardPileCopy = make(map[string]Card, DeckSize)
	for key := range d.discardPile {
		discardPileCopy[key] = d.discardPile[key].Copy()
	}

	return &deck{
		drawPile: drawPileCopy,
		discardPile: discardPileCopy,
	}
}

func (d *deck) Draw() (Card, error) {
	if len(d.drawPile) < 1 {
		// this is equivalent to shuffling the discard pile into the draw pile, because we draw randomly from the deck
		for id, card := range d.discardPile {
			delete(d.discardPile, id)
			d.drawPile[id] = card
		}
	}

	if len(d.drawPile) < 1 {
		// in any normal game, this should never happen
		return (Card)(nil), errors.New("no cards available in deck")
	}

	// because range on a map is not stable, the order of keys will vary
	keys := make([]string, 0, len(d.drawPile))
	for k := range d.drawPile {
		keys = append(keys, k)
	}

	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(keys))))
	if err != nil {
		return (Card)(nil), errors.New("failed to generate random int for draw")
	}

	key := keys[int(index.Int64())]
	card, _ := d.drawPile[key]
	delete(d.drawPile, key)

	return card, nil
}

func (d *deck) Discard(card Card) error {
	_, inDrawPile := d.drawPile[card.Id()]
	_, inDiscardPile := d.discardPile[card.Id()]

	if inDrawPile || inDiscardPile {
		return errors.New("card already exists in deck")
	}

	d.discardPile[card.Id()] = card
	return nil
}