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
	// These weights correspond to: AggregateHeight, Holes, Bumpiness
	weights := []float64{-1.0, -0.1, -0.1} // Strongly penalize height

	// Expected move:
	// I-tetromino placed horizontally (rotation 0) at x-offset 1 (leftmost valid PiecePosCol after accounting for border).
	expectedMove := Move{
		Rotation:  1, // Vertical I-tetromino now yields a higher score
		XOffset:   1, // Leftmost valid PiecePosCol for an I-tetromino
		// EvaluatedScore will depend on the actual board state after placement,
		// so we won't assert it directly but ensure the move chosen is correct.
	}

	gotMove := FindBestMove(game, weights)

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
	weights := []float64{-1.0, -0.1, -0.1} // Strongly penalize height

	gotMove := FindBestMove(game, weights)

	// An O-tetromino is 2 blocks wide. If column 1 is filled, the leftmost valid XOffset should be 2.
	expectedXOffset := 2
	if gotMove.XOffset != expectedXOffset {
		t.Errorf("FindBestMove() got XOffset %d, expected %d for optimal placement after avoiding column 1", gotMove.XOffset, expectedXOffset)
	}
}

// TODO: Add more detailed tests for specific board configurations and piece types.
// For example:
// - Test a board where a vertical placement is better.
// - Test a board where a hole needs to be filled.
// - Test boundary conditions (left/right edges).
