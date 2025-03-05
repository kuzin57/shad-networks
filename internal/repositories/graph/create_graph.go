package graph

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/internal/consts"
	"github.com/kuzin57/shad-networks/internal/entities"
	"github.com/kuzin57/shad-networks/internal/repositories/graph/queries"
	jsonutils "github.com/kuzin57/shad-networks/internal/utils/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) CreateGraph(ctx context.Context, graph entities.Graph) error {
	var (
		session = r.driver.NewSession(ctx, neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
		})
		graphNode = &entities.GraphNode{GraphID: graph.ID}
		graphEdge = &entities.GraphEdge{GraphID: graph.ID, Connection: consts.EdgeConnectionName}
	)

	defer func() {
		_ = session.Close(ctx)
	}()

	for i := range graph.AdjencyMaxtrix {
		r.log.Sugar().Infof("creating node#%d", i)

		graphNode.Number = i

		_, err := session.ExecuteWrite(
			ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				return r.createNode(ctx, tx, graphNode)
			},
		)
		if err != nil {
			return fmt.Errorf("add new node: %w", err)
		}
	}

	r.log.Info("successfully created nodes")

	for i := range graph.AdjencyMaxtrix {
		for j := i + 1; j < len(graph.AdjencyMaxtrix); j++ {
			if graph.AdjencyMaxtrix[i][j] == 0 {
				continue
			}

			graphEdge.From = i
			graphEdge.To = j
			graphEdge.Weight = graph.AdjencyMaxtrix[i][j]

			_, err := session.ExecuteWrite(
				ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					return r.createEdge(ctx, tx, graphEdge)
				},
			)
			if err != nil {
				return fmt.Errorf("creating edge from %d to %d: %w", i, j, err)
			}
		}
	}

	return nil
}

func (r *Repository) createNode(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	node *entities.GraphNode,
) (any, error) {
	records, err := tx.Run(
		ctx,
		queries.CreateNode,
		jsonutils.Serialize(node),
	)
	if err != nil {
		return nil, fmt.Errorf("create node: %w", err)
	}

	_, err = records.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("got single record: %w", err)
	}

	return nil, nil
}

func (r *Repository) createEdge(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	edge *entities.GraphEdge,
) (any, error) {
	records, err := tx.Run(
		ctx,
		queries.CreateEdge,
		jsonutils.Serialize(edge),
	)
	if err != nil {
		return nil, fmt.Errorf("create edge: %w", err)
	}

	_, err = records.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("got single record: %w", err)
	}

	return nil, nil
}
