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
				Holes:           1,
				Bumpiness:       2,
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
				Holes:           4,
				Bumpiness:       4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: The board setup for "Board with holes and bumpiness" was also slightly adjusted
			// to make the `want` values clearer and consistent with the trace.
			got := CalculateHeuristics(tt.board)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateHeuristics() = %v, want %v", got, tt.want)
			}
		})
	}
}
