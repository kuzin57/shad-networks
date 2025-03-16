package graph

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"github.com/kuzin57/shad-networks/services/graph/internal/repositories/graph/queries"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) FindPaths(ctx context.Context, graphID string, source, target, amount int) ([]entities.Path, error) {
	result, err := neo4j.ExecuteQuery(
		ctx,
		r.driver.Unwrap(),
		queries.FindKNearestYens,
		map[string]any{
			idParam:     graphID,
			sourceParam: source,
			targetParam: target,
			kParam:      amount,
		},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return nil, fmt.Errorf("find path via yens: %w", err)
	}

	return r.getPathsFromResult(result)
}
