package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/pkg/evolution"
	"superfrink.net/tetris/pkg/player"
)

type TrainedModel struct {
	Weights []float64 `json:"weights"`
	Mask    string    `json:"mask"`
}

const (
	PopulationSize      = 30 // 100 // 250
	Generations         = 200
	ElitismFactor       = 0.85 // Top N% individuals are carried over
	MutationRate        = 0.15 // 0.20 // 0.05
	GamesPerIndividual  = 20 // Number of games each individual plays per generation
	CheckpointFrequency = 10 // Save checkpoint every N generations

	GenomeLength = 8 // AggregateHeight, Holes, Bumpiness, LinesCleared, LandingHeight, Overhangs, ColumnTransitions, RowTransitions
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

func parseMask(maskStr string) ([]bool, error) {
	if len(maskStr) != GenomeLength {
		return nil, fmt.Errorf("mask must be exactly %d characters, got %d", GenomeLength, len(maskStr))
	}
	mask := make([]bool, GenomeLength)
	for i, ch := range maskStr {
		switch ch {
		case '1':
			mask[i] = true
		case '0':
			mask[i] = false
		default:
			return nil, fmt.Errorf("mask must contain only '0' or '1', got '%c' at position %d", ch, i)
		}
	}
	return mask, nil
}

func applyMaskToPopulation(pop *evolution.Population, weightMask []bool) {
	for _, ind := range pop.Individuals {
		for i, enabled := range weightMask {
			if !enabled && i < len(ind.Genome) {
				ind.Genome[i] = 0
			}
		}
	}
}

func main() {
	maskFlag := flag.String("mask", "11111111", "Weight mask: 8 characters of 0/1 indicating which weights are active.\n"+
		"Position 0 is the leftmost character.\n"+
		"Positions: 0=AggregateHeight, 1=Holes, 2=Bumpiness, 3=LinesCleared,\n"+
		"           4=LandingHeight, 5=Overhangs, 6=ColumnTransitions, 7=RowTransitions")
	flag.Parse()

	weightMask, err := parseMask(*maskFlag)
	if err != nil {
		fmt.Printf("Invalid mask: %v\n", err)
		os.Exit(1)
	}

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
	applyMaskToPopulation(pop, weightMask)

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
						move := player.FindBestMove(game, individual.Genome, weightMask)

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
		applyMaskToPopulation(pop, weightMask)
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

	model := TrainedModel{
		Weights: bestIndividual.Genome,
		Mask:    *maskFlag,
	}
	modelJSON, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling trained model to JSON: %v\n", err)
		os.Exit(1)
	}

	const TrainedModelFileName = "trained_model.json"
	if err := os.WriteFile(TrainedModelFileName, modelJSON, 0644); err != nil {
		fmt.Printf("Error writing trained model to file %s: %v\n", TrainedModelFileName, err)
		os.Exit(1)
	}

	fmt.Printf("Best individual's genome saved to %s\n", TrainedModelFileName)
}
