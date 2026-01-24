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
		Rotation:  0, // I-tetromino horizontal, default rotation
		XOffset:   1, // Leftmost valid PiecePosCol for a 4-block wide piece
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

	// An O-tetromino is 2 blocks wide. If column 0 is filled, the leftmost valid XOffset should be 1.
	// For an O-tetromino, if column 0 is blocked, it should ideally place itself starting at x=1.
	// This test simply asserts that it is not at x=0, which would be blocked.
	if gotMove.XOffset < 1 {
		t.Errorf("FindBestMove() got XOffset %d, expected XOffset >= 1 to avoid collision with filled column 0", gotMove.XOffset)
	}

	// More specific check: if only column 0 is filled, XOffset 1 is the most logical leftmost placement.
	// However, we are testing for a *preference* to move right, not just avoiding left collision.
	// Given the setup, it should definitely not pick XOffset 0.
	// Let's refine the expectation: given an empty board apart from column 0,
	// and an O-tetromino, it should settle at XOffset 1.
	expectedXOffset := 1
	if gotMove.XOffset != expectedXOffset {
		t.Errorf("FindBestMove() got XOffset %d, expected %d for optimal placement after avoiding column 0", gotMove.XOffset, expectedXOffset)
	}
}

// TODO: Add more detailed tests for specific board configurations and piece types.
// For example:
// - Test a board where a vertical placement is better.
// - Test a board where a hole needs to be filled.
// - Test boundary conditions (left/right edges).
