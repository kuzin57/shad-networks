package graph

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/internal/entities"
	"github.com/kuzin57/shad-networks/internal/repositories/graph/queries"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) GetGraph(ctx context.Context, graphID string) (graph entities.Graph, err error) {
	result, err := neo4j.ExecuteQuery(
		ctx,
		r.driver.Get(),
		queries.GetNodesWithRels,
		map[string]any{
			idParam: graphID,
		},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return entities.Graph{}, fmt.Errorf("get nodes with rels: %w", err)
	}

	verticesNimber := r.getVerticesNumber(result.Records)
	graph.AdjencyMaxtrix = make([][][]int, verticesNimber)
	for i := range verticesNimber {
		graph.AdjencyMaxtrix[i] = make([][]int, verticesNimber)
	}

	for _, record := range result.Records {
		if record == nil {
			r.log.Warn("skipping nil record from db")

			continue
		}

		var (
			source = r.getIntValue(record, sourceParam)
			target = r.getIntValue(record, targetParam)
			weight = r.getIntValue(record, weightParam)
		)

		if source < 0 || target < 0 || weight < 0 {
			return entities.Graph{}, fmt.Errorf("invalid records, source: %d, target: %d, weight: %d", source, target, weight)
		}

		graph.AdjencyMaxtrix[source][target] = append(graph.AdjencyMaxtrix[source][target], weight)
	}

	return
}

func (r *Repository) getVerticesNumber(records []*neo4j.Record) (verticesNumber int) {
	for _, record := range records {
		if record == nil {
			continue
		}

		source := r.getIntValue(record, sourceParam)
		target := r.getIntValue(record, targetParam)

		if source < 0 || target < 0 {
			return 0
		}

		verticesNumber = max(verticesNumber, source, target)
	}

	return verticesNumber + 1
}

func (r *Repository) getIntValue(record *neo4j.Record, key string) int {
	number, ok := record.Get(key)
	if !ok {
		r.log.Sugar().Warnf("no such key: %s", key)

		return -1
	}

	nodeNumber, _ := number.(float64)

	return int(nodeNumber)
}
