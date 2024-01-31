package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistanceToHome(t *testing.T) {
	// distance from home is always 0
	for _, color := range []PlayerColor{ Red, Yellow, Green } {
		assert.Equal(t, 0, DistanceToHome(pawnHome(color)))
	}

	// distance from start is always 65
	for _, color := range []PlayerColor{ Red, Yellow, Green } {
		assert.Equal(t, 65, DistanceToHome(pawnStart(color)))
	}

	// distance from within safe is always <= 5
	assert.Equal(t, 5, DistanceToHome(pawnSafe(Red, 0)))
	assert.Equal(t, 4, DistanceToHome(pawnSafe(Red, 1)))
	assert.Equal(t, 3, DistanceToHome(pawnSafe(Red, 2)))
	assert.Equal(t, 2, DistanceToHome(pawnSafe(Red, 3)))
	assert.Equal(t, 1, DistanceToHome(pawnSafe(Red, 4)))

	// distance from circle is always 64
	assert.Equal(t, 64, DistanceToHome(pawnSquare(Red, 4)))
	assert.Equal(t, 64, DistanceToHome(pawnSquare(Blue, 19)))
	assert.Equal(t, 64, DistanceToHome(pawnSquare(Yellow, 34)))
	assert.Equal(t, 64, DistanceToHome(pawnSquare(Green, 49)))

	// distance from square between turn and circle is always 65
	assert.Equal(t, 65, DistanceToHome(pawnSquare(Red, 3)))
	assert.Equal(t, 65, DistanceToHome(pawnSquare(Blue, 18)))
	assert.Equal(t, 65, DistanceToHome(pawnSquare(Yellow, 33)))
	assert.Equal(t, 65, DistanceToHome(pawnSquare(Green, 48)))

	// distance from turn is always 6
	assert.Equal(t, 6, DistanceToHome(pawnSquare(Red, 2)))
	assert.Equal(t, 6, DistanceToHome(pawnSquare(Blue, 17)))
	assert.Equal(t, 6, DistanceToHome(pawnSquare(Yellow, 32)))
	assert.Equal(t, 6, DistanceToHome(pawnSquare(Green, 47)))

	// check some arbitrary squares
	assert.Equal(t, 7, DistanceToHome(pawnSquare(Red, 1)))
	assert.Equal(t, 8, DistanceToHome(pawnSquare(Red, 0)))
	assert.Equal(t, 9, DistanceToHome(pawnSquare(Red, 59)))
	assert.Equal(t, 59, DistanceToHome(pawnSquare(Red, 9)))
	assert.Equal(t, 23, DistanceToHome(pawnSquare(Blue, 0)))
	assert.Equal(t, 13, DistanceToHome(pawnSquare(Green, 40)))
}

func TestCalculatePositionHome(t *testing.T) {
	for _, color := range PlayerColors.Members() {
		position := NewPosition(false, false, nil, nil)
		_ = position.MoveToHome()
		_, err := calculatePosition(color, position, 1)
		assert.EqualError(t, err, "pawn in home or start may not move")
	}
}

func TestCalculatePositionStart(t *testing.T) {
	for _, color := range PlayerColors.Members() {
		position := NewPosition(false, false, nil, nil)
		_ = position.MoveToStart()
		_, err := calculatePosition(color, position, 1)
		assert.EqualError(t, err, "pawn in home or start may not move")
	}
}

func TestCalculatePositionFromSafe(t *testing.T) {
	var result Position
	var err error
	var color PlayerColor

	for _, color = range PlayerColors.Members() {
		result, err = calculatePosition(color, positionSafe(0), 0)
		assert.Nil(t, err)
		assert.Equal(t, positionSafe(0), result)

		result, err = calculatePosition(color, positionSafe(3), 0)
		assert.Nil(t, err)
		assert.Equal(t, positionSafe(3), result)
	}

	for _, color = range PlayerColors.Members() {
		result, err = calculatePosition(color, positionSafe(0), 1)
		assert.Nil(t, err)
		assert.Equal(t, positionSafe(1), result)

		result, err = calculatePosition(color, positionSafe(2), 2)
		assert.Nil(t, err)
		assert.Equal(t, positionSafe(4), result)

		result, err = calculatePosition(color, positionSafe(4), 1)
		assert.Nil(t, err)
		assert.Equal(t, positionHome(), result)
	}

	for _, color = range PlayerColors.Members() {
		result, err = calculatePosition(color, positionSafe(3), 3)
		assert.EqualError(t, err, "pawn cannot move past home")

		result, err = calculatePosition(color, positionSafe(4), 2)
		assert.EqualError(t, err, "pawn cannot move past home")
	}

	for _, color = range PlayerColors.Members() {
		result, err = calculatePosition(color, positionSafe(4), -2)
		assert.Nil(t, err)
		assert.Equal(t, positionSafe(2), result)

		result, err = calculatePosition(color, positionSafe(1), -1)
		assert.Nil(t, err)
		assert.Equal(t, positionSafe(0), result)
	}

	result, err = calculatePosition(Red, positionSafe(0), -1)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(2), result)

	result, err = calculatePosition(Red, positionSafe(0), -2)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(1), result)

	result, err = calculatePosition(Red, positionSafe(0), -3)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(0), result)

	result, err = calculatePosition(Red, positionSafe(0), -4)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(59), result)

	result, err = calculatePosition(Red, positionSafe(0), -5)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(58), result)

	result, err = calculatePosition(Blue, positionSafe(0), -1)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(17), result)

	result, err = calculatePosition(Blue, positionSafe(0), -2)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(16), result)

	result, err = calculatePosition(Yellow, positionSafe(0), -1)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(32), result)

	result, err = calculatePosition(Yellow, positionSafe(0), -2)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(31), result)

	result, err = calculatePosition(Green, positionSafe(0), -1)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(47), result)

	result, err = calculatePosition(Green, positionSafe(0), -2)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(46), result)
}

func TestCalculatePositionFromSquare(t *testing.T) {
	var result Position
	var err error
	var color PlayerColor

	result, err = calculatePosition(Red, positionSquare(58), 1)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(59), result)

	result, err = calculatePosition(Red, positionSquare(59), 1)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(0), result)

	result, err = calculatePosition(Red, positionSquare(54), 5)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(59), result)

	result, err = calculatePosition(Red, positionSquare(54), 6)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(0), result)

	result, err = calculatePosition(Red, positionSquare(54), 7)
	assert.Nil(t, err)
	assert.Equal(t, positionSquare(1), result)

	for _, color = range PlayerColors.Members() {
		result, err = calculatePosition(color, positionSquare(54), 5)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(59), result)

		result, err = calculatePosition(color, positionSquare(54), 6)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(0), result)

		result, err = calculatePosition(color, positionSquare(54), 7)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(1), result)

		result, err = calculatePosition(color, positionSquare(58), 1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(59), result)

		result, err = calculatePosition(color, positionSquare(59), 1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(0), result)

		result, err = calculatePosition(color, positionSquare(0), 1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(1), result)

		result, err = calculatePosition(color, positionSquare(1), 1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(2), result)

		result, err = calculatePosition(color, positionSquare(10), 5)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(15), result)
	}

	for _, color = range PlayerColors.Members() {
		result, err = calculatePosition(color, positionSquare(59), -5)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(54), result)

		result, err = calculatePosition(color, positionSquare(0), -6)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(54), result)

		result, err = calculatePosition(color, positionSquare(1), -7)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(54), result)

		result, err = calculatePosition(color, positionSquare(59), -1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(58), result)

		result, err = calculatePosition(color, positionSquare(0), -1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(59), result)

		result, err = calculatePosition(color, positionSquare(1), -1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(0), result)

		result, err = calculatePosition(color, positionSquare(2), -1)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(1), result)

		result, err = calculatePosition(color, positionSquare(15), -5)
		assert.Nil(t, err)
		assert.Equal(t, positionSquare(10), result)
	}

	result, err = calculatePosition(Red, positionSquare(0), 3)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Red, positionSquare(1), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Red, positionSquare(2), 1)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Red, positionSquare(1), 3)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Red, positionSquare(2), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Red, positionSquare(2), 6)
	assert.Nil(t, err)
	assert.Equal(t, positionHome(), result)

	result, err = calculatePosition(Red, positionSquare(51), 12)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Red, positionSquare(52), 12)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Red, positionSquare(58), 5)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Red, positionSquare(59), 4)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Red, positionSquare(2), 7)
	assert.EqualError(t, err, "pawn cannot move past home")

	result, err = calculatePosition(Blue, positionSquare(16), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Blue, positionSquare(17), 1)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Blue, positionSquare(16), 3)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Blue, positionSquare(17), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Blue, positionSquare(17), 6)
	assert.Nil(t, err)
	assert.Equal(t, positionHome(), result)

	result, err = calculatePosition(Blue, positionSquare(17), 7)
	assert.EqualError(t, err, "pawn cannot move past home")

	result, err = calculatePosition(Yellow, positionSquare(31), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Yellow, positionSquare(32), 1)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Yellow, positionSquare(31), 3)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Yellow, positionSquare(32), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Yellow, positionSquare(32), 6)
	assert.Nil(t, err)
	assert.Equal(t, positionHome(), result)

	result, err = calculatePosition(Yellow, positionSquare(32), 7)
	assert.EqualError(t, err, "pawn cannot move past home")

	result, err = calculatePosition(Green, positionSquare(46), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Green, positionSquare(47), 1)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(0), result)

	result, err = calculatePosition(Green, positionSquare(46), 3)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Green, positionSquare(47), 2)
	assert.Nil(t, err)
	assert.Equal(t, positionSafe(1), result)

	result, err = calculatePosition(Green, positionSquare(47), 6)
	assert.Nil(t, err)
	assert.Equal(t, positionHome(), result)

	result, err = calculatePosition(Green, positionSquare(47), 7)
	assert.EqualError(t, err, "pawn cannot move past home")
}

func positionHome() Position {
	return NewPosition(false, true, nil, nil)
}

func positionStart() Position {
	return NewPosition(true, false, nil, nil)
}

func positionSafe(safe int) Position {
	return NewPosition(false, false, &safe, nil)
}

func positionSquare(square int) Position {
	return NewPosition(false, false, nil, &square)
}

func pawnHome(color PlayerColor) Pawn {
	pawn := NewPawn(color, 0)
	pawn.SetPosition(positionHome())
	return pawn
}

func pawnStart(color PlayerColor) Pawn {
	pawn := NewPawn(color, 0)
	pawn.SetPosition(positionStart())
	return pawn
}

func pawnSafe(color PlayerColor, safe int) Pawn {
	pawn := NewPawn(color, 0)
	pawn.SetPosition(positionSafe(safe))
	return pawn
}

func pawnSquare(color PlayerColor, square int) Pawn {
	pawn := NewPawn(color, 0)
	pawn.SetPosition(positionSquare(square))
	return pawn
}