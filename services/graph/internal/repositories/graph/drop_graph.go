package graph

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/services/graph/internal/repositories/graph/queries"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) DropGraph(ctx context.Context, graphID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		_ = session.Close(ctx)
	}()

	_, err := session.ExecuteWrite(
		ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			return r.dropGraph(ctx, tx, graphID)
		},
	)
	if err != nil {
		return fmt.Errorf("add new node: %w", err)
	}

	return nil
}

func (r *Repository) dropGraph(ctx context.Context, tx neo4j.ManagedTransaction, graphID string) (any, error) {
	_, err := tx.Run(
		ctx,
		queries.DropGraph,
		map[string]any{
			idParam: graphID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("project graph: %w", err)
	}

	return nil, nil
}
