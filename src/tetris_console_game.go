package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/nsf/termbox-go"
	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/streamer"
)

// DOC: Display wraps a mutex around termbox writes.
type Display struct {
	mu sync.Mutex
}

func (d *Display) Clear(a termbox.Attribute, b termbox.Attribute) {
	d.mu.Lock()
	termbox.Clear(a, b)
	d.mu.Unlock()
}

func (d *Display) Close() {
	d.mu.Lock()
	termbox.Close()
	d.mu.Unlock()
}

func (d *Display) Flush() {
	d.mu.Lock()
	termbox.Flush()
	d.mu.Unlock()
}

func (d *Display) Init() error {
	d.mu.Lock()
	val := termbox.Init()
	d.mu.Unlock()
	return val
}

func (d *Display) PollEvent() termbox.Event {
	val := termbox.PollEvent()
	return val
}

func (d *Display) TBPrint(y int, x int, str string) {
	d.mu.Lock()
	// tbprint is based on from https://github.com/jjinux/gotetris/
	for _, c := range str {
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorBlack)
		x++
	}
	d.mu.Unlock()
}

var display Display

func exitProgram() {
	display.Close()
	os.Exit(0)
}

func drawPressAnyKey() {
	display.TBPrint(13, 17, "press any key")
}

func drawGameState(gameState engine.Game) {
	// GOAL: Draw the game field
	for i := 0; i < gameState.GameRows+2; i++ {
		for j := 0; j < gameState.GameColumns+2; j++ {
			if 0 != gameState.Field[i][j] {
				display.TBPrint(i, j, "X")
			} else {
				display.TBPrint(i, j, " ")
			}
		}
	}

	// GOAL: Draw the piece in play
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != gameState.PieceMap[gameState.Piece][gameState.PieceRotation][i][j] {
				display.TBPrint(gameState.PiecePosRow+i, gameState.PiecePosCol+j, "*")
			}
		}
	}

	// GOAL: Draw the score
	display.TBPrint(2, 15, fmt.Sprintf("Pieces: %d", gameState.ScorePieceCount))
	display.TBPrint(3, 15, fmt.Sprintf("Lines:  %d", gameState.ScoreLineCount))

	if true {
		// FIXME: only show when debugging
		display.TBPrint(7, 15, fmt.Sprintf("Piece    : %2d", gameState.Piece))
		display.TBPrint(8, 15, fmt.Sprintf("Rotation : %2d", gameState.PieceRotation))
		display.TBPrint(9, 15, fmt.Sprintf("Piece row: %2d", gameState.PiecePosRow))
		display.TBPrint(10, 15, fmt.Sprintf("Piece col: %2d", gameState.PiecePosCol))
	}

	// GOAL: Check if the game is over
	if engine.StateGameOver == gameState.State {
		display.TBPrint(12, 20, "GAME OVER")
	}

	// GOAL: Update the screen
	display.Flush()
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
	err := display.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer display.Close()

	display.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// GOAL: Create a channel for user input
	localUserInputChan := make(chan rune)

	go func() {
		for {
			event := display.PollEvent()
			if termbox.EventKey == event.Type {
				localUserInputChan <- event.Ch
			}
		}
	}()

	if *viewStream {
		display.TBPrint(21, 0, "q = quit")

		for {
			select {

			case key := <-localUserInputChan:
				if key == 'q' {
					exitProgram()
				}
			case message := <-streamMesgChan:
				if message.Type == streamer.StateUpdate {
					drawGameState(message.Game)
				}
			}

			// GOAL: Update the screen
			display.Flush()
		}
	}

	// GOAL: Setup the keystroke legend
	display.TBPrint(21, 0, "q = quit\tr = rotate\th = left\tl = right\td = drop\tp = pause")

	display.Flush()
	// GOAL: Create an instance of the game
	var gameState *engine.Game
	var gameUserInputChan chan<- byte
	var gameOutputChan <-chan *engine.Game

	if *flagBucketgame {
		_, gameUserInputChan, gameOutputChan = engine.NewBucketGame()
	} else {
		_, gameUserInputChan, gameOutputChan = engine.NewGame()
	}

	// Main game

	// GOAL: Send commands from keypresses to the game
	go func(userInputChan chan rune, gameCommandChan chan<- byte, stream *streamer.Streamer) {
		quit := 0

		for {
			nop := false

			var move byte

			key := <-localUserInputChan

			switch key {
			case 'q':
				move = engine.PlayInputStop
				quit++
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

			// DOC: 2 goroutines would be using the same gameState variable
			// if *streamGame {
			// 	stream.SendMove(move, *gameState)
			// }

			switch quit {
			case 1:
				drawPressAnyKey()
			case 2:
				exitProgram()
			}
		}
	}(localUserInputChan, gameUserInputChan, &stream)

	// GOAL: Draw the game state updates to the screen
	for {
		gameState = <-gameOutputChan

		drawGameState(*gameState)

		if *streamGame {
			stream.SendGameState(*gameState)
		}
	}
}
