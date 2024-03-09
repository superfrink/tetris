package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nsf/termbox-go"
	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/streamer"
)

// tbprint is based on from https://github.com/jjinux/gotetris/
func tbprint(y int, x int, str string) {
	for _, c := range str {
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorBlack)
		x++
	}
}

func drawGameState(gameState engine.Game) {
	// GOAL: Draw the game field
	for i := 0; i < gameState.GameRows+2; i++ {
		for j := 0; j < gameState.GameColumns+2; j++ {
			if 0 != gameState.Field[i][j] {
				tbprint(i, j, "X")
			} else {
				tbprint(i, j, " ")
			}
		}
	}

	// GOAL: Draw the piece in play
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != gameState.PieceMap[gameState.Piece][gameState.PieceRotation][i][j] {
				tbprint(gameState.PiecePosRow+i, gameState.PiecePosCol+j, "*")
			}
		}
	}

	// GOAL: Draw the score
	tbprint(2, 15, fmt.Sprintf("Pieces: %d", gameState.ScorePieceCount))
	tbprint(3, 15, fmt.Sprintf("Lines:  %d", gameState.ScoreLineCount))

	if true {
		// FIXME: only show when debugging
		tbprint(7, 15, fmt.Sprintf("Piece    : %2d", gameState.Piece))
		tbprint(8, 15, fmt.Sprintf("Rotation : %2d", gameState.PieceRotation))
		tbprint(9, 15, fmt.Sprintf("Piece row: %2d", gameState.PiecePosRow))
		tbprint(10, 15, fmt.Sprintf("Piece col: %2d", gameState.PiecePosCol))
	}
}

func main() {

	var flagBucketgame = flag.Bool("b", false, "Play a bucket game instead.")
	var natsUrl = flag.String("u", "", "NATS URL")
	var natsCredFile = flag.String("c", "", "NATS credential file")
	var streamGame = flag.Bool("s", false, "Send game stream")
	var viewStream = flag.Bool("v", false, "Watch streaming game")
	flag.Parse()

	var stream streamer.Streamer
	var streamMesgChan chan streamer.Message
	if *streamGame || *viewStream {
		stream = streamer.Streamer{}
		stream.Connect(*natsUrl, *natsCredFile, "fixme")
		streamMesgChan = stream.RecvChan()
	}

	// GOAL: Setup the screen
	err := termbox.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer termbox.Close()

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// GOAL: Create a channel for user input
	localUserInputChan := make(chan rune)

	go func() {
		for {
			event := termbox.PollEvent()
			if termbox.EventKey == event.Type {
				localUserInputChan <- event.Ch
			}
		}
	}()

	if *viewStream {
		tbprint(21, 0, "q = quit")

		for {
			select {

			case key := <-localUserInputChan:
				if key == 'q' {
					termbox.Close()
					os.Exit(0)
				}
			case message := <-streamMesgChan:
				drawGameState(message.Game)
			}

			// GOAL: Update the screen
			termbox.Flush()
		}
	}

	// GOAL: Setup the keystroke legend
	// FIXME: It would be good to use variables for each key
	tbprint(21, 0, "q = quit\tr = rotate\th = left\tl = right\td = drop\tp = pause")

	termbox.Flush()
	// GOAL: Create an instance of the game
	var gameState *engine.Game
	var gameUserInputChan chan<- byte
	var gameOutputChan <-chan *engine.Game

	if *flagBucketgame {
		_, gameUserInputChan, gameOutputChan = engine.NewBucketGame()
	} else {
		_, gameUserInputChan, gameOutputChan = engine.NewGame()
	}

	// Wait until the game is ready
	gameState = <-gameOutputChan

	// Main game loop
	quit := false
	var move byte
mainloop:
	for {

		select {

		case key := <-localUserInputChan:
			nop := false

			if quit {
				break mainloop
			}

			switch key {
			case 'q':
				move = engine.PlayInputStop
			case 'd':
				move = engine.PlayInputDrop
			case 'h':
				move = engine.PlayInputMoveLeft
			case 'l':
				move = engine.PlayInputMoveRight
			case 'p':
				move = engine.PlayInputPause
			case 'r':
				move = engine.PlayInputRotate
			default:
				nop = true
			}

			if !nop {
				gameUserInputChan <- move
			}

			if *streamGame {
				stream.SendMove(move, *gameState)
			}

		case gameState = <-gameOutputChan:

			// GOAL: Draw the game field
			for i := 0; i < gameState.GameRows+2; i++ {
				for j := 0; j < gameState.GameColumns+2; j++ {
					if 0 != gameState.Field[i][j] {
						tbprint(i, j, "X")
					} else {
						tbprint(i, j, " ")
					}
				}
			}

			// GOAL: Draw the piece in play
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if 0 != gameState.PieceMap[gameState.Piece][gameState.PieceRotation][i][j] {
						tbprint(gameState.PiecePosRow+i, gameState.PiecePosCol+j, "*")
					}
				}
			}

			// GOAL: Draw the score
			tbprint(2, 15, fmt.Sprintf("Pieces: %d", gameState.ScorePieceCount))
			tbprint(3, 15, fmt.Sprintf("Lines:  %d", gameState.ScoreLineCount))

			if true {
				// FIXME: only show when debugging
				tbprint(7, 15, fmt.Sprintf("Piece    : %2d", gameState.Piece))
				tbprint(8, 15, fmt.Sprintf("Rotation : %2d", gameState.PieceRotation))
				tbprint(9, 15, fmt.Sprintf("Piece row: %2d", gameState.PiecePosRow))
				tbprint(10, 15, fmt.Sprintf("Piece col: %2d", gameState.PiecePosCol))
			}

			if *streamGame {
				stream.SendGameState(*gameState)
			}
		}

		// GOAL: Check if the game is over
		if gameState != nil && engine.StateGameOver == gameState.State {
			tbprint(12, 20, "GAME OVER")
			tbprint(13, 17, "press any key")
			quit = true
		}
		// GOAL: Update the screen
		termbox.Flush()
	}
}
