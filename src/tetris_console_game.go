package main

import (
	"./engine"
	"encoding/json"
	"fmt"
	"github.com/rthornton128/goncurses"
	"log"
	"time"
)

func PrettyPrint(v interface{}) (err error) {
	// from https://siongui.github.io/2016/01/30/go-pretty-print-variable/
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func main() {

	// GOAL: Setup the ncurses screen
	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer goncurses.End()

	goncurses.Echo(false)
	goncurses.CBreak(true)
	goncurses.Cursor(0)

	// GOAL: Setup the keystroke legend
	// FIXME: It would be good to use variables for each key
	stdscr.MovePrint(21, 0, "q = quit\tr = rotate\th = left\tl = right")

	// GOAL: Setup the output screen : game score
	// FIXME: incomplete

	stdscr.Refresh()

	// GOAL: Create a channel for a ticker to drop the pieces
	ticker := time.NewTicker(time.Millisecond * 500)

	// GOAL: Create a channel for user input
	user_input_ch := make(chan goncurses.Key)
	go func() {
		for {
			user_input_ch <- stdscr.GetChar()
		}
	}()
	var key goncurses.Key

	// GOAL: Create an instance of the game
	g := engine.NewGame()

	// Main game loop
	quit := false
gameloop:
	for !quit {
		// GOAL: Draw the game field
		for i := 0; i < engine.GameRows+2; i++ {
			for j := 0; j < engine.GameColumns+2; j++ {
				if 0 != g.Field[i][j] {
					stdscr.MovePrint(i, j, "X")
				} else {
					stdscr.MovePrint(i, j, " ")
				}
			}
		}

		// GOAL: Draw the piece in play
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if 0 != engine.PieceMap[g.Piece][g.PieceRotation][i][j] {
					stdscr.MovePrint(g.PiecePosRow+i, g.PiecePosCol+j, "*")
				}
			}
		}

		// GOAL: Draw the score
		stdscr.MovePrint(2, 15, fmt.Sprintf("Pieces: %d", g.ScorePieceCount))
		stdscr.MovePrint(3, 15, fmt.Sprintf("Lines:  %d", g.ScoreLineCount))

		if true {
			// FIXME: only show when debugging
			stdscr.MovePrint(7, 15, fmt.Sprintf("Piece    : %2d", g.Piece))
			stdscr.MovePrint(8, 15, fmt.Sprintf("Rotation : %2d", g.PieceRotation))
			stdscr.MovePrint(9, 15, fmt.Sprintf("Piece row: %2d", g.PiecePosRow))
			stdscr.MovePrint(10, 15, fmt.Sprintf("Piece col: %2d", g.PiecePosCol))
		}

		// GOAL: Update the screen
		stdscr.Refresh()

		select {
		case key = <-user_input_ch:
			switch key {
			case 'q':
				quit = true
			case 'h':
				g.MoveLeft()
			case 'l':
				g.MoveRight()
			case 'r':
				g.Rotate()
			case 'd':
				//FIXME: drop
			}

		case <-ticker.C:
			// Lower the piece and check if it collides.
			able_to_lower := g.LowerPiece()
			if !able_to_lower {
				g.PlacePiece()
				g.ClearCompletedRows()
				g.ScorePieceCount++

				if 1 == g.PiecePosRow {
					// CLAIM: game over
					break gameloop
				}
				g.NextPiece()
			}
		}
	}

	stdscr.MovePrint(12, 20, "GAME OVER")
	stdscr.GetChar()
}
