package model

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/pronovic/go-apologies/internal/enum"
	"github.com/pronovic/go-apologies/internal/jsonutil"
	"github.com/pronovic/go-apologies/internal/randomutil"
)

// AdultHand for an adult-mode game, we deal out 5 cards
const AdultHand = 5

// CardType defines all legal types of cards
// The "A" card (CardApologies) is like the "Sorry" card in the original game
type CardType struct{ value string }

func (e CardType) Value() string                         { return e.value }
func (e CardType) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
func (e *CardType) UnmarshalText(text []byte) error      { return enum.Unmarshal(e, text, CardTypes) }

var (
	CardTypes     = enum.NewValues[CardType](Card1, Card2, Card3, Card4, Card5, Card7, Card8, Card10, Card11, Card12, CardApologies)
	Card1         = CardType{"1"}
	Card2         = CardType{"2"}
	Card3         = CardType{"3"}
	Card4         = CardType{"4"}
	Card5         = CardType{"5"}
	Card7         = CardType{"7"}
	Card8         = CardType{"8"}
	Card10        = CardType{"10"}
	Card11        = CardType{"11"}
	Card12        = CardType{"12"}
	CardApologies = CardType{"A"}
)

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
	total := 0
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
	// Id Unique identifier for this card
	Id() string

	// Type The type of the card
	Type() CardType

	// Copy Return a fully-independent copy of the card.
	Copy() Card
}

type card struct {
	Xid   string   `json:"id"`
	Xtype CardType `json:"type"`
}

// NewCard constructs a new Card
func NewCard(id string, cardType CardType) Card {
	return &card{
		Xid:   id,
		Xtype: cardType,
	}
}

// NewCardFromJSON constructs a new object from JSON in an io.Reader
func NewCardFromJSON(reader io.Reader) (Card, error) {
	return jsonutil.DecodeSimpleJSON[card](reader)
}

func (c *card) Id() string {
	return c.Xid
}

func (c *card) Type() CardType {
	return c.Xtype
}

func (c *card) Copy() Card {
	return &card{
		Xid:   c.Xid,
		Xtype: c.Xtype,
	}
}

func (c *card) Equals(other Card) bool {
	return other != nil &&
		c.Xid == other.Id() &&
		c.Xtype == other.Type()
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
	XdrawPile    map[string]Card `json:"draw"`
	XdiscardPile map[string]Card `json:"discard"`
}

// NewDeck constructs a new Deck
func NewDeck() Deck {
	drawPile := make(map[string]Card, DeckSize)
	discardPile := make(map[string]Card, DeckSize)

	count := 0
	for _, c := range CardTypes.Members() {
		for i := 0; i < DeckCounts[c]; i++ {
			id := strconv.Itoa(count)
			drawPile[id] = NewCard(id, c)
			count += 1
		}
	}

	return &deck{
		XdrawPile:    drawPile,
		XdiscardPile: discardPile,
	}
}

// NewDeckFromJSON constructs a new object from JSON in an io.Reader
func NewDeckFromJSON(reader io.Reader) (Deck, error) {
	type raw struct {
		XdrawPile    map[string]json.RawMessage `json:"draw"`
		XdiscardPile map[string]json.RawMessage `json:"discard"`
	}

	var temp raw
	err := json.NewDecoder(reader).Decode(&temp)
	if err != nil {
		return nil, err
	}

	var XdrawPile map[string]Card
	XdrawPile, err = jsonutil.DecodeMapJSON(temp.XdrawPile, NewCardFromJSON)
	if err != nil {
		return nil, err
	}

	var XdiscardPile map[string]Card
	XdiscardPile, err = jsonutil.DecodeMapJSON(temp.XdiscardPile, NewCardFromJSON)
	if err != nil {
		return nil, err
	}

	obj := deck{
		XdrawPile:    XdrawPile,
		XdiscardPile: XdiscardPile,
	}

	return &obj, nil
}

// Copy Return a fully-independent copy of the deck.
func (d *deck) Copy() Deck {
	drawPileCopy := make(map[string]Card, DeckSize)
	for key := range d.XdrawPile {
		drawPileCopy[key] = d.XdrawPile[key].Copy()
	}

	discardPileCopy := make(map[string]Card, DeckSize)
	for key := range d.XdiscardPile {
		discardPileCopy[key] = d.XdiscardPile[key].Copy()
	}

	return &deck{
		XdrawPile:    drawPileCopy,
		XdiscardPile: discardPileCopy,
	}
}

func (d *deck) Draw() (Card, error) {
	if len(d.XdrawPile) < 1 {
		// this is equivalent to shuffling the discard pile into the draw pile, because we draw randomly from the deck
		for id, card := range d.XdiscardPile {
			delete(d.XdiscardPile, id)
			d.XdrawPile[id] = card
		}
	}

	if len(d.XdrawPile) < 1 {
		// in any normal game, this should never happen
		return (Card)(nil), errors.New("no cards available in deck")
	}

	// because range on a map is not stable, the order of keys will vary
	keys := make([]string, 0, len(d.XdrawPile))
	for k := range d.XdrawPile {
		keys = append(keys, k)
	}

	index, err := randomutil.RandomInt(len(keys))
	if err != nil {
		return (Card)(nil), err
	}

	key := keys[index]
	c := d.XdrawPile[key]
	delete(d.XdrawPile, key)

	return c, nil
}

func (d *deck) Discard(card Card) error {
	_, inDrawPile := d.XdrawPile[card.Id()]
	_, inDiscardPile := d.XdiscardPile[card.Id()]

	if inDrawPile || inDiscardPile {
		return errors.New("card already exists in deck")
	}

	d.XdiscardPile[card.Id()] = card
	return nil
}
