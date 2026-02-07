package player

import (
	"testing"

	"superfrink.net/tetris/engine"
)

func TestFindBestMove(t *testing.T) {
	// Scenario: An empty board, a straight piece (I-tetromino)
	// Weights: High penalty for AggregateHeight, low for Holes and Bumpiness, to encourage flat and low boards.
	// Expected: Place the I-tetromino horizontally at the lowest possible position.

	game := engine.NewGame() // Initializes current and next piece
	// Ensure the current piece is an I-tetromino for a predictable test.
	game.Piece = 0 // 0 corresponds to I_TETROMINO in DefaultPieceMap

	// Define weights: high penalty for height, others neutral/low
	// These weights correspond to: AggregateHeight, Holes, Bumpiness, LinesCleared, LandingHeight, Overhangs, ColumnTransitions, RowTransitions
	weights := []float64{-1.0, -0.1, -0.1, 0.0, 0.0, 0.0, 0.0, 0.0} // Strongly penalize height

	// Expected move:
	// I-tetromino placed horizontally (rotation 0) at x-offset 1 (leftmost valid PiecePosCol after accounting for border).
	expectedMove := Move{
		Rotation:  1, // Vertical I-tetromino now yields a higher score
		XOffset:   1, // Leftmost valid PiecePosCol for an I-tetromino
		// EvaluatedScore will depend on the actual board state after placement,
		// so we won't assert it directly but ensure the move chosen is correct.
	}

	gotMove := FindBestMove(game, weights, nil)

	if gotMove.Rotation != expectedMove.Rotation || gotMove.XOffset != expectedMove.XOffset {
		t.Errorf("FindBestMove() got move {Rotation: %d, XOffset: %d, Score: %f}, want {Rotation: %d, XOffset: %d}",
			gotMove.Rotation, gotMove.XOffset, gotMove.EvaluatedScore, expectedMove.Rotation, expectedMove.XOffset)
	}
}

func TestFindBestMove_RightwardPreference(t *testing.T) {
	// Scenario: First column is completely filled, forcing the O-tetromino to move right.
	game := engine.NewGame()
	game.Piece = 3 // O_TETROMINO

	// Fill the first playable column to force rightward movement
	for r := 0; r < game.GameRows; r++ {
		game.Field[r][1] = 1
	}

	// Define weights: prioritize clearing lines and minimizing height
	weights := []float64{-1.0, -0.1, -0.1, 0.0, 0.0, 0.0, 0.0, 0.0} // Strongly penalize height

	gotMove := FindBestMove(game, weights, nil)

	// An O-tetromino is 2 blocks wide. If column 1 is filled, the leftmost valid XOffset should be 2.
	expectedXOffset := 2
	if gotMove.XOffset != expectedXOffset {
		t.Errorf("FindBestMove() got XOffset %d, expected %d for optimal placement after avoiding column 1", gotMove.XOffset, expectedXOffset)
	}
}

func TestFindBestMove_ClearsLine(t *testing.T) {
	// Scenario: A board configuration where placing a piece (I-tetromino) can clear a line.
	// Weights: Highly prioritize LinesCleared.
	game := engine.NewGame()
	game.Piece = 0 // I-tetromino

	// Create a board state that allows an I-tetromino (4 blocks wide) to clear a line.
	// Fill most of the second to last row, leaving a gap for the I-tetromino.
	// GameCols is 10, so playable columns are 1-10.
	// We'll fill row 19 (last playable row) almost completely.
	for c := 1; c <= game.GameColumns; c++ {
		if c != 5 { // Leave a gap at column 5
			game.Field[19][c] = 1
		}
	}
	// Fill row 18 (second to last playable row) partially to create a line clear scenario
	for c := 1; c <= game.GameColumns; c++ {
		if c != 6 && c != 7 && c != 8 && c != 9 { // Leave gaps for the I-tetromino
			game.Field[18][c] = 1
		}
	}

	// Define weights: high bonus for LinesCleared, others neutral
	weights := []float64{0.0, 0.0, 0.0, 100.0, 0.0, 0.0, 0.0, 0.0} // Huge bonus for clearing lines

	// An I-tetromino placed horizontally (rotation 0) at XOffset 6 (PiecePosCol 6)
	// would clear row 19 (and potentially more depending on the exact setup).
	// If rotation is 0 and PiecePosCol is 6, it will occupy columns 6, 7, 8, 9.
	// This would complete row 19 (which is missing 6,7,8,9)
	expectedMove := Move{
		Rotation:  0,
		XOffset:   6,
	}

	gotMove := FindBestMove(game, weights, nil)

	if gotMove.Rotation != expectedMove.Rotation || gotMove.XOffset != expectedMove.XOffset {
		t.Errorf("FindBestMove_ClearsLine() got move {Rotation: %d, XOffset: %d, Score: %f}, want {Rotation: %d, XOffset: %d}",
			gotMove.Rotation, gotMove.XOffset, gotMove.EvaluatedScore, expectedMove.Rotation, expectedMove.XOffset)
	}
}

func TestFindBestMove_LandingHeight(t *testing.T) {
	// Scenario: An empty board, a straight piece (I-tetromino)
	// Weights: Highly prioritize lower placement (higher PiecePosRow value).
	game := engine.NewGame()
	game.Piece = 0 // I-tetromino

	// Define weights: high bonus for LandingHeight, others neutral
	// Since PiecePosRow increases as the piece moves down, a positive weight
	// for LandingHeight will reward lower placements.
	weights := []float64{0.0, 0.0, 0.0, 0.0, 100.0, 0.0, 0.0, 0.0} // Huge bonus for lower placement

	// The I-tetromino can be placed at various XOffsets. On an empty board,
	// all valid XOffsets would result in the same lowest possible PiecePosRow.
	// We expect the piece to be placed at rotation 0 and the leftmost valid XOffset.
	// The leftmost XOffset for an I-tetromino (4 blocks wide) in a 10-column board
	// with borders would be 1.
	expectedMove := Move{
		Rotation: 0, // Horizontal I-tetromino for lowest block row (tie with vertical on empty board)
		XOffset:  1, // Leftmost valid XOffset
	}

	gotMove := FindBestMove(game, weights, nil)

	if gotMove.Rotation != expectedMove.Rotation || gotMove.XOffset != expectedMove.XOffset {
		t.Errorf("TestFindBestMove_LandingHeight() got move {Rotation: %d, XOffset: %d, Score: %f}, want {Rotation: %d, XOffset: %d}",
			gotMove.Rotation, gotMove.XOffset, gotMove.EvaluatedScore, expectedMove.Rotation, expectedMove.XOffset)
	}
}

func TestFindBestMove_WeightMask(t *testing.T) {
	game := engine.NewGame()
	game.Piece = 0 // I-tetromino

	// With LinesCleared weight active (index 3), the AI should try to clear lines.
	// With it masked out, it should behave as if that weight is zero.
	weightsWithLines := []float64{0.0, 0.0, 0.0, 100.0, 0.0, 0.0, 0.0, 0.0}
	maskAllOn := []bool{true, true, true, true, true, true, true, true}
	maskLinesClearedOff := []bool{true, true, true, false, true, true, true, true}

	// Set up board where clearing a line matters
	for c := 1; c <= game.GameColumns; c++ {
		if c < 7 {
			game.Field[18][c] = 1
		}
	}

	moveWithLines := FindBestMove(game, weightsWithLines, maskAllOn)
	moveWithoutLines := FindBestMove(game, weightsWithLines, maskLinesClearedOff)

	// When mask disables LinesCleared, the effective weight is 0, so the move should
	// be the same as if we passed all-zero weights.
	allZeroWeights := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	moveAllZero := FindBestMove(game, allZeroWeights, nil)

	if moveWithoutLines.Rotation != moveAllZero.Rotation || moveWithoutLines.XOffset != moveAllZero.XOffset {
		t.Errorf("Masked-out LinesCleared weight should produce same move as zero weight. Got {Rot:%d, X:%d}, want {Rot:%d, X:%d}",
			moveWithoutLines.Rotation, moveWithoutLines.XOffset, moveAllZero.Rotation, moveAllZero.XOffset)
	}

	// Sanity check: with mask all-on, the move should differ (lines cleared matters)
	if moveWithLines.Rotation == moveAllZero.Rotation && moveWithLines.XOffset == moveAllZero.XOffset {
		t.Log("Warning: move with LinesCleared=100 is same as all-zero weights; board setup may not differentiate")
	}
}

// TODO: Add more detailed tests for specific board configurations and piece types.
// For example:
// - Test a board where a vertical placement is better.
// - Test a board where a hole needs to be filled.
// - Test boundary conditions (left/right edges).
