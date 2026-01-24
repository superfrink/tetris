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

	// AggregateHeight calculation
	for x := 0; x < gameCols; x++ {
		for y := 0; y < gameRows; y++ {
			if board[y][x] != 0 {
				h.AggregateHeight += (gameRows - y)
				break
			}
		}
	}

	// Holes calculation
	for x := 0; x < gameCols; x++ {
		foundBlock := false
		for y := 0; y < gameRows; y++ {
			if board[y][x] != 0 {
				foundBlock = true
			} else if foundBlock && board[y][x] == 0 {
				h.Holes++
			}
		}
	}
	// Bumpiness calculation
	columnHeights := make([]int, gameCols)
	for x := 0; x < gameCols; x++ {
		for y := 0; y < gameRows; y++ {
			if board[y][x] != 0 {
				columnHeights[x] = gameRows - y
				break
			}
		}
	}

	for x := 0; x < gameCols-1; x++ {
		diff := columnHeights[x] - columnHeights[x+1]
		if diff < 0 {
			diff = -diff
		}
		h.Bumpiness += diff
	}

	return h
}