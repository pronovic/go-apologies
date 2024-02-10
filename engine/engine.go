package engine

import (
	"errors"
	"fmt"

	"github.com/pronovic/go-apologies/internal/circularqueue"
	"github.com/pronovic/go-apologies/internal/randomutil"
	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/rules"
)

// Engine Game engine that coordinates character actions in a game.
type Engine interface {
	// Mode The game mode
	Mode() model.GameMode

	// Characters The characters playing the game
	Characters() []Character

	// First The first player, chosen randomly
	First() model.PlayerColor

	// SetFirst Override the randomly-chosen first player
	SetFirst(first model.PlayerColor) error

	// Players The number of players in the game
	Players() int

	// State String describing the state of the game
	State() string

	// Game A reference to the underlying game
	Game() model.Game

	// Started Whether the game is started.
	Started() bool

	// Completed Whether the game is completed.
	Completed() bool

	// ColorMap Map from player color to character
	ColorMap() map[model.PlayerColor]Character

	// Winner Return the winner of the game
	Winner() Character

	// Reset Reset game state
	Reset() (model.Game, error)

	// StartGame Start the game, returning game state.
	StartGame() (model.Game, error)

	// NextTurn Get the color and character for the next turn
	// This will give you a different player each time you call it.
	NextTurn() (Character, error)

	// PlayNext Play the next turn of the game, returning game state as of the end of the turn.
	PlayNext() (model.Game, error)

	// Draw Draw a random card from the game's draw pile.
	Draw() (model.Card, error)

	// Discard Discard back to the game's discard pile.
	Discard(card model.Card) error

	// ConstructLegalMoves Construct the legal moves based on a player view, using the passed-in card if provided.
	ConstructLegalMoves(view model.PlayerView, card model.Card) (model.Card, []model.Move, error)

	// ChooseNextMove Choose the next move for a character based on a player view.
	ChooseNextMove(character Character, view model.PlayerView) (model.Move, error)

	// ExecuteMove Execute a move for a player, returning true if the player's turn is done.
	ExecuteMove(color model.PlayerColor, move model.Move) (bool, error)
}

type engine struct {
	mode       model.GameMode
	characters []Character
	evaluator  rules.Rules
	players    int
	colors     []model.PlayerColor
	first      model.PlayerColor
	queue      circularqueue.CircularQueue[model.PlayerColor]
	game       model.Game
	colorMap   map[model.PlayerColor]Character
}

// NewEngine constructs a new Engine
func NewEngine(mode model.GameMode, characters []Character, evaluator rules.Rules) (Engine, error) {
	if characters == nil || len(characters) < 1 {
		return nil, errors.New("at least one character required")
	}

	if evaluator == nil {
		evaluator = rules.NewRules(nil)
	}

	players := len(characters)
	colors := model.PlayerColors.Members()[0:players]

	first, err := randomutil.RandomChoice(colors)
	if err != nil {
		return nil, err
	}

	queue := circularqueue.NewCircularQueue(colors)
	err = queue.SetFirst(first)
	if err != nil {
		return nil, err
	}

	game, err := model.NewGame(players, nil)
	if err != nil {
		return nil, err
	}

	i := 0
	colorMap := make(map[model.PlayerColor]Character, players)
	for _, color := range model.PlayerColors.Members()[0:players] {
		c := characters[i]
		c.SetColor(color)
		colorMap[c.Color()] = c
		i += 1
	}

	constructed := &engine{
		mode:       mode,
		characters: characters,
		evaluator:  evaluator,
		players:    players,
		colors:     colors,
		first:      first,
		queue:      queue,
		game:       game,
		colorMap:   colorMap,
	}

	return constructed, nil
}

func (e *engine) Mode() model.GameMode {
	return e.mode
}

func (e *engine) Characters() []Character {
	return e.characters
}

func (e *engine) First() model.PlayerColor {
	return e.first
}

func (e *engine) SetFirst(first model.PlayerColor) error {
	e.first = first
	return e.queue.SetFirst(first)
}

func (e *engine) Players() int {
	return e.players
}

func (e *engine) State() string {
	if e.Completed() {
		return "Game completed"
	} else if e.Started() {
		return "Game in progress"
	} else {
		return "Game waiting to start"
	}
}

func (e *engine) Game() model.Game {
	return e.game
}

func (e *engine) Started() bool {
	return e.game.Started()
}

func (e *engine) Completed() bool {
	return e.game.Completed()
}

func (e *engine) ColorMap() map[model.PlayerColor]Character {
	return e.colorMap
}

func (e *engine) Winner() Character {
	winner := e.game.Winner()
	if winner == nil {
		return nil
	}

	color := (*winner).Color()
	return e.colorMap[color]
}

func (e *engine) Reset() (model.Game, error) {
	game, err := model.NewGame(e.players, nil)
	if err != nil {
		return nil, err
	}

	e.game = game
	return e.game, nil
}

func (e *engine) StartGame() (model.Game, error) {
	err := e.evaluator.StartGame(e.game, e.mode)
	if err != nil {
		return nil, err
	}

	return e.game, nil
}

func (e *engine) NextTurn() (Character, error) {
	color, err := e.queue.Next()
	if err != nil {
		return nil, err
	}

	return e.colorMap[color], nil
}

func (e *engine) PlayNext() (model.Game, error) {
	if e.Completed() {
		return e.game, errors.New("game is complete")
	}

	saved := e.game.Copy()

	next, err := e.NextTurn()
	if err != nil {
		e.game = saved // put back original so failed call is idempotent
		return nil, err
	}

	color := next.Color()

	done := false
	for {
		var err error
		var view model.PlayerView
		var move model.Move

		view, err = e.game.CreatePlayerView(color)
		if err != nil {
			e.game = saved // put back original so failed call is idempotent
			return nil, err
		}

		move, err = e.ChooseNextMove(next, view)
		if err != nil {
			e.game = saved // put back original so failed call is idempotent
			return nil, err
		}

		done, err = e.ExecuteMove(color, move)
		if err != nil {
			e.game = saved // put back original so failed call is idempotent
			return nil, err
		}

		if done {
			break
		}
	}

	return e.game, nil
}

func (e *engine) Draw() (model.Card, error) {
	return e.game.Deck().Draw()
}

func (e *engine) Discard(card model.Card) error {
	if card == nil {
		return errors.New("card is nil")
	}

	return e.Game().Deck().Discard(card)
}

// ConstructLegalMoves Construct the legal moves based on a player view, using the passed-in card if provided.
func (e *engine) ConstructLegalMoves(view model.PlayerView, card model.Card) (model.Card, []model.Move, error) {
	if view == nil {
		return nil, nil, errors.New("view is nil")
	}

	if e.mode == model.StandardMode {
		if card == nil {
			drawn, err := e.game.Deck().Draw()
			if err != nil {
				return nil, nil, err
			}

			card = drawn
		}
	}

	legalMoves, err := e.evaluator.ConstructLegalMoves(view, card)
	if err != nil {
		return nil, nil, err
	}

	return card, legalMoves, nil
}

// ChooseNextMove Choose the next move for a character based on a player view.
func (e *engine) ChooseNextMove(character Character, view model.PlayerView) (model.Move, error) {
	if character == nil {
		return nil, errors.New("character is nil")
	}

	if view == nil {
		return nil, errors.New("view is nil")
	}

	_, legalMoves, err := e.ConstructLegalMoves(view, nil)
	if err != nil {
		return nil, err
	}

	move, err := character.ChooseMove(e.mode, view, legalMoves)
	if err != nil {
		return nil, err
	}

	return move, nil
}

// ExecuteMove Execute a move for a player, returning true if the player's turn is done.
func (e *engine) ExecuteMove(color model.PlayerColor, move model.Move) (bool, error) {
	if move == nil {
		return false, errors.New("model is nil")
	}

	player := e.game.Players()[color]
	if e.mode == model.AdultMode {
		return e.executeMoveAdult(player, move)
	} else {
		return e.executeMoveStandard(player, move)
	}
}

func (e *engine) executeMoveStandard(player model.Player, move model.Move) (bool, error) {
	if len(move.Actions()) == 0 {
		e.game.Track(fmt.Sprintf("Turn is forfeit; discarded card %s", move.Card().Type()), player, move.Card())

		err := e.Discard(move.Card())
		if err != nil {
			return false, err
		}

		// player's turn is done if they forfeit
		return true, nil
	} else {
		// tracks history, potentially completes game
		err := e.evaluator.ExecuteMove(e.game, player, move)
		if err != nil {
			return false, err
		}

		err = e.Discard(move.Card())
		if err != nil {
			return false, err
		}

		// player's turn is done unless they can draw again with this card or the game is completed
		done := e.Completed() || !e.evaluator.DrawAgain(move.Card())

		return done, nil
	}
}

func (e *engine) executeMoveAdult(player model.Player, move model.Move) (bool, error) {
	if len(move.Actions()) == 0 {
		e.game.Track(fmt.Sprintf("Turn is forfeit; discarded card %s", move.Card().Type()), player, move.Card())

		player.RemoveFromHand(move.Card())

		drawn, err := e.Draw()
		if err != nil {
			return false, err
		}

		err = e.Discard(move.Card())
		if err != nil {
			return false, err
		}

		player.AppendToHand(drawn)

		// player's turn is done if they forfeit
		return true, nil
	} else {
		// tracks history, potentially completes game
		err := e.evaluator.ExecuteMove(e.game, player, move)
		if err != nil {
			return false, err
		}

		player.RemoveFromHand(move.Card())

		drawn, err := e.Draw()
		if err != nil {
			return false, err
		}

		err = e.Discard(move.Card())
		if err != nil {
			return false, err
		}

		player.AppendToHand(drawn)

		// player's turn is done unless they can draw again with this card or the game is completed
		done := e.Completed() || !e.evaluator.DrawAgain(move.Card())

		return done, nil
	}
}
