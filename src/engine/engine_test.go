package engine

import (
	"log"
	"slices"
	"testing"
	"time"
)

func TestCreateStopGame(t *testing.T) {

	_, gameInput, gameOutput := NewGame()

	game := <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	if game.ScoreLineCount != 0 {
		t.Errorf("Game line count not expected.  %d", game.ScoreLineCount)
	}

	if game.GameRows != 18 {
		t.Errorf("Game rows count not expected.  %d", game.GameRows)
	}

	if game.GameColumns != 10 {
		t.Errorf("Game columns count not expected.  %d", game.GameColumns)
	}

	if game.NumberPossiblePieces != 7 {
		t.Errorf("Game possible pieces not expected.  %d", game.NumberPossiblePieces)
	}

	gameInput <- PlayInputStop
	game = <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateGameOver {
		t.Errorf("Game not over.  %d", game.State)
	}
}

func TestCreateBucketGame(t *testing.T) {

	_, gameInput, gameOutput := NewBucketGame()

	game := <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	if game.ScoreLineCount != 0 {
		t.Errorf("Game line count not expected.  %d", game.ScoreLineCount)
	}

	if game.GameRows != 10 {
		t.Errorf("Game rows count not expected.  %d", game.GameRows)
	}

	if game.GameColumns != 3 {
		t.Errorf("Game columns count not expected.  %d", game.GameColumns)
	}

	if game.NumberPossiblePieces != 1 {
		t.Errorf("Game possible pieces not expected.  %d", game.NumberPossiblePieces)
	}

	gameInput <- PlayInputStop
	game = <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateGameOver {
		t.Errorf("Game not over.  %d", game.State)
	}
}

func TestCompleteRow(t *testing.T) {

	gameRoot, gameInput, gameOutput := NewGame()

	game := <-gameOutput

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	gameInput <- PlayInputToggleDrop
	game = <-gameOutput

	gameInput <- PlayInputPause

	// GOAL: make the game look like:
	//  [X X X X X X X X X X X X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 0 0 0 0 0 0 X]
	//  [X * * * * 0 0 0 0 0 0 X]
	//  [X 0 0 0 0 X X X X X X X]
	//  [X X X X X X X X X X X X]

	gameRoot.Piece = 0
	gameRoot.PieceRotation = 0
	gameRoot.PiecePosCol = 1
	gameRoot.PiecePosRow = 17

	gameRoot.Field[18] = []int{1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1}
	//log.Printf("1: %s", gameRoot.GetDebugState())

	gameInput <- PlayInputPause
	game = <-gameOutput
	log.Printf("2: %s", game.GetDebugState())
	log.Printf("A: %+v", game.Field[18])

	gameInput <- PlayInputDrop
	game = <-gameOutput
	log.Printf("3: %s", game.GetDebugState())
	log.Printf("B: %+v", game.Field[18])
	expectedRow := []int{1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1}
	if !slices.Equal(game.Field[18], expectedRow) {
		t.Errorf("Row not as expected.  got: %+v  want %+v", game.Field[18], expectedRow)
	}

	gameInput <- PlayInputDrop
	game = <-gameOutput
	log.Printf("4: %s", game.GetDebugState())
	log.Printf("C: %+v", game.Field[18])
	expectedRow = []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if !slices.Equal(game.Field[18], expectedRow) {
		t.Errorf("Row not as expected.  got: %+v  want %+v", game.Field[18], expectedRow)
	}
}

func TestGameOver(t *testing.T) {
	gameRoot, gameInput, gameOutput := NewGame()

	gameInput <- PlayInputPause
	game := <-gameOutput

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	for i := 1; i < 18; i++ {
		gameRoot.Field[i][5] = 1
		gameRoot.Field[i][6] = 1
	}

	gameInput <- PlayInputPause
	game = <-gameOutput

	gameInput <- PlayInputDrop
	game = <-gameOutput

	if game.State != StateGameOver {
		t.Errorf("Game not over.  %d", game.State)
	}
}

//	func TestGetDebugState(t *testing.T) {
//		// FIXME
//	}

func TestMove(t *testing.T) {

	gameRoot, gameInput, gameOutput := NewGame()

	game := <-gameOutput

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	gameInput <- PlayInputToggleDrop
	game = <-gameOutput

	gameRoot.Piece = 0
	gameRoot.PieceRotation = 0

	posCol := game.PiecePosCol
	posColExpected := posCol - 1

	gameInput <- PlayInputMoveLeft
	game = <-gameOutput
	posCol = game.PiecePosCol

	if posCol != posColExpected {
		t.Errorf("Piece did not move left.  got: %d  expected: %d", posCol, posColExpected)
	}

	for i := 12; i > 0; i-- {
		gameInput <- PlayInputMoveLeft
		game = <-gameOutput
	}
	posCol = game.PiecePosCol
	posColExpected = 1

	if posCol != posColExpected {
		t.Errorf("Piece not on left wall.  got: %d  expected: %d", posCol, posColExpected)
	}

	gameInput <- PlayInputMoveRight
	game = <-gameOutput
	posCol = game.PiecePosCol
	posColExpected = 2

	if posCol != posColExpected {
		t.Errorf("Piece did not move right: got: %d  expected: %d", posCol, posColExpected)
	}

	for i := 12; i > 0; i-- {
		gameInput <- PlayInputMoveRight
		game = <-gameOutput
	}
	posCol = game.PiecePosCol
	posColExpected = 7

	if posCol != posColExpected {
		t.Errorf("Piece not on right wall.  got: %d  expected: %d", posCol, posColExpected)
	}

	gameRoot.PieceRotation = 1

	for i := 12; i > 0; i-- {
		gameInput <- PlayInputMoveLeft
		game = <-gameOutput
	}
	posCol = game.PiecePosCol
	posColExpected = 0

	if posCol != posColExpected {
		t.Errorf("Piece not on left wall.  got: %d  expected: %d", posCol, posColExpected)
	}

	for i := 12; i > 0; i-- {
		gameInput <- PlayInputMoveRight
		game = <-gameOutput
	}
	posCol = game.PiecePosCol
	posColExpected = 9

	if posCol != posColExpected {
		t.Errorf("Piece not on right wall.  got: %d  expected: %d", posCol, posColExpected)
	}

	game = <-gameOutput

	posRow := game.PiecePosRow
	posRowExpected := posRow + 1

	gameInput <- PlayInputDrop
	game = <-gameOutput
	posRow = game.PiecePosRow

	if posRow != posRowExpected {
		t.Errorf("Piece did not drop.  got: %d  expected: %d", posRow, posRowExpected)
	}
}

func TestPauseGame(t *testing.T) {

	_, gameInput, gameOutput := NewGame()

	game := <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	posRow := game.PiecePosRow
	posRowExpected := posRow

	gameInput <- PlayInputPause
	time.Sleep(3 * time.Second)

	gameInput <- PlayInputDrop

	gameInput <- PlayInputPause
	game = <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	if posRow != posRowExpected {
		t.Errorf("piece not at expected positioon.  got: %d  want: %d", posRow, posRowExpected)
	}
}

func TestRotate(t *testing.T) {

	_, gameInput, gameOutput := NewGame()

	game := <-gameOutput

	for i := 0; i < 6; i++ {
		expected := i % 4
		if expected != game.PieceRotation {
			t.Errorf("Rotation not expected. %d, %d", expected, game.PieceRotation)
		}

		gameInput <- PlayInputRotate
		game = <-gameOutput
	}
}
