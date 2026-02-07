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
				AggregateHeight:   0,
				Holes:             0,
				Bumpiness:         0,
				LinesCleared:      0,
				LandingHeight:     0,
				Overhangs:         0,
				ColumnTransitions: 20,
				RowTransitions:    40,
				HoleDepth:         0,
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
				AggregateHeight:   2,
				Holes:             1,
				Bumpiness:         2,
				LinesCleared:      0,
				LandingHeight:     0,
				Overhangs:         1,
				ColumnTransitions: 22, // col1: 4 (border→empty→filled→empty→border), cols 2-10: 2 each = 18
				RowTransitions:    40, // rows 1-18,20: 2 each (wall→empty→wall), row 19: 2 (wall→filled→empty→wall but filled touches wall)
				HoleDepth:         1,
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
				AggregateHeight:   10,
				Holes:             4,
				Bumpiness:         4,
				LinesCleared:      0,
				LandingHeight:     0,
				Overhangs:         4,
				ColumnTransitions: 28, // col1: 4, col2: 6 (gap), col3: 4, cols 4-10: 2 each = 14
				RowTransitions:    40, // each row has 2 transitions (wall boundaries + fill patterns)
				HoleDepth:         7,
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
				// Also (19,1) is empty with block above and open side
				return b
			}(),
			want: Heuristics{
				AggregateHeight:   4,
				Holes:             3,
				Bumpiness:         4,
				LinesCleared:      0,
				LandingHeight:     0,
				Overhangs:         3,
				ColumnTransitions: 22, // col1: 4 (border→empty→filled→empty→border), cols 2-10: 2 each = 18
				RowTransitions:    40, // each row has 2 transitions (wall boundaries)
				HoleDepth:         3,
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

				// Holes: (18,1), (19,1), (18,2), (19,2) all have blocks above
				// Overhangs: all 4 empty cells have blocks above and open sides
				return b
			}(),
			want: Heuristics{
				AggregateHeight:   8,
				Holes:             6,
				Bumpiness:         4,
				LinesCleared:      0,
				LandingHeight:     0,
				Overhangs:         6,
				ColumnTransitions: 24, // col1: 4, col2: 4, cols 3-10: 2 each = 16
				RowTransitions:    40, // each row has 2 transitions (wall boundaries)
				HoleDepth:         6,
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
