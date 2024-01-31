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

func pawnHome(color PlayerColor) Pawn {
	position := NewPosition(false, true, nil, nil)
	pawn := NewPawn(color, 0)
	pawn.SetPosition(position)
	return pawn
}

func pawnStart(color PlayerColor) Pawn {
	position := NewPosition(true, false, nil, nil)
	pawn := NewPawn(color, 0)
	pawn.SetPosition(position)
	return pawn
}

func pawnSafe(color PlayerColor, safe int) Pawn {
	position := NewPosition(false, false, &safe, nil)
	pawn := NewPawn(color, 0)
	pawn.SetPosition(position)
	return pawn
}

func pawnSquare(color PlayerColor, square int) Pawn {
	position := NewPosition(false, false, nil, &square)
	pawn := NewPawn(color, 0)
	pawn.SetPosition(position)
	return pawn
}