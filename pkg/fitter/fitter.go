package fitter

// Heuristics holds the calculated values for a given board state.
type Heuristics struct {
	AggregateHeight int
	Holes           int
	Bumpiness       int
}

// CalculateHeuristics computes the feature values from the board state.
func CalculateHeuristics(board [][]int) Heuristics {
	h := Heuristics{}

	gameRows := len(board)
	gameCols := len(board[0])
	playfieldCols := gameCols - 2

	// AggregateHeight calculation
	for x := 1; x <= playfieldCols; x++ {
		for y := 1; y < gameRows-1; y++ {
			if board[y][x] != 0 {
				h.AggregateHeight += (gameRows - 1 - y)
				break
			}
		}
	}

	// Holes calculation
	for x := 1; x <= playfieldCols; x++ {
		for y := 1; y < gameRows-1; y++ {
			if board[y][x] == 0 { // Found an empty cell
				// Check for a block above it in the same column
				for y_above := y - 1; y_above >= 1; y_above-- {
					if board[y_above][x] != 0 {
						h.Holes++
						break // Found a block above, count it as a hole and move to the next empty cell
					}
				}
			}
		}
	}
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

	return h
}