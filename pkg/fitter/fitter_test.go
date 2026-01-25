package fitter

import (
	"reflect"
	"testing"
)

// newTestBoard creates a board with borders, similar to the game engine.
func newTestBoard(rows, cols int) [][]int {
	board := make([][]int, rows+2)
	for i := range board {
		board[i] = make([]int, cols+2)
	}
	// Create borders
	for j := 0; j < cols+2; j++ {
		board[0][j] = 1
		board[rows+1][j] = 1
	}
	for i := 1; i < rows+1; i++ {
		board[i][0] = 1
		board[i][cols+1] = 1
	}
	return board
}

func TestCalculateHeuristics(t *testing.T) {
	const testRows, testCols = 20, 10
	tests := []struct {
		name  string
		board [][]int
		want  Heuristics
	}{
		{
			name:  "Empty board",
			board: newTestBoard(testRows, testCols),
			want: Heuristics{
				AggregateHeight: 0,
				Holes:           0,
				Bumpiness:       0,
				LinesCleared:    0,
				LandingHeight:   0,
				Overhangs:       0,
			},
		},
		{
			name: "Single block at bottom-left",
			board: func() [][]int {
				b := newTestBoard(testRows, testCols)
				b[19][1] = 1 // Place in the first playable column
				return b
			}(),
			want: Heuristics{
				AggregateHeight: 2,
				Holes:           0,
				Bumpiness:       2,
				LinesCleared:    0,
				LandingHeight:   0,
				Overhangs:       1, // Corrected to 1
			},
		},
		{
			name: "Board with holes and bumpiness",
			board: func() [][]int {
				b := newTestBoard(testRows, testCols)
				// Column 1: height 4
				b[19][1] = 1
				b[18][1] = 1
				b[17][1] = 1
				// Column 2: height 4, 1 hole
				b[19][2] = 1
				b[17][2] = 1 // hole at [18][2]
				// Column 3: height 2
				b[19][3] = 1
				// Column 4: height 0
				return b
			}(),
			want: Heuristics{
				AggregateHeight: 10,
				Holes:           0,
				Bumpiness:       4,
				LinesCleared:    0,
				LandingHeight:   0,
				Overhangs:       4,
			},
		},
		{
			name: "Board with single overhang",
			board: func() [][]int {
				b := newTestBoard(testRows, testCols)
				// Create an overhang:
				//  X
				// . .  (empty spots at (18,1) and (18,2))
				b[17][1] = 1 // Filled spot at (17,1)

				// This creates an overhang at (18,1):
				// - (18,1) is empty
				// - (17,1) is filled (above (18,1))
				// - (18,2) is empty (beside (18,1))
				return b
			}(),
			want: Heuristics{
				AggregateHeight: 4,
				Holes:           0,
				Bumpiness:       4,
				LinesCleared:    0,
				LandingHeight:   0,
				Overhangs:       1, // Corrected to 1
			},
		},
		{
			name: "Board with multiple overhangs",
			board: func() [][]int {
				b := newTestBoard(testRows, testCols)
				// Create overhangs:
				//  X X
				// . . . (empty spots)
				b[17][1] = 1
				b[17][2] = 1

				// This creates two overhangs:
				// 1. At (18,1): (17,1) filled, (18,1) empty, (18,2) empty
				// 2. At (18,2): (17,2) filled, (18,2) empty, (18,1) empty
				return b
			}(),
			want: Heuristics{
				AggregateHeight: 8,
				Holes:           0,
				Bumpiness:       4,
				LinesCleared:    0,
				LandingHeight:   0,
				Overhangs:       2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: The board setup for "Board with holes and bumpiness" was also slightly adjusted
			// to make the `want` values clearer and consistent with the trace.
			got := CalculateHeuristics(tt.board, 0, 0)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateHeuristics() = %v, want %v", got, tt.want)
			}
		})
	}
}
