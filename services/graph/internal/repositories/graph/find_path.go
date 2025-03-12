package graph

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"github.com/kuzin57/shad-networks/services/graph/internal/repositories/graph/queries"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) FindPath(ctx context.Context, graphID string, source, target int) (entities.Path, error) {
	result, err := neo4j.ExecuteQuery(
		ctx,
		r.driver.Unwrap(),
		queries.FindNearestBetweenTwoDijkstra,
		map[string]any{
			idParam:     graphID,
			sourceParam: source,
			targetParam: target,
		},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return nil, fmt.Errorf("find path via dijkstra: %w", err)
	}

	if len(result.Records) == 0 || len(result.Records) > 1 {
		return nil, fmt.Errorf("invalid records len: %d", len(result.Records))
	}

	record := result.Records[0]
	if record == nil {
		return nil, fmt.Errorf("nil record")
	}

	path, ok := record.Get(pathParam)
	if !ok {
		return nil, fmt.Errorf("path param not found")
	}

	costs, ok := record.Get(costsParam)
	if !ok {
		return nil, fmt.Errorf("costs param not found")
	}

	var (
		pathValue, _  = path.([]float64)
		costsValue, _ = costs.([]float64)
		resultPath    = make(entities.Path, 0, len(pathValue))
	)

	for i := range pathValue {
		if i == 0 {
			resultPath = append(resultPath, entities.PathPart{
				Vertex: int(pathValue[i]),
				Weight: 0},
			)

			continue
		}

		resultPath = append(resultPath, entities.PathPart{
			Vertex: int(pathValue[i]),
			Weight: int(costsValue[i]) - int(costsValue[i-1]),
		})
	}

	return resultPath, nil
}
