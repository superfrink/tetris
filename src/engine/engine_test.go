package engine

import (
	"log"
	"testing"
)

func TestCreateGame(t *testing.T) {

	_, gameInput, gameOutput := NewGame()

	game := <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateRunning {
		t.Errorf("Game state not running.  %d", game.State)
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

	gameInput <- PlayInputQuit
	game = <-gameOutput
	log.Printf("%+v", game)

	if game.State != StateGameOver {
		t.Errorf("Game state not over.  %d", game.State)
	}
}
