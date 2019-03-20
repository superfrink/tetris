package engine

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

// DOC: Constants used in the game engine
const (
	GameRows             = 18
	GameColumns          = 10
	NumberPossiblePieces = 7
)

// DOC: Possible states a game can be in
type gamestate int

const (
	StateInitializing = iota
	StateRunning
	StateGameOver
)

// DOC: Player input commands available
const (
	PlayInputQuit = iota
	PlayInputMoveLeft
	PlayInputMoveRight
	PlayInputRotate
	PlayInputDrop
)

// DOC: Data structure describing a game
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

// DOC: Mapping from piece and rotation to blocks covered
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

// NewGame creates a new instance of a tetris game.
// Returns:
// - A game struct for the new game
// - The input channel that player moves will be read from
// - An output channel that will be sent each state change
func NewGame() (*Game, chan<- byte, <-chan Game) {

	seed := time.Now().UTC().UnixNano()

	return NewSeededGame(seed);
}

// NewSeededGame creates a new instance of a game using the the PRNG seed.
// Inputs:
// - a seed to be used for the generation of pieces to be played
// Returns:
// - A game struct for the new game
// - The input channel that player moves will be read from
// - An output channel that will be sent each state change
func NewSeededGame(seed int64) (*Game, chan<- byte, <-chan Game) {
g := Game{
		Seed:          seed,
		State:         StateInitializing,
		Piece:         1,
		PieceRotation: 0,
		PiecePosCol:   4,
		PiecePosRow:   1,
	}

	source := rand.NewSource(g.Seed)
	g.PRNG = rand.New(source)
	g.nextPiece()

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

	player_input_channel := make(chan byte)
	output_state_channel := make(chan Game)

	// GOAL: Start the main game loop
	g.MainGameLoop(player_input_channel, output_state_channel)

	return &g, player_input_channel, output_state_channel
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

// pieceCollision determines whether a specified piece in the specified position and
// rotation would collide with any existing blocks on the specfied field.
// Returns:
// - true if there is a collision
// - false otherwise
func pieceCollision(field [GameRows + 2][GameColumns + 2]int, piece int, rotation int, row int, col int) bool {
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

// rotate changes the rotation of the piece only if the rotation would not collide.
func (g *Game) rotate() {

	if !pieceCollision(g.Field, g.Piece, (g.PieceRotation+1)%4, g.PiecePosRow, g.PiecePosCol) {
		g.PieceRotation = (g.PieceRotation + 1) % 4
	}
}

// moveLeft move the position to the left only if the move would not collide.
func (g *Game) moveLeft() {
	if !pieceCollision(g.Field, g.Piece, g.PieceRotation, g.PiecePosRow, g.PiecePosCol-1) {
		g.PiecePosCol--
	}
}

// moveRight move the position to the right only if the move would not collide.
func (g *Game) moveRight() {
	if !pieceCollision(g.Field, g.Piece, g.PieceRotation, g.PiecePosRow, g.PiecePosCol+1) {
		g.PiecePosCol++
	}
}

// lowerPiece lowers the position by one step only if the move would not collide.
// Returns:
// - false if a collision would occur.
// - true otherwise
func (g *Game) lowerPiece() bool {
	// Returns false if unable to lower peice because of collision.
	// Returns true otherwise.

	// GOAL: lower the piece one step

	if pieceCollision(g.Field, g.Piece, g.PieceRotation, g.PiecePosRow+1, g.PiecePosCol) {
		// CLAIM: Piece will collides if lowered.
		return false
	}

	g.PiecePosRow++
	return true
}

// placePiece updates the field to place each block from the piece onto the play field.
func (g *Game) placePiece() {

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if 0 != PieceMap[g.Piece][g.PieceRotation][i][j] {
				g.Field[g.PiecePosRow+i][g.PiecePosCol+j] = 1
			}
		}
	}
}

// nextPiece updates the game state to have a new piece in play at the top of the field.
// FIXME: when would this return an error?
func (g *Game) nextPiece() error {
	// GOAL: pick a new random piece

	x := g.PRNG.Intn(NumberPossiblePieces)

	g.Piece = int(x)

	g.PiecePosCol = 4
	g.PiecePosRow = 1
	g.PieceRotation = 0

	return nil
}

// clearCompletedRows finds completed rows in the field, removes them, and drops
// above rows down.
func (g *Game) clearCompletedRows() {

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

// ShiftRowsDown drops blocks down by one row, starting at the start_row.
// Called by clearCompletedRows().
func (g *Game) ShiftRowsDown(start_row int) {
	// Shifts down all rows above the start_row.

	for i := start_row; 1 < i; i-- {
		for j := 1; j < GameColumns+1; j++ {
			g.Field[i][j] = g.Field[i-1][j]
		}
	}
}

// MainGameLoop provides the main game loop logic.
// Reads player input from channel player_input.
// Sends game state to channel game_state.
func (g *Game) MainGameLoop(player_input <-chan byte, game_state chan<- Game) {

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
				case PlayInputQuit:
					quit = true
				case PlayInputMoveLeft:
					g.moveLeft()
				case PlayInputMoveRight:
					g.moveRight()
				case PlayInputRotate:
					g.rotate()
				case PlayInputDrop:
					//FIXME: drop
				}

			case <-ticker.C:
				// Lower the piece and check if it collides.
				able_to_lower := g.lowerPiece()
				if !able_to_lower {
					g.placePiece()
					g.clearCompletedRows()
					g.ScorePieceCount++

					if 1 == g.PiecePosRow {
						// CLAIM: game over
						quit = true
					}
					g.nextPiece()
				}
			}
			game_state <- *g
		}

		g.State = StateGameOver
		game_state <- *g
	}()
}
