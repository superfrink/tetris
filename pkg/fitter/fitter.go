package fitter

// Heuristics holds the calculated values for a given board state.
type Heuristics struct {
	AggregateHeight int
	Holes           int
	Bumpiness       int
	LinesCleared    int
	LandingHeight   int
	Overhangs       int
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

				// Check for blocks on both sides within the playfield
				// `x > 1` ensures `x-1` is not the left border
				// `x < playfieldCols` ensures `x+1` is not the right border
				hasBlockLeft := (x > 1 && board[y][x-1] != 0)
				hasBlockRight := (x <= playfieldCols && board[y][x+1] != 0) // x <= playfieldCols ensures x+1 is not outside the playfield + border

				if hasBlockAbove && hasBlockLeft && hasBlockRight {
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
				if y > 1 && board[y-1][x] != 0 { // Has a filled spot directly above it (y>1 to avoid border)
					// Check for an open spot beside it (left or right)
					// x > 1 ensures we are not checking the left border
					// x < gameCols-1 ensures we are not checking the right border
					if (x > 1 && board[y][x-1] == 0) || (x < gameCols-1 && board[y][x+1] == 0) {
						h.Overhangs++
					}
				}
			}
		}
	}

	return h
}