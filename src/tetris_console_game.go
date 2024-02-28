package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
	"superfrink.net/tetris/engine"
)

// tbprint is based on from https://github.com/jjinux/gotetris/
func tbprint(y int, x int, str string) {
	for _, c := range str {
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorBlack)
		x++
	}
}

func main() {

	var flag_bucketgame = flag.Bool("b", false, "Play a bucket game instead.")
	flag.Parse()

	// GOAL: Setup the screen
	err := termbox.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer termbox.Close()

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// GOAL: Setup the keystroke legend
	// FIXME: It would be good to use variables for each key
	tbprint(21, 0, "q = quit\tr = rotate\th = left\tl = right\td = drop\tp = pause")

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
	var game_state *engine.Game
	var game_user_input_ch chan<- byte
	var game_output_channel <-chan *engine.Game

	if *flag_bucketgame {
		_, game_user_input_ch, game_output_channel = engine.NewBucketGame()
	} else {
		_, game_user_input_ch, game_output_channel = engine.NewGame()
	}

	// Wait until the game is ready
	game_state = <-game_output_channel

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
				game_user_input_ch <- engine.PlayInputStop
			case 'd':
				game_user_input_ch <- engine.PlayInputDrop
			case 'h':
				game_user_input_ch <- engine.PlayInputMoveLeft
			case 'l':
				game_user_input_ch <- engine.PlayInputMoveRight
			case 'p':
				game_user_input_ch <- engine.PlayInputPause
			case 'r':
				game_user_input_ch <- engine.PlayInputRotate
			}

		case game_state = <-game_output_channel:

			// GOAL: Draw the game field
			for i := 0; i < game_state.GameRows+2; i++ {
				for j := 0; j < game_state.GameColumns+2; j++ {
					if 0 != game_state.Field[i][j] {
						tbprint(i, j, "X")
					} else {
						tbprint(i, j, " ")
					}
				}
			}

			// GOAL: Draw the piece in play
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if 0 != game_state.PieceMap[game_state.Piece][game_state.PieceRotation][i][j] {
						tbprint(game_state.PiecePosRow+i, game_state.PiecePosCol+j, "*")
					}
				}
			}

			// GOAL: Draw the score
			tbprint(2, 15, fmt.Sprintf("Pieces: %d", game_state.ScorePieceCount))
			tbprint(3, 15, fmt.Sprintf("Lines:  %d", game_state.ScoreLineCount))

			if true {
				// FIXME: only show when debugging
				tbprint(7, 15, fmt.Sprintf("Piece    : %2d", game_state.Piece))
				tbprint(8, 15, fmt.Sprintf("Rotation : %2d", game_state.PieceRotation))
				tbprint(9, 15, fmt.Sprintf("Piece row: %2d", game_state.PiecePosRow))
				tbprint(10, 15, fmt.Sprintf("Piece col: %2d", game_state.PiecePosCol))
			}

		}

		// GOAL: Check if the game is over
		if game_state != nil && engine.StateGameOver == game_state.State {
			tbprint(12, 20, "GAME OVER")
			tbprint(13, 17, "press any key")
			quit = true
		}
		// GOAL: Update the screen
		termbox.Flush()
	}
}
