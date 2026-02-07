package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"sync"
	"superfrink.net/tetris/engine"
	"superfrink.net/tetris/pkg/evolution"
	"superfrink.net/tetris/pkg/player"
)

const (
	DefaultPopulationSize     = 30
	DefaultGenerations        = 200
	DefaultElitismFactor      = 0.85
	DefaultMutationRate       = 0.15
	DefaultGamesPerIndividual = 20
	GenomeLength              = 9
)

var weightLabels = []string{
	"AggregateHeight",
	"Holes",
	"Bumpiness",
	"LinesCleared",
	"LandingHeight",
	"Overhangs",
	"ColumnTransitions",
	"RowTransitions",
	"HoleDepth",
}

type AblationResult struct {
	Mask        string    `json:"mask"`
	Label       string    `json:"label"`
	BestFitness int       `json:"best_fitness"`
	AvgFitness  float64   `json:"avg_fitness"`
	BestWeights []float64 `json:"weights"`
}

type TrainedModel struct {
	Weights []float64 `json:"weights"`
	Mask    string    `json:"mask"`
}

type maskEntry struct {
	Mask  string
	Label string
}

func buildMasks() []maskEntry {
	baseline := "111111111"
	entries := []maskEntry{{Mask: baseline, Label: "baseline"}}
	for i := 0; i < GenomeLength; i++ {
		m := []byte(baseline)
		m[i] = '0'
		entries = append(entries, maskEntry{
			Mask:  string(m),
			Label: "-" + weightLabels[i],
		})
	}
	return entries
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

func runTraining(mask string, label string, index int, total int, popSize int, generations int, gamesPerIndividual int, elitism float64, mutationRate float64) AblationResult {
	weightMask, err := parseMask(mask)
	if err != nil {
		fmt.Printf("Invalid mask %s: %v\n", mask, err)
		return AblationResult{Mask: mask, Label: label}
	}

	fmt.Printf("\n=== Training mask %s (%s) [%d/%d] ===\n", mask, label, index, total)

	pop := evolution.NewPopulation(popSize, GenomeLength)
	applyMaskToPopulation(pop, weightMask)

	var lastBestFitness int
	var lastAvgFitness float64

	for gen := 0; gen < generations; gen++ {
		fmt.Printf("Generation %d: ", gen)

		var wg sync.WaitGroup
		for i := range pop.Individuals {
			wg.Add(1)
			go func(individual *evolution.Individual) {
				defer wg.Done()
				totalIndividualFitness := 0

				for gameNum := 0; gameNum < gamesPerIndividual; gameNum++ {
					game := engine.NewGame()
					game.State = engine.StateRunning

					for game.State != engine.StateGameOver {
						move := player.FindBestMove(game, individual.Genome, weightMask)

						for i := 0; i < move.Rotation; i++ {
							game.RotatePiece()
						}
						game.PiecePosCol = move.XOffset

						for {
							currentPiece := game.Piece
							currentRow := game.PiecePosRow
							game.Step(engine.PlayInputDrop)
							if game.Piece != currentPiece || game.PiecePosRow == currentRow {
								break
							}
						}

						if game.State == engine.StateGameOver {
							break
						}
					}
					totalIndividualFitness += (game.ScoreLineCount * 100) + (game.ScorePieceCount * 1)
				}
				individual.Fitness = totalIndividualFitness / gamesPerIndividual
			}(pop.Individuals[i])
		}
		wg.Wait()

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

		lastBestFitness = bestFitness
		lastAvgFitness = averageFitness

		pop = evolution.Evolve(pop, elitism, mutationRate)
		applyMaskToPopulation(pop, weightMask)
	}

	pop.Sort()
	bestIndividual := pop.Individuals[0]

	return AblationResult{
		Mask:        mask,
		Label:       label,
		BestFitness: lastBestFitness,
		AvgFitness:  lastAvgFitness,
		BestWeights: bestIndividual.Genome,
	}
}

func main() {
	quick := flag.Bool("quick", false, "Use reduced parameters (50 gen, 10 pop, 10 games) for fast testing")
	genFlag := flag.Int("generations", 0, "Override generation count")
	popFlag := flag.Int("population", 0, "Override population size")
	gamesFlag := flag.Int("games", 0, "Override games per individual")
	flag.Parse()

	popSize := DefaultPopulationSize
	generations := DefaultGenerations
	gamesPerIndividual := DefaultGamesPerIndividual

	if *quick {
		popSize = 10
		generations = 50
		gamesPerIndividual = 10
	}

	if *genFlag > 0 {
		generations = *genFlag
	}
	if *popFlag > 0 {
		popSize = *popFlag
	}
	if *gamesFlag > 0 {
		gamesPerIndividual = *gamesFlag
	}

	fmt.Printf("Parameters: population=%d, generations=%d, games=%d, elitism=%.2f, mutation=%.2f\n",
		popSize, generations, gamesPerIndividual, DefaultElitismFactor, DefaultMutationRate)

	masks := buildMasks()
	results := make([]AblationResult, 0, len(masks))

	for i, entry := range masks {
		result := runTraining(entry.Mask, entry.Label, i+1, len(masks),
			popSize, generations, gamesPerIndividual, DefaultElitismFactor, DefaultMutationRate)
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].BestFitness > results[j].BestFitness
	})

	fmt.Printf("\nABLATION RESULTS (sorted by Best Fitness)\n")
	fmt.Printf("%-10s | %-20s | %12s | %11s\n", "Mask", "Label", "Best Fitness", "Avg Fitness")
	fmt.Printf("%-10s-|-%-20s-|-%-12s-|-%-11s\n", "----------", "--------------------", "------------", "-----------")
	for _, r := range results {
		fmt.Printf("%-10s | %-20s | %12d | %11.2f\n", r.Mask, r.Label, r.BestFitness, r.AvgFitness)
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling results: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile("ablation_results.json", data, 0644); err != nil {
		fmt.Printf("Error writing ablation_results.json: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nResults written to ablation_results.json")

	for _, r := range results {
		model := TrainedModel{Weights: r.BestWeights, Mask: r.Mask}
		modelJSON, err := json.MarshalIndent(model, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling model for mask %s: %v\n", r.Mask, err)
			continue
		}
		filename := fmt.Sprintf("trained_model-%s.json", r.Mask)
		if err := os.WriteFile(filename, modelJSON, 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", filename, err)
			continue
		}
	}
	fmt.Printf("Saved %d trained models as trained_model-MASK.json\n", len(results))
}
