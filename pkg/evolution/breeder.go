package evolution

import (
	"log/slog"
	"math/rand"
	"sort" // Import the sort package
)

// Individual represents a single AI with its Tetris playing strategy (genome) and its fitness score.
type Individual struct {
	Genome  []float64
	Fitness int
}

// Population represents a collection of individuals for a given generation.
type Population struct {
	Individuals []*Individual
	Generation  int
}

// NewPopulation creates a new population with randomly initialized genomes.
func NewPopulation(size int, genomeLength int) *Population {
	individuals := make([]*Individual, size)
	for i := 0; i < size; i++ {
		genome := make([]float64, genomeLength)
		for j := 0; j < genomeLength; j++ {
			// Initialize weights in a reasonable range, e.g., -1.0 to 1.0
			genome[j] = rand.Float64()
		}
		individuals[i] = &Individual{
			Genome:  genome,
			Fitness: 0, // Fitness will be evaluated later
		}
	}
	return &Population{
		Individuals: individuals,
		Generation:  0,
	}
}

// Sort sorts the individuals in the population by fitness in descending order.
func (p *Population) Sort() {
	sort.Slice(p.Individuals, func(i, j int) bool {
		return p.Individuals[i].Fitness > p.Individuals[j].Fitness
	})
}

// Evolve creates the next generation of the population through elitism, crossover, and mutation.
func Evolve(pop *Population, elitismFactor float64, mutationRate float64) *Population {
	slog.Debug("Evolve", "populationSize", len(pop.Individuals), "elitismFactor", elitismFactor)

	// Sort the current population by fitness (descending)
	pop.Sort()

	newIndividuals := make([]*Individual, 0, len(pop.Individuals))

	// Elitism: carry over the top individuals directly
	numElite := int(float64(len(pop.Individuals)) * elitismFactor)
	for i := 0; i < numElite; i++ {
		newIndividuals = append(newIndividuals, pop.Individuals[i])
	}

	// Crossover and Mutation: fill the rest of the new generation
	for len(newIndividuals) < len(pop.Individuals) {
		// Select two parents (e.g., from the elite or using a selection mechanism)
		parent1 := pop.Individuals[rand.Intn(numElite)] // Simple selection from elite
		parent2 := pop.Individuals[rand.Intn(numElite)]

		// Crossover: create offspring genome
		offspringGenome := make([]float64, len(parent1.Genome))
		for i := range offspringGenome {
			if rand.Float64() < 0.5 { // 50% chance to inherit from parent1 or parent2
				offspringGenome[i] = parent1.Genome[i]
			} else {
				offspringGenome[i] = parent2.Genome[i]
			}
		}

		// Mutation: introduce small random changes
		for i := range offspringGenome {
			if rand.Float64() < mutationRate {
				offspringGenome[i] += (rand.Float64()*2 - 1) * 0.2 // Small random change
				if offspringGenome[i] < 0 {
					offspringGenome[i] = 0
				}
			}
		}

		newIndividuals = append(newIndividuals, &Individual{
			Genome:  offspringGenome,
			Fitness: 0, // Fitness will be evaluated in the next generation
		})
	}

	return &Population{
		Individuals: newIndividuals,
		Generation:  pop.Generation + 1,
	}
}