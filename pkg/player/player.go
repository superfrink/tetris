package player

import (
	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/pkg/fitter"
)

// Move represents a potential move for a Tetris piece.
type Move struct {
	Rotation       int     // How many times the piece is rotated (0-3)
	XOffset        int     // The x-position of the left-most block of the piece
	EvaluatedScore float64 // The score of the board after this move, based on heuristics
}

// FindBestMove determines the optimal move for the current piece based on the AI's weights.
func FindBestMove(g *engine.Game, weights []float64) Move {
	bestMove := Move{EvaluatedScore: -999999999} // Initialize with a very low score

	for rot := 0; rot < 4; rot++ { // rot is the desired final rotation count (0 to 3)
		rotatedGame := g.CopyOfState()
		rotatedGame.PieceRotation = rot

		// If after desired rotations, the piece is in a colliding state, skip this rotation
		if rotatedGame.CheckCollision(rotatedGame.Piece, rotatedGame.PieceRotation, rotatedGame.PiecePosRow, rotatedGame.PiecePosCol) {
			continue
		}

		// Simulate moving left as much as possible
		tempGame := rotatedGame.CopyOfState() // Use tempGame for horizontal movement simulation
		for {
			if tempGame.CheckCollision(tempGame.Piece, tempGame.PieceRotation, tempGame.PiecePosRow, tempGame.PiecePosCol-1) {
				break // Cannot move further left
			}
			tempGame.MoveLeftPiece()
		}

		// Now iterate from the leftmost valid position to the rightmost
		for {
			movedGame := tempGame.CopyOfState() // Copy for this x-position evaluation

			// If the piece is colliding in this x-position, break the loop
			if movedGame.CheckCollision(movedGame.Piece, movedGame.PieceRotation, movedGame.PiecePosRow, movedGame.PiecePosCol) {
				break
			}

			// Drop the piece
			finalGame := movedGame.CopyOfState()
			for {
				if !finalGame.LowerPiece() {
					break // Piece can't move down further
				}
			}
			finalGame.PlacePiece() // Place the piece after it has landed in simulation
			linesCleared := finalGame.ClearCompletedRows()

			// Calculate landingHeight as the row of the lowest placed block in the piece
			maxIRelativeOffset := 0
			pieceMap := finalGame.PieceMap[finalGame.Piece][finalGame.PieceRotation]
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if pieceMap[i][j] != 0 {
						if i > maxIRelativeOffset {
							maxIRelativeOffset = i
						}
					}
				}
			}
			landingHeight := finalGame.PiecePosRow + maxIRelativeOffset

			// Evaluate this finalGame state
			heuristics := fitter.CalculateHeuristics(finalGame.Field, linesCleared, landingHeight)
			currentScore := -(weights[0]*float64(heuristics.AggregateHeight) +
				weights[1]*float64(heuristics.Holes) +
				weights[2]*float64(heuristics.Bumpiness)) +
				weights[3]*float64(heuristics.LinesCleared) +
				weights[4]*float64(heuristics.LandingHeight) +
				weights[5]*float64(heuristics.Overhangs)

			if currentScore > bestMove.EvaluatedScore {
				bestMove.EvaluatedScore = currentScore
				bestMove.Rotation = rot // Store the desired rotation
				bestMove.XOffset = movedGame.PiecePosCol // Store the PiecePosCol of the evaluated move
			}

			// Try to move one step right for the next iteration
			if tempGame.CheckCollision(tempGame.Piece, tempGame.PieceRotation, tempGame.PiecePosRow, tempGame.PiecePosCol+1) {
				break // Cannot move further right
			}
			tempGame.MoveRightPiece()
		}
	}

	return bestMove
}