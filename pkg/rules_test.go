package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAction(t *testing.T) {
	pawn := NewPawn(Red, 0)
	position := NewPosition(false, false, nil, nil)
	obj := NewAction(MoveToPosition, pawn, position)
	assert.Equal(t, MoveToPosition, obj.Type())
	assert.Same(t, pawn, obj.Pawn())
	assert.Same(t, position, obj.Position())
}

func TestActionSetPosition(t *testing.T) {
	pawn := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	position2 := NewPosition(true, false, nil, nil)
	obj := NewAction(MoveToPosition, pawn, position1)
	obj.SetPosition(position2)
	assert.Same(t, position2, obj.Position())
}

func TestActionEquals(t *testing.T) {
	pawn1 := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	obj1 := NewAction(MoveToPosition, pawn1, position1)

	pawn2 := NewPawn(Red, 0)
	position2 := NewPosition(false, false, nil, nil)
	obj2 := NewAction(MoveToStart, pawn2, position2)

	position3 := NewPosition(false, false, nil, nil)
	obj3 := NewAction(MoveToStart, nil, position3)

	pawn4 := NewPawn(Red, 0)
	obj4 := NewAction(MoveToStart, pawn4, nil)

	obj5 := NewAction(MoveToPosition, nil, nil)

	assert.Equal(t, obj1, obj1)
	assert.Equal(t, obj2, obj2)
	assert.Equal(t, obj3, obj3)
	assert.Equal(t, obj4, obj4)
	assert.Equal(t, obj5, obj5)

	assert.NotEqual(t, obj1, obj2)
	assert.NotEqual(t, obj1, obj3)
	assert.NotEqual(t, obj1, obj4)
	assert.NotEqual(t, obj1, obj5)

	assert.NotEqual(t, obj5, obj1)
	assert.NotEqual(t, obj5, obj2)
	assert.NotEqual(t, obj5, obj3)
	assert.NotEqual(t, obj5, obj4)
}

func TestNewMove(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 1, 2)
	sideEffects := make([]Action, 2, 3)
	obj := NewMove(card, actions, sideEffects)
	assert.NotEmpty(t, obj.Id()) // filled in with a UUID
	assert.Equal(t, card, obj.Card())
	assert.Equal(t, actions, obj.Actions())
	assert.Equal(t, sideEffects, obj.SideEffects())

	obj2 := NewMove(card, actions, sideEffects)
	assert.NotEqual(t, obj.Id(), obj2.Id()) // just make sure we get a unique UUID each time
}

func TestMoveAddSideEffect(t *testing.T) {
	card := NewCard("1", Card1)
	actions := make([]Action, 0)
	sideEffects := make([]Action, 0)
	obj := NewMove(card, actions, sideEffects)

	pawn := NewPawn(Red, 0)
	position := NewPosition(false, false, nil, nil)
	sideEffect := NewAction(MoveToPosition, pawn, position)
	obj.AddSideEffect(sideEffect)
	assert.Equal(t, []Action {sideEffect}, obj.SideEffects())
}

func TestMoveMergedActions(t *testing.T) {
	pawn1 := NewPawn(Red, 0)
	position1 := NewPosition(false, false, nil, nil)
	action1 := NewAction(MoveToPosition, pawn1, position1)

	pawn2 := NewPawn(Red, 0)
	position2 := NewPosition(false, false, nil, nil)
	action2 := NewAction(MoveToStart, pawn2, position2)

	position3 := NewPosition(false, false, nil, nil)
	action3 := NewAction(MoveToStart, nil, position3)

	pawn4 := NewPawn(Red, 0)
	action4 := NewAction(MoveToStart, pawn4, nil)

	card := NewCard("1", Card1)
	actions := []Action { action1, action2 }
	sideEffects := []Action { action3, action4 }
	expected := []Action { action1, action2, action3, action4 }
	obj := NewMove(card, actions, sideEffects)
	assert.Equal(t, expected, obj.MergedActions())
}

func TestStartGameStandardMode(t *testing.T) {
	game, _ := NewGame(2)
	err := StartGame(game, StandardMode)

	assert.Nil(t, err)
	assert.True(t, game.Started())

	assert.Equal(t, Red, game.Players()[Red].Color())
	assert.Equal(t, 0, len(game.Players()[Red].Hand()))

	assert.Equal(t, Yellow, game.Players()[Yellow].Color())
	assert.Equal(t, 0, len(game.Players()[Yellow].Hand()))

	err = StartGame(game, StandardMode)
	assert.EqualError(t, err, "game is already started")
}

func TestStartGameAdultMode(t *testing.T) {
	game, _ := NewGame(4)
	err := StartGame(game, AdultMode)

	assert.Equal(t, Red, game.Players()[Red].Color())
	assert.Equal(t, AdultHand, len(game.Players()[Red].Hand()))
	assert.Equal(t, 4, *game.Players()[Red].Pawns()[0].Position().Square())

	assert.Equal(t, Yellow, game.Players()[Yellow].Color())
	assert.Equal(t, AdultHand, len(game.Players()[Yellow].Hand()))
	assert.Equal(t, 34, *game.Players()[Yellow].Pawns()[0].Position().Square())

	assert.Equal(t, Green, game.Players()[Green].Color())
	assert.Equal(t, AdultHand, len(game.Players()[Green].Hand()))
	assert.Equal(t, 49, *game.Players()[Green].Pawns()[0].Position().Square())

	assert.Equal(t, Blue, game.Players()[Blue].Color())
	assert.Equal(t, AdultHand, len(game.Players()[Blue].Hand()))
	assert.Equal(t, 19, *game.Players()[Blue].Pawns()[0].Position().Square())

	err = StartGame(game, AdultMode)
	assert.EqualError(t, err, "game is already started")
}

func TestExecuteMove(t *testing.T) {
	actions := []Action {
		NewAction(MoveToPosition, NewPawn(Red, 1), positionSquare(10)),
		NewAction(MoveToPosition, NewPawn(Yellow, 3), positionSquare(11)),
	}

	sideEffects := []Action {
		NewAction(MoveToStart, NewPawn(Blue, 2), nil),
		NewAction(MoveToPosition, NewPawn(Green, 0), positionSquare(12)),
	}

	move := NewMove(NewCard("1", Card1), actions, sideEffects)

	game, _ := NewGame(4)
	player := game.Players()[Red]

	err := ExecuteMove(game, player, move)
	assert.Nil(t, err)

	assert.Equal(t, 10, *game.Players()[Red].Pawns()[1].Position().Square())
	assert.Equal(t, 11, *game.Players()[Yellow].Pawns()[3].Position().Square())
	assert.True(t, game.Players()[Blue].Pawns()[2].Position().Start())
	assert.Equal(t, 12, *game.Players()[Green].Pawns()[0].Position().Square())
}

func TestEvaluateMove(t *testing.T) {
	var err error
	var result PlayerView

	actions := []Action {
		NewAction(MoveToPosition, NewPawn(Red, 1), positionSquare(10)),
		NewAction(MoveToPosition, NewPawn(Yellow, 3), positionSquare(11)),
	}

	sideEffects := []Action {
		NewAction(MoveToStart, NewPawn(Blue, 2), nil),
		NewAction(MoveToPosition, NewPawn(Green, 0), positionSquare(12)),
	}

	move := NewMove(NewCard("1", Card1), actions, sideEffects)

	game, _ := NewGame(4)
	view, err := game.CreatePlayerView(Red)
	assert.Nil(t, err)

	expected := view.Copy()

	err = expected.Player().Pawns()[1].Position().MoveToSquare(10)
	assert.Nil(t, err)

	err = expected.Opponents()[Yellow].Pawns()[3].Position().MoveToSquare(11)
	assert.Nil(t, err)

	err = expected.Opponents()[Yellow].Pawns()[2].Position().MoveToStart()
	assert.Nil(t, err)

	err = expected.Opponents()[Green].Pawns()[0].Position().MoveToSquare(12)
	assert.Nil(t, err)

	result, err = EvaluateMove(view, move)
	assert.Equal(t, expected, result)
}

func TestDistanceToHome(t *testing.T) {
	// distance from home is always 0
	for _, color := range []PlayerColor{ Red, Yellow, Green } {
		assert.Equal(t, 0, distanceToHome(pawnHome(color)))
	}

	// distance from start is always 65
	for _, color := range []PlayerColor{ Red, Yellow, Green } {
		assert.Equal(t, 65, distanceToHome(pawnStart(color)))
	}

	// distance from within safe is always <= 5
	assert.Equal(t, 5, distanceToHome(pawnSafe(Red, 0)))
	assert.Equal(t, 4, distanceToHome(pawnSafe(Red, 1)))
	assert.Equal(t, 3, distanceToHome(pawnSafe(Red, 2)))
	assert.Equal(t, 2, distanceToHome(pawnSafe(Red, 3)))
	assert.Equal(t, 1, distanceToHome(pawnSafe(Red, 4)))

	// distance from circle is always 64
	assert.Equal(t, 64, distanceToHome(pawnSquare(Red, 4)))
	assert.Equal(t, 64, distanceToHome(pawnSquare(Blue, 19)))
	assert.Equal(t, 64, distanceToHome(pawnSquare(Yellow, 34)))
	assert.Equal(t, 64, distanceToHome(pawnSquare(Green, 49)))

	// distance from square between turn and circle is always 65
	assert.Equal(t, 65, distanceToHome(pawnSquare(Red, 3)))
	assert.Equal(t, 65, distanceToHome(pawnSquare(Blue, 18)))
	assert.Equal(t, 65, distanceToHome(pawnSquare(Yellow, 33)))
	assert.Equal(t, 65, distanceToHome(pawnSquare(Green, 48)))

	// distance from turn is always 6
	assert.Equal(t, 6, distanceToHome(pawnSquare(Red, 2)))
	assert.Equal(t, 6, distanceToHome(pawnSquare(Blue, 17)))
	assert.Equal(t, 6, distanceToHome(pawnSquare(Yellow, 32)))
	assert.Equal(t, 6, distanceToHome(pawnSquare(Green, 47)))

	// check some arbitrary squares
	assert.Equal(t, 7, distanceToHome(pawnSquare(Red, 1)))
	assert.Equal(t, 8, distanceToHome(pawnSquare(Red, 0)))
	assert.Equal(t, 9, distanceToHome(pawnSquare(Red, 59)))
	assert.Equal(t, 59, distanceToHome(pawnSquare(Red, 9)))
	assert.Equal(t, 23, distanceToHome(pawnSquare(Blue, 0)))
	assert.Equal(t, 13, distanceToHome(pawnSquare(Green, 40)))
}

func TestCalculatePositionHome(t *testing.T) {
	for _, color := range PlayerColors.Members() {
		calculatePositionFailure(t, color, positionHome(), 1, "pawn in home or start may not move")
	}
}

func TestCalculatePositionStart(t *testing.T) {
	for _, color := range PlayerColors.Members() {
		calculatePositionFailure(t, color, positionStart(), 1, "pawn in home or start may not move")
	}
}

func TestCalculatePositionFromSafe(t *testing.T) {
	var color PlayerColor

	for _, color = range PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSafe(0), 0, positionSafe(0))
		calculatePositionSuccess(t, color, positionSafe(3), 0, positionSafe(3))
	}

	for _, color = range PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSafe(0), 1, positionSafe(1))
		calculatePositionSuccess(t, color, positionSafe(2), 2, positionSafe(4))
		calculatePositionSuccess(t, color, positionSafe(4), 1, positionHome())
	}

	for _, color = range PlayerColors.Members() {
		calculatePositionFailure(t, color, positionSafe(3), 3, "pawn cannot move past home")
		calculatePositionFailure(t, color, positionSafe(4), 2, "pawn cannot move past home")
	}

	for _, color = range PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSafe(4), -2, positionSafe(2))
		calculatePositionSuccess(t, color, positionSafe(1), -1, positionSafe(0))
	}

	calculatePositionSuccess(t, Red, positionSafe(0), -1, positionSquare(2))
	calculatePositionSuccess(t, Red, positionSafe(0), -2, positionSquare(1))
	calculatePositionSuccess(t, Red, positionSafe(0), -3, positionSquare(0))
	calculatePositionSuccess(t, Red, positionSafe(0), -4, positionSquare(59))
	calculatePositionSuccess(t, Red, positionSafe(0), -5, positionSquare(58))

	calculatePositionSuccess(t, Blue, positionSafe(0), -1, positionSquare(17))
	calculatePositionSuccess(t, Blue, positionSafe(0), -2, positionSquare(16))

	calculatePositionSuccess(t, Yellow, positionSafe(0), -1, positionSquare(32))
	calculatePositionSuccess(t, Yellow, positionSafe(0), -2, positionSquare(31))

	calculatePositionSuccess(t, Green, positionSafe(0), -1, positionSquare(47))
	calculatePositionSuccess(t, Green, positionSafe(0), -2, positionSquare(46))
}

func TestCalculatePositionFromSquare(t *testing.T) {
	var color PlayerColor

	calculatePositionSuccess(t, Red, positionSquare(58), 1, positionSquare(59))
	calculatePositionSuccess(t, Red, positionSquare(59), 1, positionSquare(0))
	calculatePositionSuccess(t, Red, positionSquare(54), 5, positionSquare(59))
	calculatePositionSuccess(t, Red, positionSquare(54), 6, positionSquare(0))
	calculatePositionSuccess(t, Red, positionSquare(54), 7, positionSquare(1))

	for _, color = range PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSquare(54), 5, positionSquare(59))
		calculatePositionSuccess(t, color, positionSquare(54), 6, positionSquare(0))
		calculatePositionSuccess(t, color, positionSquare(54), 7, positionSquare(1))
		calculatePositionSuccess(t, color, positionSquare(58), 1, positionSquare(59))
		calculatePositionSuccess(t, color, positionSquare(59), 1, positionSquare(0))
		calculatePositionSuccess(t, color, positionSquare(0), 1, positionSquare(1))
		calculatePositionSuccess(t, color, positionSquare(1), 1, positionSquare(2))
		calculatePositionSuccess(t, color, positionSquare(10), 5, positionSquare(15))
	}

	for _, color = range PlayerColors.Members() {
		calculatePositionSuccess(t, color, positionSquare(59), -5, positionSquare(54))
		calculatePositionSuccess(t, color, positionSquare(0), -6, positionSquare(54))
		calculatePositionSuccess(t, color, positionSquare(1), -7, positionSquare(54))
		calculatePositionSuccess(t, color, positionSquare(59), -1, positionSquare(58))
		calculatePositionSuccess(t, color, positionSquare(0), -1, positionSquare(59))
		calculatePositionSuccess(t, color, positionSquare(1), -1, positionSquare(0))
		calculatePositionSuccess(t, color, positionSquare(2), -1, positionSquare(1))
		calculatePositionSuccess(t, color, positionSquare(15), -5, positionSquare(10))
	}

	calculatePositionSuccess(t, Red, positionSquare(0), 3, positionSafe(0))
	calculatePositionSuccess(t, Red, positionSquare(1), 2, positionSafe(0))
	calculatePositionSuccess(t, Red, positionSquare(2), 1, positionSafe(0))
	calculatePositionSuccess(t, Red, positionSquare(1), 3, positionSafe(1))
	calculatePositionSuccess(t, Red, positionSquare(2), 2, positionSafe(1))
	calculatePositionSuccess(t, Red, positionSquare(2), 6, positionHome())
	calculatePositionSuccess(t, Red, positionSquare(51), 12, positionSafe(0))
	calculatePositionSuccess(t, Red, positionSquare(52), 12, positionSafe(1))
	calculatePositionSuccess(t, Red, positionSquare(58), 5, positionSafe(0))
	calculatePositionSuccess(t, Red, positionSquare(59), 4, positionSafe(0))
	calculatePositionFailure(t, Red, positionSquare(2), 7, "pawn cannot move past home")

	calculatePositionSuccess(t, Blue, positionSquare(16), 2, positionSafe(0))
	calculatePositionSuccess(t, Blue, positionSquare(17), 1, positionSafe(0))
	calculatePositionSuccess(t, Blue, positionSquare(16), 3, positionSafe(1))
	calculatePositionSuccess(t, Blue, positionSquare(17), 2, positionSafe(1))
	calculatePositionSuccess(t, Blue, positionSquare(17), 6, positionHome())
	calculatePositionFailure(t, Blue, positionSquare(17), 7, "pawn cannot move past home")

	calculatePositionSuccess(t, Yellow, positionSquare(31), 2, positionSafe(0))
	calculatePositionSuccess(t, Yellow, positionSquare(32), 1, positionSafe(0))
	calculatePositionSuccess(t, Yellow, positionSquare(31), 3, positionSafe(1))
	calculatePositionSuccess(t, Yellow, positionSquare(32), 2, positionSafe(1))
	calculatePositionSuccess(t, Yellow, positionSquare(32), 6, positionHome())
	calculatePositionFailure(t, Yellow, positionSquare(32), 7, "pawn cannot move past home")

	calculatePositionSuccess(t, Green, positionSquare(46), 2, positionSafe(0))
	calculatePositionSuccess(t, Green, positionSquare(47), 1, positionSafe(0))
	calculatePositionSuccess(t, Green, positionSquare(46), 3, positionSafe(1))
	calculatePositionSuccess(t, Green, positionSquare(47), 2, positionSafe(1))
	calculatePositionSuccess(t, Green, positionSquare(47), 6, positionHome())
	calculatePositionFailure(t, Green, positionSquare(47), 7, "pawn cannot move past home")
}

func TestConstructLegalMovesNoMovesWithCard(t *testing.T) {

}

func TestConstructLegalMovesNoMovesNoCard(t *testing.T) {

}

func TestConstructLegalMovesWithMovesWithCard(t *testing.T) {

}

func TestConstructLegalMovesWithMovesNoCard(t *testing.T) {

}

func TestConstructLegalMovesCard1(t *testing.T) {

}

func TestConstructLegalMovesCard2(t *testing.T) {

}

func TestConstructLegalMovesCard3(t *testing.T) {

}

func TestConstructLegalMovesCard4(t *testing.T) {

}

func TestConstructLegalMovesCard5(t *testing.T) {

}

func TestConstructLegalMovesCard7(t *testing.T) {

}

func TestConstructLegalMovesCard8(t *testing.T) {

}

func TestConstructLegalMovesCard10(t *testing.T) {

}

func TestConstructLegalMovesCard11(t *testing.T) {

}

func TestConstructLegalMovesCard12(t *testing.T) {

}

func TestConstructLegalMovesCardApologies(t *testing.T) {

}

func TestConstructLegalMovesCardSpecial(t *testing.T) {

}

func calculatePositionSuccess(t *testing.T, color PlayerColor, start Position, squares int, expected Position) {
	result, err := calculatePosition(color, start, squares)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func calculatePositionFailure(t *testing.T, color PlayerColor, start Position, squares int, expected string) {
	_, err := calculatePosition(color, start, squares)
	assert.EqualError(t, err, expected)
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
