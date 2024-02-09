package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/internal/enum"
	"github.com/pronovic/go-apologies/internal/timestamp"
	"io"
)

// GameMode defines legal game modes
type GameMode struct{ value string }
func (e GameMode) Value() string { return e.value }
func (e GameMode) MarshalText() (text []byte, err error) { return enum.Marshal(e) }
func (e *GameMode) UnmarshalText(text []byte) error { return enum.Unmarshal(e, text, GameModes) }
var GameModes = enum.NewValues[GameMode](AdultMode, StandardMode)
var StandardMode = GameMode{"StandardMode"}
var AdultMode = GameMode{"AdultMode"}

// History Tracks an action taken during the game.
type History interface {

	// Action String describing the action
	Action() string

	// Color Color of the player associated with the action
	Color() *PlayerColor  // optional

	// Card Card associated with the action
	Card() *CardType	// optional

	// Timestamp Timestamp tied to the action (defaults to current time)
	Timestamp() timestamp.Timestamp

	// Copy Return a fully-independent copy of the history.
	Copy() History

}

type history struct {
	Xaction   string `json:"action"`
	Xcolor    *PlayerColor `json:"color"`
	Xcard      *CardType `json:"card"`
	Xtimestamp timestamp.Timestamp `json:"timestamp"`
}

// NewHistory constructs a new History, optionally accepting a timestamp factory
func NewHistory(action string, color *PlayerColor, card *CardType, factory timestamp.Factory) History {
	if factory == nil {
		factory = timestamp.NewFactory()
	}

	return &history{
		Xaction:    action,
		Xcolor:     color,
		Xcard:      card,
		Xtimestamp: factory.CurrentTime(),
	}
}

// NewHistoryFromJSON constructs a new object from JSON in an io.Reader
func NewHistoryFromJSON(reader io.Reader) (History, error) {
	var obj history

	err := json.NewDecoder(reader).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func (h *history) Action() string {
	return h.Xaction
}

func (h *history) Color() *PlayerColor { // optional
	return h.Xcolor
}

func (h *history) Card() *CardType { // optional
	return h.Xcard
}

func (h *history) Timestamp() timestamp.Timestamp {
	return h.Xtimestamp
}

func (h *history) Copy() History {
	return &history{
		Xaction:    h.Xaction,
		Xcolor:     h.Xcolor,
		Xcard:      h.Xcard,
		Xtimestamp: h.Xtimestamp,
	}
}

func (h *history) String() string {
	now := h.Xtimestamp.Format()
	color := "General"
	if h.Xcolor != nil {
		color = h.Xcolor.Value()
	}
	action := h.Xaction
	return fmt.Sprintf("[%s] %s - %s", now, color, action)
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
	XplayerCount int  `json:"playercount"`
	Xplayers map[PlayerColor]Player `json:"players"`
	Xdeck    Deck      `json:"deck"`
	Xhistory []History `json:"history"`
	factory  timestamp.Factory
}

// NewGame constructs a new Game, optionally accepting a timestamp factory
func NewGame(playerCount int, factory timestamp.Factory) (Game, error) {
	if factory == nil {
		factory = timestamp.NewFactory()
	}

	if playerCount < MinPlayers || playerCount > MaxPlayers {
		return (*game)(nil), errors.New("invalid number of players")
	}

	players := make(map[PlayerColor]Player, playerCount)
	for i := 0; i < playerCount; i++ {
		color := PlayerColors.Members()[i]
		players[color] = NewPlayer(color)
	}

	game := &game{
		XplayerCount: playerCount,
		Xplayers:     players,
		Xdeck:        NewDeck(),
		Xhistory:     make([]History, 0),
		factory:      factory,
	}

	return game, nil
}

// NewGameFromJSON constructs a new object from JSON in an io.Reader
func NewGameFromJSON(reader io.Reader) (Game, error) {
	type raw struct {
		XplayerCount int  `json:"playercount"`
		Xplayers map[PlayerColor]json.RawMessage `json:"players"`
		Xdeck    json.RawMessage `json:"deck"`
		Xhistory []json.RawMessage `json:"history"`
	}

	var temp raw
	err := json.NewDecoder(reader).Decode(&temp)
	if err != nil {
		return nil, err
	}

	var Xplayers = make(map[PlayerColor]Player, len(temp.Xplayers))
	for key := range temp.Xplayers {
		value := temp.Xplayers[key]
		if value == nil || string(value) == "null" {
			Xplayers[key] = nil
		} else {
			element, err := NewPlayerFromJSON(bytes.NewReader(value))
			if err != nil {
				return nil, err
			}
			Xplayers[key] = element
		}
	}

	var Xdeck Deck
	if temp.Xdeck != nil || string(temp.Xdeck) == "null" {
		Xdeck, err = NewDeckFromJSON(bytes.NewReader(temp.Xdeck))
		if err != nil {
			return nil, err
		}
	}

	var Xhistory = make([]History, len(temp.Xhistory))
	for i := range temp.Xhistory {
		value := temp.Xhistory[i]
		if value == nil || string(value) == "null" {
			Xhistory[i] = nil
		} else {
			element, err := NewHistoryFromJSON(bytes.NewReader(value))
			if err != nil {
				return nil, err
			}
			Xhistory[i] = element
		}
	}

	obj := game {
		XplayerCount: temp.XplayerCount,
		Xplayers:     Xplayers,
		Xdeck:        Xdeck,
		Xhistory:     Xhistory,
		factory:      timestamp.NewFactory(),
	}

	return &obj, nil
}

func (g *game) PlayerCount() int {
	return g.XplayerCount
}

func (g *game) Players() map[PlayerColor]Player {
	return g.Xplayers
}

func (g *game) Deck() Deck {
	return g.Xdeck
}

func (g *game) History() []History {
	return g.Xhistory
}

func (g *game) Copy() Game {
	var playersCopy = make(map[PlayerColor]Player, len(g.Xplayers))

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		player, exists := g.Xplayers[color]
		if exists {
			playersCopy[color] = player.Copy()
		}
	}

	var historyCopy = make([]History, 0, len(g.Xhistory))
	for i := range g.Xhistory {
		historyCopy = append(historyCopy, g.Xhistory[i])
	}

	return &game{
		XplayerCount: g.XplayerCount,
		Xplayers:     playersCopy,
		Xdeck:        g.Xdeck.Copy(),
		Xhistory:     historyCopy,
		factory:      g.factory,
	}
}

func (g *game) Started() bool {
	return len(g.Xhistory) > 0 // if there is any history the game has been started
}

func (g *game) Completed() bool {
	for _, players := range g.Xplayers {
		if players.AllPawnsInHome() {
			return true
		}
	}

	return false
}

func (g *game) Winner() *Player {
	for _, players := range g.Xplayers {
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

	var history = NewHistory(action, color, cardtype, g.factory)
	g.Xhistory = append(g.Xhistory, history)

	if player != nil {
		g.Xplayers[player.Color()].IncrementTurns()
	}
}

func (g *game) CreatePlayerView(color PlayerColor) (PlayerView, error) {
	player, ok := g.Xplayers[color]
	if ! ok {
		return (PlayerView)(nil), errors.New("invalid color")
	}

	copied := player.Copy()

	opponents := make(map[PlayerColor]Player, len(g.Xplayers))

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		opponent, exists := g.Xplayers[color]
		if exists {
			if opponent.Color() != player.Color() {
				opponents[color] = opponent.PublicData()
			}
		}
	}

	return NewPlayerView(copied, opponents), nil
}

