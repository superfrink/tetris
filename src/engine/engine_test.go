package engine

import (
	"log"
	"testing"
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

func TestMove(t *testing.T) {

	_, gameInput, gameOutput := NewGame()

	game := <-gameOutput

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	gameInput <- PlayInputToggleDrop
	game = <-gameOutput
	log.Printf("%+v", game)

	posCol1 := game.PiecePosCol

	gameInput <- PlayInputMoveLeft
	game = <-gameOutput
	posCol2 := game.PiecePosCol
	log.Printf("%+v", game)

	if posCol1-1 != posCol2 {
		t.Errorf("Piece did not move left: %d -> %d", posCol1, posCol2)
	}

	gameInput <- PlayInputMoveRight
	game = <-gameOutput
	posCol3 := game.PiecePosCol

	if posCol3-1 != posCol2 {
		t.Errorf("Piece did not move right: %d -> %d", posCol2, posCol3)
	}

	posCol4 := game.PiecePosCol
	for i := posCol3; i > 1; i-- {
		gameInput <- PlayInputMoveLeft
		game = <-gameOutput
		posCol4 = game.PiecePosCol
	}
	if posCol4 != 1 {
		t.Errorf("Piece should be at pos 1: %d", posCol4)
	}

	gameInput <- PlayInputMoveLeft
	game = <-gameOutput
	posCol4 = game.PiecePosCol
	if posCol4 != 1 {
		// FIXME: test moving all the way to the left.  colum number depends on peice and rotation.
		// t.Errorf("Piece should still be at pos 1: %d", posCol4)
	}

	// FIXME: test moving all the way to the right.  colum number depends on peice and rotation.

	gameInput <- PlayInputPause
	game = <-gameOutput

	if game.State != StatePaused {
		t.Errorf("Game not paused.  %d", game.State)
	}
	gameInput <- PlayInputPause
	game = <-gameOutput

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}
}

func TestPauseGame(t *testing.T) {

	_, gameInput, gameOutput := NewGame()

	game := <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
	}

	gameInput <- PlayInputPause
	game = <-gameOutput
	log.Printf("%+v", game)

	if game.State != StatePaused {
		t.Errorf("Game not paused.  %d", game.State)
	}

	gameInput <- PlayInputPause
	game = <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateRunning {
		t.Errorf("Game not running.  %d", game.State)
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
