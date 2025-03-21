package graph

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/services/graph/internal/consts"
	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"github.com/kuzin57/shad-networks/services/graph/internal/repositories/graph/queries"
	jsonutils "github.com/kuzin57/shad-networks/services/graph/internal/utils/json"
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

	var edgeNumber int
	for i := range graph.AdjencyMaxtrix {
		for j := i + 1; j < len(graph.AdjencyMaxtrix); j++ {
			if len(graph.AdjencyMaxtrix[i][j]) == 0 {
				continue
			}

			for _, edgeWeight := range graph.AdjencyMaxtrix[i][j] {
				graphEdge.From = i
				graphEdge.To = j
				graphEdge.Weight = edgeWeight
				graphEdge.Number = edgeNumber

				_, err := session.ExecuteWrite(
					ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						return r.СreateEdge(ctx, tx, graphEdge)
					},
				)
				if err != nil {
					return fmt.Errorf("creating edge from %d to %d: %w", i, j, err)
				}

				edgeNumber++
			}
		}
	}

	r.log.Info("success fully created edges")

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.projectGraph(ctx, tx, graph.ID)
	}); err != nil {
		return err
	}

	r.log.Info("successfully created graph")

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
