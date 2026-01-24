package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/streamer"

	"github.com/nsf/termbox-go"
)

type MainState struct {
	mu             sync.Mutex
	display        Display               // for drawing the game
	gameState      *engine.Game          // stores intermediate game state for rendering
	stream         streamer.Streamer     // NATS streamer
	streamMesgChan chan streamer.Message // messages from streamer
}

var mainState MainState

func (m *MainState) startDisplayingRemoteGame(localUserInputChan chan rune) {
	m.display.TBPrint(21, 0, "q = quit")
	m.display.Flush()

	for {
		select {

		case key := <-localUserInputChan:
			if key == 'q' {
				exitProgram()
			}
		case message := <-m.streamMesgChan:
			if message.Type == streamer.StateUpdate {
				drawGameState(message.Game)
			}
		}
	}
}

func exitProgram() {
	mainState.display.Close()
	os.Exit(0)
}

func drawPressAnyKey() {
	mainState.display.TBPrint(13, 17, "press any key")
}

func drawGameState(gameState engine.Game) {
	// GOAL: Draw the game field
	for i := 0; i < gameState.GameRows+2; i++ {
		for j := 0; j < gameState.GameColumns+2; j++ {
			if 0 != gameState.Field[i][j] {
				mainState.display.TBPrint(i, j, "X")
			} else {
				mainState.display.TBPrint(i, j, " ")
			}
		}
	}

	// GOAL: Draw the piece in play
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != gameState.PieceMap[gameState.Piece][gameState.PieceRotation][i][j] {
				mainState.display.TBPrint(gameState.PiecePosRow+i, gameState.PiecePosCol+j, "*")
			}
		}
	}

	// GOAL: Draw the score
	mainState.display.TBPrint(2, 15, fmt.Sprintf("Pieces: %d", gameState.ScorePieceCount))
	mainState.display.TBPrint(3, 15, fmt.Sprintf("Lines:  %d", gameState.ScoreLineCount))

	if true {
		// FIXME: only show when debugging
		mainState.display.TBPrint(7, 15, fmt.Sprintf("Piece    : %2d", gameState.Piece))
		mainState.display.TBPrint(8, 15, fmt.Sprintf("Rotation : %2d", gameState.PieceRotation))
		mainState.display.TBPrint(9, 15, fmt.Sprintf("Piece row: %2d", gameState.PiecePosRow))
		mainState.display.TBPrint(10, 15, fmt.Sprintf("Piece col: %2d", gameState.PiecePosCol))
	}

	// GOAL: Check if the game is over
	if engine.StateGameOver == gameState.State {
		mainState.display.TBPrint(12, 20, "GAME OVER")
	} else {
		// DOC: The stream-viewing display experiences this when a new game is started.
		mainState.display.TBPrint(12, 20, "         ")
	}

	// GOAL: Update the screen
	mainState.display.Flush()
}

func main() {

	var flagBucketgame = flag.Bool("b", false, "Play a bucket game instead.")
	var natsUrl = flag.String("u", "", "NATS URL")
	var natsCredFile = flag.String("c", "", "NATS credential file")
	var streamGame = flag.Bool("s", false, "Send game stream")
	var viewStream = flag.Bool("v", false, "Watch streaming game")
	flag.Parse()

	if *streamGame || *viewStream {
		mainState.stream = streamer.Streamer{}
		mainState.stream.Connect(*natsUrl, *natsCredFile, "fixme")
		mainState.streamMesgChan = mainState.stream.RecvChan()
	}

	// GOAL: Setup the screen
	err := mainState.display.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer mainState.display.Close()

	mainState.display.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// GOAL: Create a channel for user input
	localUserInputChan := make(chan rune)

	go func() {
		for {
			event := mainState.display.PollEvent()
			if termbox.EventKey == event.Type {
				localUserInputChan <- event.Ch
			}
		}
	}()

	// GOAL: Play the game viewer if requested
	if *viewStream {
		mainState.startDisplayingRemoteGame(localUserInputChan) // does not return
	}

	// CLAIM: Game viewer not requested so play the game instead

	// GOAL: Create an instance of the game
	if *flagBucketgame {
		mainState.gameState = engine.NewBucketGame()
	} else {
		mainState.gameState = engine.NewGame()
	}

	// Main game

	// GOAL: Setup the keystroke legend
	mainState.display.TBPrint(21, 0, "q = quit\tr = rotate\th = left\tl = right\td = drop\tp = pause")
	mainState.display.Flush()

	// GOAL: start the game loop
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	quitCount := 0

	for {
		// Draw the initial state
		drawGameState(*mainState.gameState)

		select {
		case key := <-localUserInputChan:
			var move byte
			nop := false
			switch key {
			case 'q':
				quitCount++
				if quitCount > 1 {
					exitProgram()
				}
				drawPressAnyKey()
				nop = true
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
				mainState.gameState.Step(move)
				if *streamGame {
					mainState.stream.SendMove(move, *mainState.gameState)
				}
			}

		case <-ticker.C:
			mainState.gameState.Step(engine.PlayInputDrop)
			if *streamGame {
				mainState.stream.SendGameState(*mainState.gameState)
			}
		}

		if mainState.gameState.State == engine.StateGameOver {
			// let the loop run one more time to draw the game over message
			<-localUserInputChan
			exitProgram()
		}
	}
}
