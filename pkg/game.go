package pkg

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/pkg/util/enum"
	"github.com/pronovic/go-apologies/pkg/util/timestamp"
	"math/big"
	"slices"
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
var PlayerColors = enum.NewValues[PlayerColor](Red, Yellow, Green, Blue)

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

// StartCircles defines the start circles for each color
var StartCircles = map[PlayerColor]Position {
	Red: newPositionAtSquare(4),
	Blue: newPositionAtSquare(19),
	Yellow: newPositionAtSquare(34),
	Green: newPositionAtSquare(49),
}

// TurnSquares defines the turn squares for each color, where forward movement turns into the safe zone
var TurnSquares = map[PlayerColor]Position {
	Red: newPositionAtSquare(2),
	Blue: newPositionAtSquare(17),
	Yellow: newPositionAtSquare(32),
	Green: newPositionAtSquare(47),
}

// Slides defines the start positions for each color
var Slides = map[PlayerColor][]Slide {
	Red: {newSlide(1, 4), newSlide(9, 13)},
	Blue: {newSlide(16, 19), newSlide(24, 28)},
	Yellow: {newSlide(31, 34), newSlide(39, 43)},
	Green: {newSlide(46, 49), newSlide(54, 58)},
}

// Slide defines the start and end positions of a slide on the board
type Slide interface {
	// Start is the start of the slide
	Start() int

	// End is the end of a the slide
	End() int
}

type slide struct {
	start int
	end int
}

// newSlide creates a new slide, for defining constants
func newSlide(start int, end int) Slide {
	return &slide{start, end }
}

func (s *slide) Start() int {
	return s.start
}

func (s *slide) End() int {
	return s.end
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
		// this is equivalent to shuffling the discard pile into the draw pile
		for id, card := range d.discardPile {
			delete(d.discardPile, id)
			d.drawPile[id] = card
		}
	}

	if len(d.drawPile) < 1 {
		// in any normal game, this should never happen
		return *new(Card), errors.New("no cards available in deck")
	}

	keys := make([]string, 0, len(d.drawPile))
	for k := range d.drawPile {
		keys = append(keys, k)
	}

	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(keys))))
	if err != nil {
		return *new(Card), errors.New("failed to generate random int for draw")
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

// NewPosition constructs a new Position
func NewPosition(start bool, home bool, safe *int, square *int) Position {
	return &position{
		start: start,
		home: home,
		safe: safe,
		square: square,
	}
}

// emptyPosition creates a new position in the start, for internal use
func emptyPosition() Position {
	return &position{
		start: true,
		home: false,
		safe: nil,
		square: nil,
	}
}

// newPositionAtSquare creates a new position at a particular square, for defining constants
func newPositionAtSquare(square int) Position {
	p := NewPosition(false, false, nil, nil)

	err := p.MoveToSquare(square)
	if err != nil {
		// panic is appropriate here, because this is used internally to set up constants, and if those are broken, we can't run
		panic("invalid square for new p")
	}

	return p
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

	if position.Start() {
		fields += 1
	}

	if position.Home() {
		fields += 1
	}

	if position.Safe() != nil {
		fields += 1
	}

	if position.Square() != nil {
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
		return errors.New("invalid safe square")
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
		return fmt.Sprintf("safe %v", *p.safe)
	} else if p.square != nil {
		return fmt.Sprintf("square %v", *p.square)
	} else {
		return "uninitialized"
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

	// SetPosition Set the position of this pawn on the board
	SetPosition(position Position)

	// Copy Return a fully-independent copy of the pawn.
	Copy() Pawn
}

type pawn struct {
	color    PlayerColor
	index    int
	name     string
	position Position
}

// NewPawn constructs a new Pawn
func NewPawn(color PlayerColor, index int) Pawn {
	return &pawn{
		color: color,
		index: index,
		name: fmt.Sprintf("%s%d", color.value, index),
		position: emptyPosition(),
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

func (p *pawn) Copy() Pawn {
	return &pawn{
		color: p.color,
		index: p.index,
		name: p.name,
		position: p.position.Copy(),
	}
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

	// AppendToHand appends a card into the player's hand
	AppendToHand(card Card)

	// RemoveFromHand removes a card from the player's hand
	RemoveFromHand(card Card)

	// FindFirstPawnInStart Find the first pawn in the start area, if any.
	FindFirstPawnInStart() *Pawn // optional

	// AllPawnsInHome Whether all of this user's pawns are in home.
	AllPawnsInHome() bool

	// IncrementTurns the number of turns for a player
	IncrementTurns()
}

type player struct {
	color    PlayerColor
	hand     []Card
	pawns    []Pawn
	turns    int
}

// NewPlayer constructs a new Player
func NewPlayer(color PlayerColor) Player {
	pawns := make([]Pawn, 0, Pawns)
	for i := 0; i < Pawns; i++ {
		pawns = append(pawns, NewPawn(color, i))
	}

	return &player{
		color: color,
		hand: make([]Card, 0, DeckSize),
		pawns: pawns,
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
	handCopy := make([]Card, 0, DeckSize)
	for i := range p.hand {
		handCopy = append(handCopy, p.hand[i].Copy())
	}

	pawnsCopy := make([]Pawn, 0, Pawns)
	for i := range p.pawns {
		pawnsCopy = append(pawnsCopy, p.pawns[i].Copy())
	}

	return &player{
		color: p.color,
		hand:  handCopy,
		pawns: pawnsCopy,
		turns: p.turns,
	}
}

func (p *player) PublicData() Player {
	handCopy := make([]Card, 0, DeckSize) // other players should not see this player's hand when making decisions

	pawnsCopy := make([]Pawn, 0, Pawns)
	for i := range p.pawns {
		pawnsCopy = append(pawnsCopy, p.pawns[i].Copy())
	}

	return &player{
		color: p.color,
		hand:  handCopy,
		pawns: pawnsCopy,
		turns: p.turns,
	}
}

func (p *player) AppendToHand(card Card) {
	p.hand = append(p.hand, card)
}

func (p *player) RemoveFromHand(card Card) {
	for i := 0; i < len(p.hand); i++ {
		found := p.hand[i]
		if found == card {
			p.hand = slices.Delete(p.hand, i, i+1)
			return
		}
	}

	return
}

func (p *player) FindFirstPawnInStart() *Pawn { // optional
	for i := range p.pawns {
		if p.pawns[i].Position().Start() {
			return &p.pawns[i]
		}
	}

	return nil
}

func (p *player) AllPawnsInHome() bool {
	for i := range p.pawns {
		if ! p.pawns[i].Position().Home() {
			return false
		}
	}

	return true
}

func (p *player) IncrementTurns() {
	p.turns += 1
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

	// Copy Return a fully-independent copy of the history.
	Copy() History

}

type history struct {
	action string
	color *PlayerColor
	card *CardType
	timestamp time.Time
}

// NewHistory constructs a new History
func NewHistory(action string, color *PlayerColor, card *CardType) History {
	return &history{
		action: action,
		color: color,
		card: card,
		timestamp: timestamp.CurrentTime(),
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

func (h *history) Copy() History {
	return &history{
		action: h.action,
		color: h.color,
		card: h.card,
		timestamp: h.timestamp,
	}
}

func (h *history) String() string {
	now := h.timestamp.Format(timestamp.Layout)
	color := "General"
	if h.color != nil {
		color = h.color.Value()
	}
	action := h.action
	return fmt.Sprintf("[%s] %s - %s", now, color, action)
}

// PlayerView A player-specific view of the game, showing only the information a player would have available on their turn.
type PlayerView interface {

	// Player The player associated with the view.
	Player() Player

	// Opponents The player's opponents, with private information stripped
	Opponents() map[PlayerColor]Player

	// Copy Return a fully-independent copy of the player view.
	Copy() PlayerView

	// GetPawn Return the pawn from this view with the same color and index, if any
	GetPawn(prototype Pawn) *Pawn

	// AllPawns Return a list of all pawns on the board.
	AllPawns() []Pawn
}

type playerView struct {
	player Player
	opponents map[PlayerColor]Player
}

// NewPlayerView contructs a new PlayerView
func NewPlayerView(player Player, opponents map[PlayerColor]Player) PlayerView {
	return &playerView{
		player: player,
		opponents: opponents,
	}
}

func (v *playerView) Player() Player {
	return v.player
}

func (v *playerView) Opponents() map[PlayerColor]Player {
	return v.opponents
}

func (v *playerView) Copy() PlayerView {
	opponentsCopy := make(map[PlayerColor]Player, len(v.opponents))
	for key := range v.opponents {
		opponentsCopy[key] = v.opponents[key].Copy()
	}

	return &playerView{
		player: v.player.Copy(),
		opponents: opponentsCopy,
	}
}

func (v *playerView) GetPawn(prototype Pawn) *Pawn {
	all := v.AllPawns()
	for i := range all {
		if all[i].Color() == prototype.Color() && all[i].Index() == prototype.Index() {
			return &all[i]
		}
	}

	return nil
}

func (v *playerView) AllPawns() []Pawn {
	total := 0
	total += len(v.player.Pawns())
	for key := range v.opponents {
		total += len(v.opponents[key].Pawns())
	}

	all := make([]Pawn, 0, total)
	for i := range v.player.Pawns() {
		all = append(all, v.player.Pawns()[i])
	}

	for key := range v.opponents {
		pawns := v.opponents[key].Pawns()
		for i := range pawns {
			all = append(all, pawns[i])
		}
	}

	return all
}

// Game The game, consisting of state for a set of players.
type Game interface {

	// PlayerCount Number of players in the game
	PlayerCount() int

	// Players All players in the game
	Players() map[PlayerColor]Player

	// Deck The deck of cards for the game
	Deck() Deck

	// History Game history
	History() []History

	// Copy Return a fully-independent copy of the game.
	Copy() Game

	// Started Whether the game has been started.
	Started() bool

	// Completed Whether the game is completed.
	Completed() bool

	// Winner The winner of the game, if any.
	Winner() *Player

	// Track Tracks an action taken during the game, optionally tracking player and/or card
	Track(action string, player Player, card Card)

	// CreatePlayerView Return a player-specific view of the game, showing only the information a player would have available on their turn.
	CreatePlayerView(color PlayerColor) (PlayerView, error)

}

type game struct {
	playerCount int
	players map[PlayerColor]Player
	deck Deck
	history []History
}

// NewGame constructs a new Game
func NewGame(playerCount int) (Game, error) {
	if playerCount < MinPlayers || playerCount > MaxPlayers {
		return *new(Game), errors.New("invalid number of players")
	}

	players := make(map[PlayerColor]Player, playerCount)
	for i := 0; i < playerCount; i++ {
		color := PlayerColors.Members()[i]
		players[color] = NewPlayer(color)
	}

	game := &game{
		playerCount: playerCount,
		players:     players,
		deck:        NewDeck(),
		history:     make([]History, 0),
	}

	return game, nil
}

func (g *game) PlayerCount() int {
	return g.playerCount
}

func (g *game) Players() map[PlayerColor]Player {
	return g.players
}

func (g *game) Deck() Deck {
	return g.deck
}

func (g *game) History() []History {
	return g.history
}

func (g *game) Copy() Game {
	var playersCopy = make(map[PlayerColor]Player, len(g.players))
	for key := range g.players {
		playersCopy[key] = g.players[key].Copy()
	}

	var historyCopy = make([]History, 0, len(g.history))
	for i := range g.history {
		historyCopy = append(historyCopy, g.history[i])
	}

	return &game{
		playerCount: g.playerCount,
		players: playersCopy,
		deck: g.deck.Copy(),
		history: historyCopy,
	}
}

func (g *game) Started() bool {
	return len(g.history) > 0 // if there is any history the game has been started
}

func (g *game) Completed() bool {
	for _, players := range g.players {
		if players.AllPawnsInHome() {
			return true
		}
	}

	return false
}

func (g *game) Winner() *Player {
	for _, players := range g.players {
		if players.AllPawnsInHome() {
			return &players
		}
	}

	return nil
}

func (g *game) Track(action string, player Player, card Card) {
	var color *PlayerColor = nil
	if player != nil {
		x := player.Color()
		color = &x
	}

	var cardtype *CardType = nil
	if card != nil {
		tmp := card.Type()
		cardtype = &tmp
	}

	var history = NewHistory(action, color, cardtype)
	g.history = append(g.history, history)

	if player != nil {
		g.players[player.Color()].IncrementTurns()
	}
}

func (g *game) CreatePlayerView(color PlayerColor) (PlayerView, error) {
	player, ok := g.players[color]
	if ! ok {
		return *new(PlayerView), errors.New("invalid color")
	}

	copied := player.Copy()

	opponents := make(map[PlayerColor]Player, len(g.players))
	for i := range g.players {
		if g.players[i].Color() != player.Color() {
			opponents[g.players[i].Color()] = g.players[i].PublicData()
		}
	}

	return NewPlayerView(copied, opponents), nil
}