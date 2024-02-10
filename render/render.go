package render

import (
	"errors"

	"github.com/pronovic/go-apologies/model"
)

// Index into boardText where a pawn can be placed into a specific square
// Squares are numbered from starting from the upper left, in a clockwise direction
// Only a single pawn can occupy a square
var squareIndexes = []int{
	177, 182, 187, 192, 197, 202, 207, 212, 217, 222, 227, 232, 237, 242, 247, 252,
	507, 765, 1023, 1281, 1539, 1797, 2055, 2313, 2571, 2829, 3087, 3345, 3603, 3861,
	4119, 4114, 4109, 4104, 4099, 4094, 4089, 4084, 4079, 4074, 4069, 4064, 4059, 4054, 4049, 4044,
	3786, 3528, 3270, 3012, 2754, 2496, 2238, 1980, 1722, 1464, 1206, 948, 690, 432,
}

// Player names as displayed on the board
var playerNames = map[model.PlayerColor]rune{
	model.Red:    'r',
	model.Blue:   'b',
	model.Yellow: 'y',
	model.Green:  'g',
}

// Indexes in boardText where a pawn can be placed into a start location, for each player
// There are 4 arbitrary spaces for each of the 4 available pawn
var startIndexes = map[model.PlayerColor][]int{
	model.Red:    {623, 625, 627, 629},
	model.Blue:   {1183, 1185, 1187, 1189},
	model.Yellow: {3579, 3581, 3583, 3585},
	model.Green:  {3019, 3021, 3023, 3025},
}

// Indexes in boardText where a pawn can be placed into a safe location, for each player
// There are 5 safe squares per color; only a single pawn can occupy a safe square
var safeIndexes = map[model.PlayerColor][]int{
	model.Red:    {442, 700, 958, 1216, 1474},
	model.Blue:   {760, 755, 750, 745, 740},
	model.Yellow: {3851, 3593, 3335, 3077, 2819},
	model.Green:  {3533, 3538, 3543, 3548, 3553},
}

// Indexes in boardText where a pawn can be placed into a home location, for each player
// There are 4 arbitrary home spaces for the 4 available pawns per player
var homeIndexes = map[model.PlayerColor][]int{
	model.Red:    {1817, 1819, 1821, 1823},
	model.Blue:   {643, 645, 647, 649},
	model.Yellow: {2388, 2390, 2392, 2394},
	model.Green:  {3559, 3561, 3563, 3565},
}

// Hardcoded representation of the board, built by hand
const boardText = `

      0    1    2    3    4    5    6    7    8    9    10   11   12   13   14   15
    ┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐
    │   │| ▶ || ◼ || ◼ || ● ||   ||   ||   ||   || ▶ || ◼ || ◼ || ◼ || ● ||   ||   |
    └───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘
    ┌───┐     ┌───┐                                                            ┌───┐
 59 │   │   0 │   │  ┌───────────┐       ┌───────────┐                         | ▼ | 16
    └───┘     └───┘  │ S T A R T │       │  H O M E  │  4    3    2    1    0  └───┘
    ┌───┐     ┌───┐  │           │       │           │┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐
 58 │ ● │   1 │   │  │  - - - -  │       │  - - - -  │|   ||   ||   ||   ||   || ◼ | 17
    └───┘     └───┘  │  0 1 2 3  │       │  0 1 2 3  │└───┘└───┘└───┘└───┘└───┘└───┘
    ┌───┐     ┌───┐  └───────────┘       └───────────┘                         ┌───┐
 57 │ ◼ │   2 │   │                                              ┌───────────┐ | ◼ | 18
    └───┘     └───┘                                              │ S T A R T │ └───┘
    ┌───┐     ┌───┐                                              │           │ ┌───┐
 56 │ ◼ │   3 │   │                                              │  - - - -  │ | ● | 19
    └───┘     └───┘                                              │  0 1 2 3  │ └───┘
    ┌───┐     ┌───┐                                              └───────────┘ ┌───┐
 55 │ ◼ │   4 │   │                                                            |   | 20
    └───┘     └───┘                                                            └───┘
    ┌───┐ ┌───────────┐                                                        ┌───┐
 54 │ ▲ │ │  H O M E  │                                                        |   | 21
    └───┘ │           │                                                        └───┘
    ┌───┐ │  - - - -  │                                                        ┌───┐
 53 │   │ │  0 1 2 3  │                                                        |   | 22
    └───┘ └───────────┘                                                        └───┘
    ┌───┐                                                        ┌───────────┐ ┌───┐
 52 │   │                                                        │  H O M E  │ |   | 23
    └───┘                                                        │           │ └───┘
    ┌───┐                                                        │  - - - -  │ ┌───┐
 51 │   │                                                        │  0 1 2 3  │ | ▼ | 24
    └───┘                                                        └───────────┘ └───┘
    ┌───┐                                                            ┌───┐     ┌───┐
 50 │   │                                                            │   │ 4   | ◼ | 25
    └───┘ ┌───────────┐                                              └───┘     └───┘
    ┌───┐ │ S T A R T │                                              ┌───┐     ┌───┐
 49 │ ● │ │           │                                              │   │ 3   | ◼ | 26
    └───┘ │  - - - -  │                                              └───┘     └───┘
    ┌───┐ │  0 1 2 3  │                                              ┌───┐     ┌───┐
 48 │ ◼ │ └───────────┘                                              │   │ 2   | ◼ | 27
    └───┘                         ┌───────────┐       ┌───────────┐  └───┘     └───┘
    ┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐│  H O M E  │       │ S T A R T │  ┌───┐     ┌───┐
 47 │ ◼ │|   ||   ||   ||   ||   |│           │       │           │  │   │ 1   | ● | 28
    └───┘└───┘└───┘└───┘└───┘└───┘│  - - - -  │       │  - - - -  │  └───┘     └───┘
    ┌───┐  0    1    2    3    4  │  0 1 2 3  │       │  0 1 2 3  │  ┌───┐     ┌───┐
 46 │ ▲ │                         └───────────┘       └───────────┘  │   │ 0   |   | 29
    └───┘                                                            └───┘     └───┘
    ┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐┌───┐
    │   │|   || ● || ◼ || ◼ || ◼ || ◀ ||   ||   ||   ||   || ● || ◼ || ◼ || ◀ ||   | 
    └───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘└───┘
      45   44   43   42   41   40   39   38   37   36   35   34   33   32   31   30
`

func Board(game model.Game) (string, error) {
	board := []rune(boardText)

	for _, player := range game.Players() {
		for _, pawn := range player.Pawns() {
			if pawn.Position().Start() {
				index := startIndexes[pawn.Color()][pawn.Index()]
				board[index] = playerNames[pawn.Color()]
			} else if pawn.Position().Home() {
				index := homeIndexes[pawn.Color()][pawn.Index()]
				board[index] = playerNames[pawn.Color()]
			} else if pawn.Position().Safe() != nil {
				index := safeIndexes[pawn.Color()][*pawn.Position().Safe()]
				board[index] = playerNames[pawn.Color()]
			} else if pawn.Position().Square() != nil {
				index := squareIndexes[*pawn.Position().Square()]
				board[index] = playerNames[pawn.Color()]
			} else {
				return "", errors.New("pawn is not in a valid state")
			}
		}
	}

	return string(board), nil
}
