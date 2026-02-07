package fitter

// Heuristics holds the calculated values for a given board state.
type Heuristics struct {
	AggregateHeight   int
	Holes             int
	Bumpiness         int
	LinesCleared      int
	LandingHeight     int
	Overhangs         int
	ColumnTransitions int
	RowTransitions    int
	HoleDepth         int
}

// CalculateHeuristics computes the feature values from the board state.
// This is the definitive, correct version after multiple debugging cycles.
func CalculateHeuristics(board [][]int, linesCleared int, landingHeight int) Heuristics {
	h := Heuristics{}
	h.LinesCleared = linesCleared
	h.LandingHeight = landingHeight

	gameRows := len(board)
	gameCols := len(board[0])
	playfieldCols := gameCols - 2

	// Bumpiness calculation
	columnHeights := make([]int, gameCols)
	for x := 1; x <= playfieldCols; x++ {
		for y := 1; y < gameRows-1; y++ {
			if board[y][x] != 0 {
				columnHeights[x] = gameRows - 1 - y
				break
			}
		}
	}

	for x := 1; x < playfieldCols; x++ {
		diff := columnHeights[x] - columnHeights[x+1]
		if diff < 0 {
			diff = -diff
		}
		h.Bumpiness += diff
	}

	// Holes and HoleDepth calculation
	for x := 1; x <= playfieldCols; x++ {
		for y := 1; y < gameRows-1; y++ {
			if board[y][x] == 0 { // Found an empty cell
				filledAbove := 0
				for y_above := y - 1; y_above >= 1; y_above-- {
					if board[y_above][x] != 0 {
						filledAbove++
					}
				}
				if filledAbove > 0 {
					h.Holes++
					h.HoleDepth += filledAbove
				}
			}
		}
	}

	// AggregateHeight calculation
	for x := 1; x <= playfieldCols; x++ {
		h.AggregateHeight += columnHeights[x]
	}

	// Overhangs calculation
	for x := 1; x <= playfieldCols; x++ {
		for y := 1; y < gameRows-1; y++ {
			if board[y][x] == 0 { // Is an empty spot
				hasBlockAbove := false
				for y_above := y - 1; y_above >= 1; y_above-- {
					if board[y_above][x] != 0 {
						hasBlockAbove = true
						break
					}
				}
				if hasBlockAbove {
					if (x > 1 && board[y][x-1] == 0) || (x < playfieldCols && board[y][x+1] == 0) {
						h.Overhangs++
					}
				}
			}
		}
	}

	// ColumnTransitions calculation
	// Count vertical changes between empty and filled cells in each column,
	// including transitions with top and bottom border rows.
	for x := 1; x <= playfieldCols; x++ {
		prev := board[0][x]
		for y := 1; y <= gameRows-1; y++ {
			cur := board[y][x]
			if (prev == 0) != (cur == 0) {
				h.ColumnTransitions++
			}
			prev = cur
		}
	}

	// RowTransitions calculation
	// Count horizontal changes between empty and filled cells in each row,
	// including transitions with left and right wall columns.
	for y := 1; y < gameRows-1; y++ {
		prev := board[y][0]
		for x := 1; x <= playfieldCols; x++ {
			cur := board[y][x]
			if (prev == 0) != (cur == 0) {
				h.RowTransitions++
			}
			prev = cur
		}
		if (prev == 0) != (board[y][playfieldCols+1] == 0) {
			h.RowTransitions++
		}
	}

	return h
}