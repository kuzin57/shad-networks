package graphgenerator

import (
	"math/rand"
	"slices"

	"github.com/kuzin57/shad-networks/internal/entities"
	"go.uber.org/zap"
)

const (
	maxRetries = 50
)

type Generator struct {
	log *zap.Logger
}

func NewGenerator(log *zap.Logger) *Generator {
	generator := &Generator{
		log: log,
	}

	return generator
}

func (g *Generator) Generate(params entities.GraphParams) entities.Graph {
	adjMatrix := make([][][]int, 0, params.VerticesCount)
	for range params.VerticesCount {
		adjMatrix = append(adjMatrix, make([][]int, params.VerticesCount))
	}

	slices.Sort(params.Degrees)

	g.log.Info("generating graph", zap.Int("vertices count", int(params.VerticesCount)))

	maxDegree := int(params.Degrees[len(params.Degrees)-1])

	for i := range adjMatrix {
		deg := g.countDeg(adjMatrix, i)
		if deg >= maxDegree {
			continue
		}

		newDeg := g.chooseDeg(params.Degrees, deg)

		for retry := 0; deg < newDeg && retry < maxRetries; {
			neighbor := rand.Intn(len(adjMatrix))

			if len(adjMatrix[i][neighbor]) > 0 {
				retry++

				continue
			}

			weightInd := rand.Intn(len(params.Weights))
			adjMatrix[i][neighbor] = []int{int(params.Weights[weightInd])}
			adjMatrix[neighbor][i] = []int{int(params.Weights[weightInd])}
			retry = 0

			deg++
		}
	}

	g.log.Info("graph successfully generated")

	return entities.Graph{
		AdjencyMaxtrix: adjMatrix,
	}
}

func (g *Generator) countDeg(adjMatrix [][][]int, vertex int) int {
	var cnt int

	for _, weights := range adjMatrix[vertex] {
		cnt += len(weights)
	}

	return cnt
}

func (g *Generator) chooseDeg(degrees []uint32, deg int) int {
	index, _ := slices.BinarySearch(degrees, uint32(deg))

	return int(degrees[index+rand.Intn(len(degrees)-index)])
}
