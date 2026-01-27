package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync" // Import the sync package
	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/pkg/evolution"
	"superfrink.net/tetris/pkg/player"
)

const (
	PopulationSize      = 50 // 100 // 250
	Generations         = 100
	ElitismFactor       = 0.85 // Top N% individuals are carried over
	MutationRate        = 0.10 // 0.20 // 0.05
	GamesPerIndividual  = 100 // Number of games each individual plays per generation
	CheckpointFrequency = 5 // Save checkpoint every N generations

	GenomeLength = 6 // AggregateHeight, Holes, Bumpiness, LinesCleared, LandingHeight, Overhangs
	CheckpointFileName = "population_checkpoint.json"
)

// loadCheckpoint attempts to load a Population from a checkpoint file.
// Returns the loaded Population and the starting generation, or nil and an error.
func loadCheckpoint() (*evolution.Population, error) {
	data, err := os.ReadFile(CheckpointFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Checkpoint file %s not found. Starting new training.\n", CheckpointFileName)
			return nil, nil // No error, just no checkpoint
		}
		return nil, fmt.Errorf("error reading checkpoint file: %w", err)
	}

	var pop evolution.Population
	if err := json.Unmarshal(data, &pop); err != nil {
		return nil, fmt.Errorf("error unmarshalling checkpoint data: %w", err)
	}

	fmt.Printf("Loaded checkpoint from %s. Resuming from Generation %d.\n", CheckpointFileName, pop.Generation)
	return &pop, nil
}

// saveCheckpoint marshals and writes the current population state to a checkpoint file.
func saveCheckpoint(pop *evolution.Population) error {
	data, err := json.MarshalIndent(pop, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling population for checkpoint: %w", err)
	}

	if err := os.WriteFile(CheckpointFileName, data, 0644); err != nil {
		return fmt.Errorf("error writing checkpoint file %s: %w", CheckpointFileName, err)
	}

	fmt.Printf("Checkpoint saved to %s for Generation %d.\n", CheckpointFileName, pop.Generation)
	return nil
}

func main() {
	var pop *evolution.Population
	loadedPop, err := loadCheckpoint()
	if err != nil {
		fmt.Printf("Fatal error loading checkpoint: %v\n", err)
		os.Exit(1)
	}

	if loadedPop != nil {
		pop = loadedPop
	} else {
		fmt.Println("Initializing new population...")
		pop = evolution.NewPopulation(PopulationSize, GenomeLength)
	}

	fmt.Printf("Playing %d games per individual per generation.\n", GamesPerIndividual)

	for gen := pop.Generation; gen < Generations; gen++ {
		fmt.Printf("Generation %d: ", gen)

		var wg sync.WaitGroup
		for i := range pop.Individuals {
			wg.Add(1)
			go func(individual *evolution.Individual) {
				defer wg.Done()
				totalIndividualFitness := 0

				for gameNum := 0; gameNum < GamesPerIndividual; gameNum++ {
					game := engine.NewGame() // Start a new headless game for each individual
					game.State = engine.StateRunning // Ensure game is running

					for game.State != engine.StateGameOver {
						// Find the best move based on the current game state and individual's genome
						move := player.FindBestMove(game, individual.Genome)

						// Apply the rotation
						for i := 0; i < move.Rotation; i++ {
							game.RotatePiece()
						}
						// Set the x-position
						game.PiecePosCol = move.XOffset

						// Drop the piece until it lands
						for {
							// Remember the current piece ID and row before stepping
							currentPiece := game.Piece
							currentRow := game.PiecePosRow

							game.Step(engine.PlayInputDrop)

							// If the piece ID changed, a new piece has spawned, so the old one landed.
							// Also break if the row didn't change, which can happen at the top of the board
							// in a game over state before the piece ID changes.
							if game.Piece != currentPiece || game.PiecePosRow == currentRow {
								break
							}
						}

						// Check for game over state after the piece has landed
						if game.State == engine.StateGameOver {
							break
						}
					}
					totalIndividualFitness += (game.ScoreLineCount * 100) + (game.ScorePieceCount * 1)
				}
				individual.Fitness = totalIndividualFitness / GamesPerIndividual
			}(pop.Individuals[i])
		}
		wg.Wait() // Wait for all individuals in the current generation to be evaluated

		// Calculate and print performance summary
		bestFitness := 0
		totalFitness := 0
		for _, individual := range pop.Individuals {
			if individual.Fitness > bestFitness {
				bestFitness = individual.Fitness
			}
			totalFitness += individual.Fitness
		}
		averageFitness := float64(totalFitness) / float64(len(pop.Individuals))

		fmt.Printf("  Best Fitness: %d, Average Fitness: %.2f\n", bestFitness, averageFitness)

		// Evolve population
		pop = evolution.Evolve(pop, ElitismFactor, MutationRate)
		if (pop.Generation)%CheckpointFrequency == 0 { // Check after evaluation and evolution for current gen
			if err := saveCheckpoint(pop); err != nil {
				fmt.Printf("Error saving checkpoint: %v\n", err)
				// Optionally, you might want to exit or retry here
			}
		}
	}

	// After training, find the best individual
	pop.Sort()
	bestIndividual := pop.Individuals[0]

	// Marshal the best individual's genome to JSON
	genomeJSON, err := json.MarshalIndent(bestIndividual.Genome, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling best genome to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write the JSON to a file
	const TrainedModelFileName = "trained_model.json"
	if err := os.WriteFile(TrainedModelFileName, genomeJSON, 0644); err != nil {
		fmt.Printf("Error writing trained model to file %s: %v\n", TrainedModelFileName, err)
		os.Exit(1)
	}

	fmt.Printf("Best individual's genome saved to %s\n", TrainedModelFileName)
}
