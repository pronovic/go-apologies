package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSlide(t *testing.T) {
	obj := newSlide(1, 2)
	assert.Equal(t, 1, obj.Start())
	assert.Equal(t, 2, obj.End())
}

func TestNewPosition(t *testing.T) {
	var obj Position

	obj = NewPosition(true, false, nil, nil)
	assert.Equal(t, true, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "start", fmt.Sprintf("%s", obj))

	obj = NewPosition(false, true, nil, nil)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, true, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "home", fmt.Sprintf("%s", obj))

	obj = NewPosition(false, false, nil, nil)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "uninitialized", fmt.Sprintf("%s", obj))

	square := 5
	obj = NewPosition(false, false, nil, &square)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Nil(t, obj.Safe())
	assert.Equal(t, &square, obj.Square())
	assert.Equal(t, "square 5", fmt.Sprintf("%s", obj))

	safe := 10
	obj = NewPosition(false, false, &safe, nil)
	assert.Equal(t, false, obj.Start())
	assert.Equal(t, false, obj.Home())
	assert.Equal(t, &safe, obj.Safe())
	assert.Nil(t, obj.Square())
	assert.Equal(t, "safe 10", fmt.Sprintf("%s", obj))
}

func TestPositionCopy(t *testing.T) {
	var obj Position
	var copied Position

	obj = NewPosition(true, false, nil, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	obj = NewPosition(false, true, nil, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	obj = NewPosition(false, false, nil, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	square := 5
	obj = NewPosition(false, false, nil, &square)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)

	safe := 10
	obj = NewPosition(false, false, &safe, nil)
	copied = obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPositionEquals(t *testing.T) {
	safe := 1
	square := 2

	p1 := NewPosition(true, false, nil, nil)
	p2 := NewPosition(false, true, nil, nil)
	p3 := NewPosition(false, false, &safe, nil)
	p4 := NewPosition(false, false, nil, &square)

	// note: it is important to test with assert.True()/assert.False() and x.Equals(y)
	// because assert.Equals() and assert.NotEquals() are not aware of our equality by value concept

	assert.True(t, p1.Equals(p1))
	assert.True(t, p2.Equals(p2))
	assert.True(t, p3.Equals(p3))
	assert.True(t, p4.Equals(p4))

	assert.False(t, p1.Equals(nil))
	assert.False(t, p2.Equals(nil))
	assert.False(t, p3.Equals(nil))
	assert.False(t, p4.Equals(nil))

	assert.False(t, p1.Equals(p2))
	assert.False(t, p1.Equals(p3))
	assert.False(t, p1.Equals(p4))

	assert.False(t, p2.Equals(p1))
	assert.False(t, p2.Equals(p3))
	assert.False(t, p2.Equals(p4))

	assert.False(t, p3.Equals(p1))
	assert.False(t, p3.Equals(p2))
	assert.False(t, p3.Equals(p4))

	assert.False(t, p4.Equals(p1))
	assert.False(t, p4.Equals(p2))
	assert.False(t, p4.Equals(p3))
}

func TestPositionMoveToPositionValidStart(t *testing.T) {
	target := NewPosition(true, false, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "start", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionValidHome(t *testing.T) {
	target := NewPosition(false, true, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "home", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionValidSafe(t *testing.T) {
	safe := 3
	target := NewPosition(false, false, &safe, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "safe 3", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionValidSquare(t *testing.T) {
	square := 10
	target := NewPosition(false, false, nil, &square)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.Nil(t, err)
	assert.Equal(t, target, position)
	assert.NotSame(t, target, position)
	assert.Equal(t, "square 10", fmt.Sprintf("%s", position))
}

func TestPositionMoveToPositionInvalidMultiple(t *testing.T) {
	one := 1
	for _, target := range []Position {
		NewPosition(true, true, nil, nil),
		NewPosition(true, false, &one, nil),
		NewPosition(true, false, nil, &one),
		NewPosition(false, true, &one, nil),
		NewPosition(false, true, nil, &one),
		NewPosition(false, false, &one, &one),
	} {
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToPosition(target)
		assert.EqualError(t, err, "invalid position")
	}
}

func TestPositionMoveToPositionInvalidNone(t *testing.T) {
	target := NewPosition(false, false, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToPosition(target)
	assert.EqualError(t, err, "invalid position")
}

func TestPositionMoveToPositionInvalidSafe(t *testing.T) {
	for _, safe := range []int {-1000, -2, -1, 5, 6, 1000 } {
		target := NewPosition(false, false, &safe, nil)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToPosition(target)
		assert.EqualError(t, err, "invalid safe square")
	}
}

func TestPositionMoveToPositionInvalidSquare(t *testing.T) {
	for _, square := range []int { -1000, -2, -1, 60, 61, 1000 } {
		target := NewPosition(false, false, nil, &square)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToPosition(target)
		assert.EqualError(t, err, "invalid square")
	}
}

func TestPositionMoveToStart(t *testing.T) {
	expected := NewPosition(true, false, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToStart()
	assert.Nil(t, err)
	assert.Equal(t, expected, position)
}

func TestPositionMoveToHome(t *testing.T) {
	expected := NewPosition(false, true, nil, nil)
	position := NewPosition(false, false, nil, nil)
	err := position.MoveToHome()
	assert.Nil(t, err)
	assert.Equal(t, expected, position)
}

func TestPositionMoveToSafeValid(t *testing.T) {
	for safe := 0; safe < SafeSquares; safe++ {
		expected := NewPosition(false, false, &safe, nil)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSafe(safe)
		assert.Nil(t, err)
		assert.Equal(t, expected, position)
	}
}

func TestPositionMoveToSafeInvalid(t *testing.T) {
	for _, safe := range []int {-1000, -2, -1, 5, 6, 1000 } {
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSafe(safe)
		assert.EqualError(t, err, "invalid safe square")
	}
}

func TestPositionMoveToSquareValid(t *testing.T) {
	for square := 0; square < BoardSquares; square++ {
		expected := NewPosition(false, false, nil, &square)
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSquare(square)
		assert.Nil(t, err)
		assert.Equal(t, expected, position)
	}
}

func TestPositionMoveToSquareInvalid(t *testing.T) {
	for _, square := range []int { -1000, -2, -1, 60, 61, 1000 } {
		position := NewPosition(false, false, nil, nil)
		err := position.MoveToSquare(square)
		assert.EqualError(t, err, "invalid square")
	}
}

func TestNewPawn(t *testing.T) {
	obj := NewPawn(Red, 13)
	assert.Equal(t, Red, obj.Color())
	assert.Equal(t, 13, obj.Index())
	assert.Equal(t, "Red13", obj.Name())
	assert.Equal(t, NewPosition(true, false, nil, nil), obj.Position())
	assert.Equal(t, "Red13->start", fmt.Sprintf("%s", obj))
}

func TestPawnEquals(t *testing.T) {
	p1 := NewPawn(Red, 13)
	p2 := NewPawn(Blue, 4)

	// note: it is important to test with assert.True()/assert.False() and x.Equals(y)
	// because assert.Equals() and assert.NotEquals() are not aware of our equality by value concept

	assert.True(t, p1.Equals(p1))
	assert.True(t, p2.Equals(p2))

	assert.False(t, p1.Equals(nil))
	assert.False(t, p1.Equals(p2))

	assert.False(t, p2.Equals(nil))
	assert.False(t, p2.Equals(p1))
}

func TestPawnCopy(t *testing.T) {
	obj := NewPawn(Red, 13)
	copied := obj.Copy()
	assert.Equal(t, obj, copied)
	assert.NotSame(t, obj, copied)
}

func TestPawnSetPosition(t *testing.T) {
	obj := NewPawn(Red, 13)
	target := NewPosition(false, true, nil, nil)
	obj.SetPosition(target)
	assert.Equal(t, target, obj.Position())
	assert.Equal(t, "Red13->home", fmt.Sprintf("%s", obj))
}