package engine

import (
	"bytes"
	// "errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	GameRows             = 18
	GameColumns          = 10
	NumberPossiblePieces = 7
)

type gamestate int

const (
	StateInitializing = iota
	StateRunning
	StateGameOver
)

type Game struct {
	Seed            int64
	State           gamestate
	PRNG            *rand.Rand
	Piece           int
	PieceRotation   int
	PiecePosCol     int
	PiecePosRow     int
	Field           [GameRows + 2][GameColumns + 2]int
	ScorePieceCount int
	ScoreLineCount  int
}

var PieceMap = [7][4][4][4]int{
	{
		// I piece
		{
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
		},
		{
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
		},
	},
	{
		// J piece
		{
			{1, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 0, 0, 0},
			{1, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 1, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
	},
	{
		// L piece
		{
			{1, 1, 1, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 1, 0},
			{1, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
		},
	},
	{
		// O piece
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
	},
	{
		// S piece
		{
			{0, 1, 1, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		},
	},
	{
		// T piece
		{
			{1, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{1, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 0, 0, 0},
			{1, 1, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
	},
	{
		// Z piece
		{
			{1, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
	},
}

func NewGame(player_input <-chan byte, game_state chan<- Game) *Game {
	g := Game{
		Seed:          0,
		State:         StateInitializing,
		Piece:         1,
		PieceRotation: 0,
		PiecePosCol:   4,
		PiecePosRow:   1,
	}

	source := rand.NewSource(g.Seed)
	g.PRNG = rand.New(source)

	for j := 0; j < GameColumns+2; j++ {
		g.Field[0][j] = 1
		g.Field[GameRows+1][j] = 1
	}
	for i := 1; i < GameRows+1; i++ {
		g.Field[i][0] = 1
		g.Field[i][GameColumns+1] = 1
		for j := 1; j < GameColumns+1; j++ {
			g.Field[i][j] = 0
		}
	}

	// GOAL: Start the main game loop
	g.MainGameLoop(player_input, game_state)

	return &g
}

func (g *Game) GetDebugState() string {

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintln("Game {"))

	buffer.WriteString(fmt.Sprintf("Seed: %d\n", g.Seed))
	buffer.WriteString(fmt.Sprintf("Piece: %d\t", g.Piece))
	buffer.WriteString(fmt.Sprintf("Rotation: %d\n", g.PieceRotation))
	buffer.WriteString(fmt.Sprintf("PosCol: %d\n", g.PiecePosCol))
	buffer.WriteString(fmt.Sprintf("PosRow: %d\n", g.PiecePosRow))

	// draw the piece on the field
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != PieceMap[g.Piece][g.PieceRotation][i][j] {
				g.Field[g.PiecePosRow+i+1][g.PiecePosCol+j+1] = 2
			}
		}
	}

	buffer.WriteString("Field:\n")
	for i := 0; i < GameRows+2; i++ {
		buffer.WriteString("    ")
		for j := 0; j < GameColumns+2; j++ {
			if 0 == g.Field[i][j] {
				buffer.WriteString(" ")
			} else {
				buffer.WriteString("X")
			}
		}
		buffer.WriteString("\n")
	}

	// remove the piece on the field data structure
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != PieceMap[g.Piece][g.PieceRotation][i][j] {
				g.Field[g.PiecePosRow+i+1][g.PiecePosCol+j+1] = 0
			}
		}
	}

	buffer.WriteString(fmt.Sprintln("}"))

	return buffer.String()
}

func PieceCollision(field [GameRows + 2][GameColumns + 2]int, piece int, rotation int, row int, col int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != PieceMap[piece][rotation][i][j] {
				if 0 != field[row+i][col+j] {
					return true
				}
			}
		}
	}

	// CLAIM: no collision found
	return false
}

func (g *Game) Rotate() {

	if !PieceCollision(g.Field, g.Piece, (g.PieceRotation+1)%4, g.PiecePosRow, g.PiecePosCol) {
		g.PieceRotation = (g.PieceRotation + 1) % 4
	}
}

func (g *Game) MoveLeft() {
	if !PieceCollision(g.Field, g.Piece, g.PieceRotation, g.PiecePosRow, g.PiecePosCol-1) {
		g.PiecePosCol--
	}
}

func (g *Game) MoveRight() {
	if !PieceCollision(g.Field, g.Piece, g.PieceRotation, g.PiecePosRow, g.PiecePosCol+1) {
		g.PiecePosCol++
	}
}

func (g *Game) LowerPiece() bool {
	// Returns false if unable to lower peice because of collision.
	// Returns true otherwise.

	// GOAL: lower the piece one step

	if PieceCollision(g.Field, g.Piece, g.PieceRotation, g.PiecePosRow+1, g.PiecePosCol) {
		// CLAIM: Piece will collides if lowered.
		return false
	}

	g.PiecePosRow++
	return true
}

func (g *Game) PlacePiece() {

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != PieceMap[g.Piece][g.PieceRotation][i][j] {
				g.Field[g.PiecePosRow+i][g.PiecePosCol+j] = 1
			}
		}
	}
}

func (g *Game) NextPiece() error {
	// GOAL: pick a new random piece

	x := g.PRNG.Intn(NumberPossiblePieces)

	g.Piece = int(x)

	g.PiecePosCol = 4
	g.PiecePosRow = 1
	g.PieceRotation = 0

	return nil
}

func (g *Game) ClearCompletedRows() {

	for i := 1; i < GameRows+1; i++ {

		row_complete := true

		for j := 1; j < GameColumns+1; j++ {
			if 0 == g.Field[i][j] {
				row_complete = false
			}
		}

		if row_complete {

			// GOAL: drop all rows above this one down one row.
			g.ShiftRowsDown(i)

			g.ScoreLineCount++
		}
	}
}

func (g *Game) ShiftRowsDown(start_row int) {
	// Shifts down all rows above the start_row.

	for i := start_row; 1 < i; i-- {
		for j := 1; j < GameColumns+1; j++ {
			g.Field[i][j] = g.Field[i-1][j]
		}
	}
}

func (g *Game) MainGameLoop(player_input <-chan byte, game_state chan<- Game) {
	// This function provides the main game loop logic.
	// It reads player input from channel player_input.
	// It sends game state to channel game_state.

	// GOAL: Create a channel for a ticker to drop the pieces
	ticker := time.NewTicker(time.Millisecond * 500) // FIXME: use a global/config for drop speed

	var key byte
	go func() {

		quit := false
		g.State = StateRunning

		for !quit {

			select {
			case key = <-player_input:
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
						quit = true
					}
					g.NextPiece()
				}
			}
			game_state <- *g
		}

		g.State = StateGameOver
		game_state <- *g
	}()
}
