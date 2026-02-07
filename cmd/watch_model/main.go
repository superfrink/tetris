package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/pkg/player"
)

func printGame(g *engine.Game) {
	// Clear the console
	fmt.Print("\033[H\033[2J")

	// Create a buffer to draw the game state
	field := make([][]rune, g.GameRows+2)
	for i := range field {
		field[i] = make([]rune, g.GameColumns+2)
		for j := range field[i] {
			if g.Field[i][j] != 0 {
				field[i][j] = 'X'
			} else {
				field[i][j] = ' '
			}
		}
	}
	// Draw the borders
	for i := 0; i < g.GameRows+2; i++ {
		field[i][0] = '|'
		field[i][g.GameColumns+1] = '|'
	}
	for j := 0; j < g.GameColumns+2; j++ {
		field[0][j] = '-'
		field[g.GameRows+1][j] = '-'
	}
	field[0][0] = '+'
	field[0][g.GameColumns+1] = '+'
	field[g.GameRows+1][0] = '+'
	field[g.GameRows+1][g.GameColumns+1] = '+'


	// Draw the current piece
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if g.PieceMap[g.Piece][g.PieceRotation][i][j] != 0 {
				if g.PiecePosRow+i < len(field) && g.PiecePosCol+j < len(field[0]) {
					field[g.PiecePosRow+i][g.PiecePosCol+j] = '*'
				}
			}
		}
	}

	// Print the buffer
	for _, row := range field {
		fmt.Println(string(row))
	}

	// Print score
	fmt.Printf("Lines: %d\n", g.ScoreLineCount)
}

func main() {
	// Load the trained model
	genome, mask, err := loadGenome("trained_model.json")
	if err != nil {
		fmt.Printf("Error loading trained model: %v\n", err)
		os.Exit(1)
	}

	// Initialize the game
	game := engine.NewGame()
	game.State = engine.StateRunning

	// Main game loop
	for game.State != engine.StateGameOver {
		// Find the best move for the current piece
		move := player.FindBestMove(game, genome, mask)

		// Apply the chosen rotation and horizontal position
		for i := 0; i < move.Rotation; i++ {
			game.RotatePiece()
		}
		game.PiecePosCol = move.XOffset

		// Animate the piece dropping until it lands
		for {
			printGame(game)
			time.Sleep(50 * time.Millisecond)

			currentPiece := game.Piece
			currentRow := game.PiecePosRow

			game.Step(engine.PlayInputDrop)

			if game.State == engine.StateGameOver || game.Piece != currentPiece || game.PiecePosRow == currentRow {
				break
			}
		}
	}

	// Print final state and score
	printGame(game)
	fmt.Println("GAME OVER")
	fmt.Printf("Final Score (Lines): %d\n", game.ScoreLineCount)
}

type TrainedModel struct {
	Weights []float64 `json:"weights"`
	Mask    string    `json:"mask"`
}

func loadGenome(filename string) ([]float64, []bool, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	var model TrainedModel
	if err := json.Unmarshal(data, &model); err == nil && len(model.Weights) > 0 {
		var mask []bool
		if model.Mask != "" {
			mask = make([]bool, len(model.Mask))
			for i, c := range model.Mask {
				mask[i] = c == '1'
			}
		}
		return model.Weights, mask, nil
	}

	var genome []float64
	if err := json.Unmarshal(data, &genome); err != nil {
		return nil, nil, err
	}

	return genome, nil, nil
}
