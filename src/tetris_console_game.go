package main

import (
	"./engine"
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
)

// tbprint is based on from https://github.com/jjinux/gotetris/
func tbprint(y int, x int, str string) {
	for _, c := range str {
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorBlack)
		x++
	}
}

func main() {

	// GOAL: Setup the screen
	err := termbox.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer termbox.Close()

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// GOAL: Setup the keystroke legend
	// FIXME: It would be good to use variables for each key
	tbprint(21, 0, "q = quit\tr = rotate\th = left\tl = right")

	termbox.Flush()

	// GOAL: Create a channel for user input
	local_user_input_ch := make(chan rune)

	go func() {
		for {
			event := termbox.PollEvent()
			if termbox.EventKey == event.Type {
				local_user_input_ch <- event.Ch
			}
		}
	}()
	var key rune

	// GOAL: Create an instance of the game
	g, game_user_input_ch, game_output_channel := engine.NewGame()

	// Main game loop
	quit := false
mainloop:
	for {
		select {

		case key = <-local_user_input_ch:
			if quit {
				break mainloop
			}

			switch key {
			case 'q':
				game_user_input_ch <- engine.PlayInputQuit
			case 'h':
				game_user_input_ch <- engine.PlayInputMoveLeft
			case 'l':
				game_user_input_ch <- engine.PlayInputMoveRight
			case 'r':
				game_user_input_ch <- engine.PlayInputRotate
			}

		case _ = <-game_output_channel:

			// GOAL: Draw the game field
			for i := 0; i < engine.GameRows+2; i++ {
				for j := 0; j < engine.GameColumns+2; j++ {
					if 0 != g.Field[i][j] {
						tbprint(i, j, "X")
					} else {
						tbprint(i, j, " ")
					}
				}
			}

			// GOAL: Draw the piece in play
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if 0 != engine.PieceMap[g.Piece][g.PieceRotation][i][j] {
						tbprint(g.PiecePosRow+i, g.PiecePosCol+j, "*")
					}
				}
			}

			// GOAL: Draw the score
			tbprint(2, 15, fmt.Sprintf("Pieces: %d", g.ScorePieceCount))
			tbprint(3, 15, fmt.Sprintf("Lines:  %d", g.ScoreLineCount))

			if true {
				// FIXME: only show when debugging
				tbprint(7, 15, fmt.Sprintf("Piece    : %2d", g.Piece))
				tbprint(8, 15, fmt.Sprintf("Rotation : %2d", g.PieceRotation))
				tbprint(9, 15, fmt.Sprintf("Piece row: %2d", g.PiecePosRow))
				tbprint(10, 15, fmt.Sprintf("Piece col: %2d", g.PiecePosCol))
			}

		}

		// GOAL: Update the screen
		termbox.Flush()

		// GOAL: Check if the game is over
		if engine.StateGameOver == g.State {
			tbprint(12, 20, "GAME OVER")
			tbprint(13, 17, "press any key")
			quit = true
		}
	}
}
