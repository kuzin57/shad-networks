package graph

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/services/graph/internal/repositories/graph/queries"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) ProjectIfNotExists(ctx context.Context, graphID string) error {
	result, err := neo4j.ExecuteQuery(
		ctx,
		r.driver.Unwrap(),
		queries.CheckGraphExists,
		map[string]any{
			idParam: graphID,
		},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return fmt.Errorf("check exists: %w", err)
	}

	if len(result.Records) == 0 || len(result.Records) > 1 {
		return fmt.Errorf("got invalid records amount on exists query: %d", len(result.Records))
	}

	record := result.Records[0]
	if record == nil {
		return fmt.Errorf("found nil record")
	}

	exists, ok := record.Get(existsParam)
	if !ok {
		return fmt.Errorf("exists property not found in result record")
	}

	if v, _ := exists.(bool); v {
		return nil
	}

	return r.restoreGraph(ctx, graphID)
}

func (r *Repository) restoreGraph(ctx context.Context, graphID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close(ctx)
	}()

	_, err := session.ExecuteWrite(
		ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			return r.projectGraph(ctx, tx, graphID)
		},
	)
	if err != nil {
		return fmt.Errorf("add new node: %w", err)
	}

	return nil
}

func (r *Repository) projectGraph(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	graphID string,
) (any, error) {
	_, err := tx.Run(
		ctx,
		queries.ProjectGraph,
		map[string]any{
			idParam: graphID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("project graph: %w", err)
	}

	return nil, nil
}
