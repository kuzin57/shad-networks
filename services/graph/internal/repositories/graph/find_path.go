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
		queries.FindNearestDijkstra,
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

	paths, err := r.getPathsFromResult(result)
	switch {
	case err != nil:
		return nil, err
	case len(paths) == 0:
		return nil, nil
	default:
		return paths[0], nil
	}
}
