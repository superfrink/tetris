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

	// Holes calculation
	for x := 1; x <= playfieldCols; x++ {
		for y := 1; y < gameRows-1; y++ {
			if board[y][x] == 0 { // Found an empty cell
				hasBlockAbove := false
				for y_above := y - 1; y_above >= 1; y_above-- {
					if board[y_above][x] != 0 {
						hasBlockAbove = true
						break
					}
				}
				if hasBlockAbove {
					h.Holes++
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
	// Count vertical changes between empty and filled cells in each column
	for x := 1; x <= playfieldCols; x++ {
		for y := 2; y < gameRows-1; y++ {
			previousCell := board[y-1][x]
			currentCell := board[y][x]
			if (previousCell == 0 && currentCell != 0) || (previousCell != 0 && currentCell == 0) {
				h.ColumnTransitions++
			}
		}
	}

	// RowTransitions calculation
	// Count horizontal changes between empty and filled cells in each row
	for y := 1; y < gameRows-1; y++ {
		for x := 2; x <= playfieldCols; x++ {
			previousCell := board[y][x-1]
			currentCell := board[y][x]
			if (previousCell == 0 && currentCell != 0) || (previousCell != 0 && currentCell == 0) {
				h.RowTransitions++
			}
		}
	}

	return h
}