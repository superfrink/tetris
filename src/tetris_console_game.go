package main

import (
	"./engine"
	"fmt"
	"github.com/rthornton128/goncurses"
	"log"
)

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

	stdscr.Refresh()

	// GOAL: Create a channel for user input
	local_user_input_ch := make(chan goncurses.Key)

	go func() {
		for {
			local_user_input_ch <- stdscr.GetChar()
		}
	}()
	var key goncurses.Key

	// GOAL: Create an instance of the game
	g, game_user_input_ch, game_output_channel := engine.NewGame()

	// Main game loop
	quit := false
	for !quit {
		select {
		case key = <-local_user_input_ch:
			game_user_input_ch <- byte(key)

		case _ = <-game_output_channel:

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

		}

		// GOAL: Update the screen
		stdscr.Refresh()

		// GOAL: Check if the game is over
		if engine.StateGameOver == g.State {
			quit = true
		}
	}
	stdscr.MovePrint(12, 20, "GAME OVER")
	stdscr.MovePrint(13, 17, "press any key")
	stdscr.GetChar()
}
