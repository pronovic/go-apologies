package model

import (
	"errors"
	"fmt"
	"github.com/pronovic/go-apologies/internal/equality"
)

// SafeSquares there are 5 safe squares for each color, numbered 0-4
const SafeSquares = 5

// BoardSquares there are 60 squares around the outside of the board, numbered 0-59
const BoardSquares = 60

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

// Position is the position of a pawn on the board.
type Position interface {

	equality.EqualsByValue[Position]  // This interface implements equality by value

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
	MoveToSafe(safe int) error

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

// emptyPosition creates a new empty position in the start, for internal use
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

func (p *position) Equals(other Position) bool {
	return p.start == other.Start() &&
		p.home == other.Home() &&
		equality.IntPointerEquals(p.safe, other.Safe()) &&
		equality.IntPointerEquals(p.square, other.Square())
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

func (p *position) MoveToSafe(safe int) error {
	if safe < 0 || safe >= SafeSquares {
		return errors.New("invalid safe square")
	}

	p.start = false
	p.home = false
	p.safe = &safe
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

	equality.EqualsByValue[Pawn]  // This interface implements equality by value

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

func (p *pawn) Equals(other Pawn) bool {
	return p.color == other.Color() &&
		p.index == other.Index() &&
		p.name == other.Name() &&
		equality.ByValueEquals[Position](p.position, other.Position())
}

func (p *pawn) SetPosition(position Position) {
	p.position = position
}

func (p *pawn) String() string {
	return fmt.Sprintf("%s->%s", p.name, p.position)
}
