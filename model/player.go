package model

import (
	"encoding/json"
	"io"
	"slices"

	"github.com/pronovic/go-apologies/internal/enum"
	"github.com/pronovic/go-apologies/internal/equality"
	"github.com/pronovic/go-apologies/internal/jsonutil"
)

// MinPlayers a game consists of at least 2 players
const MinPlayers = 2

// MaxPlayers a game consists of no more than 4 players
const MaxPlayers = 4

// Pawns there are 4 pawns per player, numbered 0-3
const Pawns = 4

// PlayerColor defines all legal player colors, enumerated in order of use
type PlayerColor struct{ value string }

func (e PlayerColor) Value() string                         { return e.value }
func (e PlayerColor) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
func (e *PlayerColor) UnmarshalText(text []byte) error      { return enum.Unmarshal(e, text, PlayerColors) }

var (
	PlayerColors = enum.NewValues[PlayerColor](Red, Yellow, Green, Blue)
	Red          = PlayerColor{"Red"}
	Yellow       = PlayerColor{"Yellow"}
	Blue         = PlayerColor{"Blue"}
	Green        = PlayerColor{"Green"}
)

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
	Xcolor PlayerColor `json:"color"`
	Xhand  []Card      `json:"hand"`
	Xpawns []Pawn      `json:"pawns"`
	Xturns int         `json:"turns"`
}

// NewPlayer constructs a new Player
func NewPlayer(color PlayerColor) Player {
	pawns := make([]Pawn, 0, Pawns)
	for i := 0; i < Pawns; i++ {
		pawns = append(pawns, NewPawn(color, i))
	}

	return &player{
		Xcolor: color,
		Xhand:  make([]Card, 0, DeckSize),
		Xpawns: pawns,
		Xturns: 0,
	}
}

// NewPlayerFromJSON constructs a new object from JSON in an io.Reader
func NewPlayerFromJSON(reader io.Reader) (Player, error) {
	type raw struct {
		Xcolor PlayerColor       `json:"color"`
		Xhand  []json.RawMessage `json:"hand"`
		Xpawns []json.RawMessage `json:"pawns"`
		Xturns int               `json:"turns"`
	}

	var temp raw
	err := json.NewDecoder(reader).Decode(&temp)
	if err != nil {
		return nil, err
	}

	var Xhand []Card
	Xhand, err = jsonutil.DecodeSliceJSON(temp.Xhand, NewCardFromJSON)
	if err != nil {
		return nil, err
	}

	var Xpawns []Pawn
	Xpawns, err = jsonutil.DecodeSliceJSON(temp.Xpawns, NewPawnFromJSON)
	if err != nil {
		return nil, err
	}

	obj := player{
		Xcolor: temp.Xcolor,
		Xhand:  Xhand,
		Xpawns: Xpawns,
		Xturns: temp.Xturns,
	}

	return &obj, nil
}

func (p *player) Color() PlayerColor {
	return p.Xcolor
}

func (p *player) Hand() []Card {
	return p.Xhand
}

func (p *player) Pawns() []Pawn {
	return p.Xpawns
}

func (p *player) Turns() int {
	return p.Xturns
}

func (p *player) Copy() Player {
	handCopy := make([]Card, 0, DeckSize)
	for i := range p.Xhand {
		handCopy = append(handCopy, p.Xhand[i].Copy())
	}

	pawnsCopy := make([]Pawn, 0, Pawns)
	for i := range p.Xpawns {
		pawnsCopy = append(pawnsCopy, p.Xpawns[i].Copy())
	}

	return &player{
		Xcolor: p.Xcolor,
		Xhand:  handCopy,
		Xpawns: pawnsCopy,
		Xturns: p.Xturns,
	}
}

func (p *player) PublicData() Player {
	handCopy := make([]Card, 0, DeckSize) // other players should not see this player's hand when making decisions

	pawnsCopy := make([]Pawn, 0, Pawns)
	for i := range p.Xpawns {
		pawnsCopy = append(pawnsCopy, p.Xpawns[i].Copy())
	}

	return &player{
		Xcolor: p.Xcolor,
		Xhand:  handCopy,
		Xpawns: pawnsCopy,
		Xturns: p.Xturns,
	}
}

func (p *player) AppendToHand(card Card) {
	p.Xhand = append(p.Xhand, card)
}

func (p *player) RemoveFromHand(card Card) {
	for i := 0; i < len(p.Xhand); i++ {
		found := p.Xhand[i]
		if equality.EqualByValue(card, found) {
			p.Xhand = slices.Delete(p.Xhand, i, i+1)
			return
		}
	}

	return
}

func (p *player) FindFirstPawnInStart() *Pawn { // optional
	for i := range p.Xpawns {
		if p.Xpawns[i].Position().Start() {
			return &p.Xpawns[i]
		}
	}

	return nil
}

func (p *player) AllPawnsInHome() bool {
	for i := range p.Xpawns {
		if !p.Xpawns[i].Position().Home() {
			return false
		}
	}

	return true
}

func (p *player) IncrementTurns() {
	p.Xturns += 1
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
	Xplayer    Player                 `json:"player"`
	Xopponents map[PlayerColor]Player `json:"opponents"`
}

// NewPlayerView contructs a new PlayerView
func NewPlayerView(player Player, opponents map[PlayerColor]Player) PlayerView {
	return &playerView{
		Xplayer:    player,
		Xopponents: opponents,
	}
}

// NewPlayerViewFromJSON constructs a new object from JSON in an io.Reader
func NewPlayerViewFromJSON(reader io.Reader) (PlayerView, error) {
	type raw struct {
		Xplayer    json.RawMessage                 `json:"player"`
		Xopponents map[PlayerColor]json.RawMessage `json:"opponents"`
	}

	var temp raw
	err := json.NewDecoder(reader).Decode(&temp)
	if err != nil {
		return nil, err
	}

	var Xplayer Player
	Xplayer, err = jsonutil.DecodeInterfaceJSON(temp.Xplayer, NewPlayerFromJSON)
	if err != nil {
		return nil, err
	}

	var Xopponents map[PlayerColor]Player
	Xopponents, err = jsonutil.DecodeMapJSON(temp.Xopponents, NewPlayerFromJSON)
	if err != nil {
		return nil, err
	}

	obj := playerView{
		Xplayer:    Xplayer,
		Xopponents: Xopponents,
	}

	return &obj, nil
}

func (v *playerView) Player() Player {
	return v.Xplayer
}

func (v *playerView) Opponents() map[PlayerColor]Player {
	return v.Xopponents
}

func (v *playerView) Copy() PlayerView {
	opponentsCopy := make(map[PlayerColor]Player, len(v.Xopponents))

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		opponent, exists := v.Xopponents[color]
		if exists {
			opponentsCopy[color] = opponent.Copy()
		}
	}

	return &playerView{
		Xplayer:    v.Xplayer.Copy(),
		Xopponents: opponentsCopy,
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

	for i := range v.Xplayer.Pawns() {
		all = append(all, v.Xplayer.Pawns()[i])
	}

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		opponent, exists := v.Xopponents[color]
		if exists {
			pawns := opponent.Pawns()
			for i := range pawns {
				all = append(all, pawns[i])
			}
		}
	}

	return all
}
