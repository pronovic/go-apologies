package main

import "github.com/pronovic/go-apologies/util/enum"

// MinPlayers a game consists of at least 2 players
const MinPlayers = 2

// MaxPlayers a game consists of no more than 4 players
const MaxPlayers = 4

// Pawns there are 4 pawns per player, numbered 0-3
const Pawns = 3

// SafeSquares there are 5 safe squares for each color, numbered 0-4
const SafeSquares = 5

// BoardSquares there are 60 squares around the outside of the board, numbered 0-59
const BoardSquares = 60

// GameMode defines legal game modes
type GameMode struct{ value string }

func (e GameMode) Value() string { return e.value }

var Standard = GameMode{"Standard"}
var Adult = GameMode{"Adult"}
var GameModes = enum.Values[GameMode](Standard, Adult)

// PlayerColor defines all legal player colors, in order of use
type PlayerColor struct{ value string }

func (e PlayerColor) Value() string { return e.value }

var Red = PlayerColor{"Red"}
var Yellow = PlayerColor{"Yellow"}
var Blue = PlayerColor{"Blue"}
var Green = PlayerColor{"Green"}
var PlayerColors = enum.Values[PlayerColor](Red, Yellow, Blue, Green)

// CardType defines all legal types of cards
// The "A" card (CardApologies) is like the "Sorry" card in the original game
type CardType struct{ value string }

func (e CardType) Value() string { return e.value }

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
var CardTypes = enum.Values[CardType](Card1, Card2, Card3, Card4, Card5, Card7, Card8, Card10, Card11, Card12, CardApologies)

// AdultHand for an adult-mode game, we deal out 5 cards
const AdultHand = 5

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
var DeckSize = func(m map[CardType]int) int {
	var total = 0
	for _, v := range m {
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
type Card struct {
	Id   string
	Type CardType
}
