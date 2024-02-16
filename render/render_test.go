package render

import (
	"fmt"
	"os"
	"testing"

	"github.com/pronovic/go-apologies/model"
	"github.com/stretchr/testify/assert"
)

func TestEmpty2Player(t *testing.T) {
	game := empty(2)
	executeTest(t, game, "empty2")
}

func TestEmpty3Player(t *testing.T) {
	game := empty(3)
	executeTest(t, game, "empty3")
}

func TestEmpty4Player(t *testing.T) {
	game := empty(4)
	executeTest(t, game, "empty4")
}

func TestHome(t *testing.T) {
	game := fillHome()
	executeTest(t, game, "home")
}

func TestSafe03(t *testing.T) {
	game := fillSafe(0)
	executeTest(t, game, "safe_03")
}

func TestSafe14(t *testing.T) {
	game := fillSafe(1)
	executeTest(t, game, "safe_14")
}

func TestTop(t *testing.T) {
	game := fillSquares(0, 15)
	executeTest(t, game, "top")
}

func TestRight(t *testing.T) {
	game := fillSquares(16, 29)
	executeTest(t, game, "right")
}

func TestBottom(t *testing.T) {
	game := fillSquares(30, 45)
	executeTest(t, game, "bottom")
}

func TestLeft(t *testing.T) {
	game := fillSquares(46, 59)
	executeTest(t, game, "left")
}

func executeTest(t *testing.T, game model.Game, testdata string) {
	raw, err := os.ReadFile(fmt.Sprintf("../testdata/render/%s", testdata))
	assert.NoError(t, err)
	expected := string(raw)
	actual, err := Board(game)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

// empty Create an empty game with some number of players
func empty(players int) model.Game {
	game, _ := model.NewGame(players, nil)
	return game
}

// fillHome Create a game with all players in home
func fillHome() model.Game {
	game := empty(4)

	for _, color := range []model.PlayerColor{model.Blue, model.Red, model.Yellow, model.Green} {
		for pawn := 0; pawn < 4; pawn++ {
			_ = game.Players()[color].Pawns()[pawn].Position().MoveToHome()
		}
	}

	return game
}

// fillSafe Create a game with all players in the safe zone
func fillSafe(start int) model.Game {
	game := empty(4)

	for _, color := range []model.PlayerColor{model.Blue, model.Red, model.Yellow, model.Green} {
		for pawn := 0; pawn < 4; pawn++ {
			_ = game.Players()[color].Pawns()[pawn].Position().MoveToSafe(pawn + start)
		}
	}

	return game
}

// fillSquares Fill a range of squares on the board with pieces from various players
func fillSquares(start int, end int) model.Game {
	game := empty(4)

	square := 0
	for pawn := 0; pawn < 4; pawn++ {
		for _, color := range []model.PlayerColor{model.Blue, model.Red, model.Yellow, model.Green} {
			if square+start <= end {
				_ = game.Players()[color].Pawns()[pawn].Position().MoveToSquare(square + start)
				square += 1
			}
		}
	}

	return game
}
