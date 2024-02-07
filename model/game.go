package model

import (
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/internal/timestamp"
	"time"
)

// GameMode defines legal game modes
type GameMode struct{ value string }

// Value implements the enum.Enum interface for GameMode.
func (e GameMode) Value() string { return e.value }

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

// NewHistory constructs a new History, optionally accepting a timestamp factory
func NewHistory(action string, color *PlayerColor, card *CardType, factory timestamp.Factory) History {
	if factory == nil {
		factory = timestamp.NewFactory()
	}

	return &history{
		action: action,
		color: color,
		card: card,
		timestamp: factory.CurrentTime(),
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
	factory timestamp.Factory
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
		playerCount: playerCount,
		players:     players,
		deck:        NewDeck(),
		history:     make([]History, 0),
		factory:     factory,
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

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		player, exists := g.players[color]
		if exists {
			playersCopy[color] = player.Copy()
		}
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
		factory: g.factory,
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

	var history = NewHistory(action, color, cardtype, g.factory)
	g.history = append(g.history, history)

	if player != nil {
		g.players[player.Color()].IncrementTurns()
	}
}

func (g *game) CreatePlayerView(color PlayerColor) (PlayerView, error) {
	player, ok := g.players[color]
	if ! ok {
		return (PlayerView)(nil), errors.New("invalid color")
	}

	copied := player.Copy()

	opponents := make(map[PlayerColor]Player, len(g.players))

	// range on a map explicitly does *not* return keys in a stable order, so we iterate on colors instead
	for _, color := range PlayerColors.Members() {
		opponent, exists := g.players[color]
		if exists {
			if opponent.Color() != player.Color() {
				opponents[color] = opponent.PublicData()
			}
		}
	}

	return NewPlayerView(copied, opponents), nil
}

