package fitter

import (
	"reflect"
	"testing"
)

func TestCalculateHeuristics(t *testing.T) {
	tests := []struct {
		name  string
		board [][]int
		want  Heuristics
	}{
		{
			name:  "Empty board",
			board: make([][]int, 20),
			want: Heuristics{
				AggregateHeight: 0,
				Holes:           0,
				Bumpiness:       0,
			},
		},
		{
			name: "Single block at bottom-left",
			board: func() [][]int {
				b := make([][]int, 20)
				for i := range b {
					b[i] = make([]int, 10)
				}
				b[19][0] = 1
				return b
			}(),
			want: Heuristics{
				AggregateHeight: 1, // 20 - 19
				Holes:           0,
				Bumpiness:       1, // difference between col 0 (height 1) and col 1 (height 0)
			},
		},
		{
			name: "Board with holes and bumpiness",
			board: func() [][]int {
				b := make([][]int, 20)
				for i := range b {
					b[i] = make([]int, 10)
				}
				// Column 0: height 3
				b[19][0] = 1
				b[18][0] = 1
				b[17][0] = 1
				// Column 1: height 2, 1 hole
				b[19][1] = 1
				b[17][1] = 1 // hole at [18][1]
				// Column 2: height 1
				b[19][2] = 1
				// Column 3: height 0

				return b
			}(),
			want: Heuristics{
				AggregateHeight: 3 + 3 + 1, // (20-17) for col0, (20-17) for col1, (20-19) for col2
				Holes:           1,
				Bumpiness:       (3-2) + (2-1) + (1-0), // (3-2)+(2-1)+(1-0) = 1+1+1 = 3
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateHeuristics(tt.board)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateHeuristics() = %v, want %v", got, tt.want)
			}
		})
	}
}
