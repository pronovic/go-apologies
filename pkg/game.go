package pkg

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/pkg/util/enum"
	"math/big"
	"strconv"
	"time"
)

// MinPlayers a game consists of at least 2 players
const MinPlayers = 2

// MaxPlayers a game consists of no more than 4 players
const MaxPlayers = 4

// Pawns there are 4 pawns per player, numbered 0-3
const Pawns = 4

// SafeSquares there are 5 safe squares for each color, numbered 0-4
const SafeSquares = 5

// BoardSquares there are 60 squares around the outside of the board, numbered 0-59
const BoardSquares = 60

// GameMode defines legal game modes
type GameMode struct{ value string }

// Value implements the enum.Enum interface for GameMode.
func (e GameMode) Value() string { return e.value }

var Standard = GameMode{"Standard"}
var Adult = GameMode{"Adult"}

// GameModes is the list of all legal GameMode enumerations
var GameModes = enum.NewValues[GameMode](Standard, Adult)

// PlayerColor defines all legal player colors
type PlayerColor struct{ value string }

// Value implements the enum.Enum interface for PlayerColor.
func (e PlayerColor) Value() string { return e.value }

var Red = PlayerColor{"Red"}
var Yellow = PlayerColor{"Yellow"}
var Blue = PlayerColor{"Blue"}
var Green = PlayerColor{"Green"}

// PlayerColors is the list of all legal PlayerColor enumerations, in order of use
var PlayerColors = enum.NewValues[PlayerColor](Red, Yellow, Blue, Green)

// CardType defines all legal types of cards
// The "A" card (CardApologies) is like the "Sorry" card in the original game
type CardType struct{ value string }

// Value implements the enum.Enum interface for CardType.
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

// CardTypes is the list of all legal CardType enumerations
var CardTypes = enum.NewValues[CardType](Card1, Card2, Card3, Card4, Card5, Card7, Card8, Card10, Card11, Card12, CardApologies)

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

	// Id Unique identifier for this card
	Id() string

	// Type The type of the card
	Type() CardType
}

type card struct {
	id       string
	cardType CardType
}

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

// Deck The deck of cards associated with a game.
type Deck interface {
	Draw() (Card, error)
	Discard(card Card) error
}

type deck struct {
	drawPile map[string]Card
	discardPile map[string]Card
}

func NewDeck() Deck {
	var drawPile = make(map[string]Card)
	var discardPile = make(map[string]Card)

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

func (p *deck) Draw() (Card, error) {
	if len(p.drawPile) < 1 {
		// this is equivalent to shuffling the discard pile into the draw pile
		for id, card := range p.discardPile {
			delete(p.discardPile, id)
			p.drawPile[id] = card
		}
	}

	if len(p.drawPile) < 1 {
		// in any normal game, this should never happen
		return *new(Card), errors.New("no cards available in deck")
	}

	keys := make([]string, 0, len(p.drawPile))
	for k := range p.drawPile {
		keys = append(keys, k)
	}

	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(keys))))
	if err != nil {
		return *new(Card), errors.New("failed to generate random int for draw")
	}

	key := keys[int(index.Int64())]
	card, _ := p.drawPile[key]
	delete(p.drawPile, key)

	return card, nil
}

func (p *deck) Discard(card Card) error {
	_, inDrawPile := p.drawPile[card.Id()]
	_, inDiscardPile := p.discardPile[card.Id()]

	if inDrawPile || inDiscardPile {
		return errors.New("card already exists in deck")
	}

	p.discardPile[card.Id()] = card
	return nil
}

// Position is the position of a pawn on the board.
type Position interface {

	// Start Whether this pawn resides in its start area
	Start() bool

	// Home Whether this pawn resides in its home area
	Home() bool

	// Safe Zero-based index of the square in the safe area where this pawn resides
	Safe() *int // optional

	// Square Zero-based index of the square on the board where this pawn resides
	Square() *int // optional

	// Copy Return a fully-independent copy of the position.
	Copy() Position

	// MoveToPosition Move the pawn to a specific position on the board.
	MoveToPosition(position Position) error

	// MoveToStart Move the pawn back to its start area.
	MoveToStart() error

	// MoveToHome Move the pawn to its home area.
	MoveToHome() error

	// MoveToSafe Move the pawn to a square in its safe area.
	MoveToSafe(square int) error

	// MoveToSquare Move the pawn to a square on the board.
	MoveToSquare(square int) error
}

type position struct {
	start bool
	home bool
	safe *int
	square *int
}

func NewPosition(start bool, home bool, safe *int, square *int) Position {
	return &position{
		start: start,
		home: home,
		safe: safe,
		square: square,
	}
}

func (p *position) Start() bool {
	return p.start
}

func (p *position) Home() bool {
	return p.home
}

func (p *position) Safe() *int {
	return p.safe
}

func (p *position) Square() *int {
	return p.square
}

func (p *position) Copy() Position {
	return &position {
		start: p.start,
		home: p.home,
		safe: p.safe,
		square: p.square,
	}
}

func (p *position) MoveToPosition(position Position) error {
	var fields = 0

	if p.Start() {
		fields += 1
	}

	if p.Home() {
		fields += 1
	}

	if p.Safe() != nil {
		fields += 1
	}

	if p.Square() != nil {
		fields += 1
	}

	if fields != 1 {
		return errors.New("invalid position")
	}

	if position.Start() {
		return p.MoveToStart()
	} else if position.Home() {
		return p.MoveToHome()
	} else if position.Safe() != nil {
		return p.MoveToSafe(*position.Safe())
	} else if position.Square() != nil {
		return p.MoveToSquare(*position.Square())
	} else {
		return errors.New("invalid position")
	}
}

func (p *position) MoveToStart() error {
	p.start = true
	p.home = false
	p.safe = nil
	p.square = nil

	return nil
}

func (p *position) MoveToHome() error {
	p.start = false
	p.home = true
	p.safe = nil
	p.square = nil

	return nil
}

func (p *position) MoveToSafe(square int) error {
	if square < 0 || square >= SafeSquares {
		return errors.New("invalid square")
	}

	p.start = false
	p.home = false
	p.safe = &square
	p.square = nil

	return nil
}

func (p *position) MoveToSquare(square int) error {
	if square < 0 || square >= BoardSquares {
		return errors.New("invalid square")
	}

	p.start = false
	p.home = false
	p.safe = nil
	p.square = &square

	return nil
}

func (p *position) String() string {
	if p.home {
		return "home"
	} else if p.start {
		return "start"
	} else if p.safe != nil {
		return fmt.Sprintf("safe %d", p.safe)
	} else {
		return fmt.Sprintf("square %s", p.square)
	}
}

// Pawn is a pawn on the board, belonging to a player.
type Pawn interface {

	// Color the color of this pawn
	Color() PlayerColor

	// Index Zero-based index of this pawn for a given user
	Index() int

	// Name The full name of this pawn as "colorindex"
	Name() string

	// Position The position of this pawn on the board
	Position() Position
}

type pawn struct {
	color    PlayerColor
	index    int
	name     string
	position Position
}

func NewPawn(color PlayerColor, index int) Pawn {
	return &pawn{
		color: color,
		index: index,
		name: fmt.Sprintf("%s%d", color, index),
		position: *new(Position),
	}
}

func (p *pawn) Color() PlayerColor {
	return p.color
}

func (p *pawn) Index() int {
	return p.index
}

func (p *pawn) Name() string {
	return p.name
}

func (p *pawn) Position() Position {
	return p.position
}

func (p *pawn) SetPosition(position Position) {
	p.position = position
}

func (p *pawn) String() string {
	return fmt.Sprintf("%s->%s", p.name, p.position)
}

// Player A player, which has a color and a set of pawns.
type Player interface {

	// Color the color of this player
	Color() PlayerColor

	// Hand List of cards in the player's hand
	Hand() []Card

	// Pawns List of all pawns belonging to the player
	Pawns() []Pawn

	// Turns number of turns for this player
	Turns() int

	// Copy Return a fully-independent copy of the player.
	Copy() Player

	// PublicData Return a fully-independent copy of the player with only public data visible.
	PublicData() Player

	// FindFirstPawnInStart Find the first pawn in the start area, if any.
	FindFirstPawnInStart() *Pawn // optional

	// AllPawnsInHome Whether all of this user's pawns are in home.
	AllPawnsInHome() bool
}

type player struct {
	color    PlayerColor
	hand     []Card
	pawns    []Pawn
	turns    int
}

func NewPlayer(color PlayerColor) Player {
	return &player{
		color: color,
		hand: make([]Card, 0, DeckSize),
		pawns: make([]Pawn, 0, Pawns),
		turns: 0,
	}
}

func (p *player) Color() PlayerColor {
	return p.color
}

func (p *player) Hand() []Card {
	return p.hand
}

func (p *player) Pawns() []Pawn {
	return p.pawns
}

func (p *player) Turns() int {
	return p.turns
}

func (p *player) Copy() Player {
	// TODO: implement Copy()
	return *new(Player)
}

func (p *player) PublicData() Player {
	// TODO: implement PublicData()
	return *new(Player)
}

func (p *player) FindFirstPawnInStart() *Pawn { // optional
	// TODO: implement FindFirstPawnInStart()
	return nil
}

func (p *player) AllPawnsInHome() bool {
	// TODO: implement AllPawnsInHome
	return false
}

// History Tracks an action taken during the game.
type History interface {

	// Action String describing the action
	Action() string

	// Color Color of the player associated with the action
	Color() *PlayerColor  // optional

	// Card Card associated with the action
	Card() *CardType	// optional

	// Timestamp Timestamp tied to the action (defaults to current time)
	Timestamp() time.Time

}

type history struct {
	action string
	color *PlayerColor
	card *CardType
	timestamp time.Time
}

func NewHistory(action string, color *PlayerColor, card *CardType) History {
	return &history{
		action: action,
		color: color,
		card: card,
		timestamp: time.Now().UTC(),
	}
}

func (h *history) Action() string {
	return h.action
}

func (h *history) Color() *PlayerColor { // optional
	return h.color
}

func (h *history) Card() *CardType { // optional
	return h.card
}

func (h *history) Timestamp() time.Time {
	return h.timestamp
}

// PlayerView A player-specific view of the game, showing only the information a player would have available on their turn.
type PlayerView interface {

	// Player The player associated with the view.
	Player() Player

	// Opponents The player's opponents, with private information stripped
	Opponents() map[Player]PlayerColor

	// Copy Return a fully-independent copy of the player view.
	Copy() PlayerView

	// GetPawn Return the pawn from this view with the same color and index, if any
	GetPawn(prototype Pawn) *Pawn

	// AllPawns Return a list of all pawns on the board.
	AllPawns() []Pawn
}

type playerView struct {
	player Player
	opponents map[Player]PlayerColor
}

func NewPlayerView(player Player, opponents map[Player]PlayerColor) PlayerView {
	return &playerView{
		player: player,
		opponents: opponents,
	}
}

func (v *playerView) Player() Player {
	return v.player
}

func (v *playerView) Opponents() map[Player]PlayerColor {
	return v.opponents
}

func (v *playerView) Copy() PlayerView {
	return *new(PlayerView) // TODO: implement Copy()
}

func (v *playerView) GetPawn(prototype Pawn) *Pawn {
	return new(Pawn) // TODO: implement GetPawn()
}

func (v *playerView) AllPawns() []Pawn {
	return make([]Pawn, 0) // TODO: implement AllPawns()
}

// Game The game, consisting of state for a set of players.
type Game interface {

	// PlayerCount Number of players in the game
	PlayerCount() int

	// Players All players in the game
	Players() map[Player]PlayerColor

	// Deck The deck of cards for the game
	Deck() Deck

	// History Game history
	History() []History

	// Started Whether the game has been started.
	Started() bool

	// Completed Whether the game is completed.
	Completed() bool

	// Winner The winner of the game, if any.
	Winner() *Player

	// Copy Return a fully-independent copy of the game.
	Copy() Game

	// Track Tracks an action taken during the game.
	Track(action string, player *Player, card *Card)

	// CreatePlayerView Return a player-specific view of the game, showing only the information a player would have available on their turn.
	CreatePlayerView(color PlayerColor) PlayerView

}

type game struct {
	playerCount int
	players map[Player]PlayerColor
	deck Deck
	history []History
}

func NewGame(playerCount int, players map[Player]PlayerColor, deck Deck) Game {
	return &game{
		playerCount: playerCount,
		players: players,
		deck: deck,
		history: make([]History, 0),
	}
}

func (g *game) PlayerCount() int {
	return g.playerCount
}

func (g *game) Players() map[Player]PlayerColor {
	return g.players
}

func (g *game) Deck() Deck {
	return g.deck
}

func (g *game) History() []History {
	return g.history
}

func (g *game) Started() bool {
	return false // TODO: implement Started()
}

func (g *game) Completed() bool {
	return false // TODO: implement Completed()
}

func (g *game) Winner() *Player {
	return nil // TODO: implement Winner()
}

func (g *game) Copy() Game {
	return *new(Game) // TODO : implement Copy()
}

func (g *game) Track(action string, player *Player, card *Card) {
	// TODO: implement Track()
}

func (g *game) CreatePlayerView(color PlayerColor) PlayerView {
	return *new(PlayerView) // implement CreatePlayerView()
}