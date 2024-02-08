package model

import (
	"github.com/pronovic/go-apologies/internal/enum"
	"github.com/pronovic/go-apologies/internal/equality"
	"slices"
)

// MinPlayers a game consists of at least 2 players
const MinPlayers = 2

// MaxPlayers a game consists of no more than 4 players
const MaxPlayers = 4

// Pawns there are 4 pawns per player, numbered 0-3
const Pawns = 4

// PlayerColor defines all legal player colors, enumerated in order of use
type PlayerColor struct{ value string }
func (e PlayerColor) Value() string { return e.value }
func (e PlayerColor) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
func (e *PlayerColor) UnmarshalText(text []byte) error { return enum.Unmarshal(e, text, PlayerColors) }
var PlayerColors = enum.NewValues[PlayerColor](Red, Yellow, Green, Blue)
var Red = PlayerColor{"Red"}
var Yellow = PlayerColor{"Yellow"}
var Blue = PlayerColor{"Blue"}
var Green = PlayerColor{"Green"}

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
		if equality.ByValueEquals[Card](card, found) {
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

// PlayerView A player-specific view of the game, showing only the information a player would have available on their turn.
type PlayerView interface {

	// Player The player associated with the view.
	Player() Player

	// Opponents The player's opponents, with private information stripped
	Opponents() map[PlayerColor]Player

	// Copy Return a fully-independent copy of the player view.
	Copy() PlayerView

	// GetPawn Return the pawn from this view with the same color and index, possibly nil
	GetPawn(prototype Pawn) Pawn

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

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		opponent, exists := v.opponents[color]
		if exists {
			opponentsCopy[color] = opponent.Copy()
		}
	}

	return &playerView{
		player: v.player.Copy(),
		opponents: opponentsCopy,
	}
}

func (v *playerView) GetPawn(prototype Pawn) Pawn {
	all := v.AllPawns()
	for i := range all {
		if all[i].Color() == prototype.Color() && all[i].Index() == prototype.Index() {
			return all[i]
		}
	}

	return nil
}

func (v *playerView) AllPawns() []Pawn {
	all := make([]Pawn, 0)

	for i := range v.player.Pawns() {
		all = append(all, v.player.Pawns()[i])
	}

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		opponent, exists := v.opponents[color]
		if exists {
			pawns := opponent.Pawns()
			for i := range pawns {
				all = append(all, pawns[i])
			}
		}
	}

	return all
}
