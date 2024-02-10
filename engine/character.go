package engine

import (
	"github.com/pronovic/go-apologies/model"
	"github.com/pronovic/go-apologies/source"
)

// Character is a character that plays a game, which could be a person or could be computer-driven.
type Character interface {
	// Name the name of this character
	Name() string

	// Source The character input source from which moves are taken
	Source() source.CharacterInputSource

	// Color Get the color associated with this character
	Color() model.PlayerColor

	// SetColor Set the color associated with this character
	SetColor(color model.PlayerColor)

	// ChooseMove Choose the next move for a character via the input source
	ChooseMove(mode model.GameMode, view model.PlayerView, legalMoves []model.Move) (model.Move, error)
}

type character struct {
	name   string
	source source.CharacterInputSource
	color  model.PlayerColor
}

// NewCharacter constructs a new Character
func NewCharacter(name string, source source.CharacterInputSource) Character {
	return &character{
		name:   name,
		source: source,
		color:  *new(model.PlayerColor),
	}
}

func (c *character) Name() string {
	return c.name
}

func (c *character) Source() source.CharacterInputSource {
	return c.source
}

// Color Get the color associated with this character
func (c *character) Color() model.PlayerColor {
	return c.color
}

// SetColor Set the color associated with this character
func (c *character) SetColor(color model.PlayerColor) {
	c.color = color
}

func (c *character) ChooseMove(mode model.GameMode, view model.PlayerView, legalMoves []model.Move) (model.Move, error) {
	return c.source.ChooseMove(mode, view, legalMoves)
}
