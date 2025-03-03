package repository

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/internal/config"
	"github.com/kuzin57/shad-networks/internal/consts"
	"github.com/kuzin57/shad-networks/internal/entities"
	"github.com/kuzin57/shad-networks/internal/repositories"
	"github.com/kuzin57/shad-networks/internal/repositories/graph/queries"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

type Repository struct {
	log    *zap.Logger
	driver *repositories.Neo4jDriver
	conf   *config.DBConfig
}

func NewRepository(driver *repositories.Neo4jDriver, log *zap.Logger, conf *config.Config) *Repository {
	return &Repository{
		driver: driver,
		conf:   conf.DB,
		log:    log,
	}
}

func (r *Repository) CreateGraph(ctx context.Context, graph entities.Graph) error {
	driver := r.driver.DB()

	for i := range graph.AdjencyMaxtrix {
		r.log.Sugar().Infof("creating node#%d", i)

		_, err := neo4j.ExecuteQuery(
			ctx,
			driver,
			queries.CreateNode,
			map[string]any{
				"graphID": graph.ID,
				"number":  i,
			},
			neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase(consts.DatabaseName),
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

			if err := r.createEdge(ctx, i, j, &graph); err != nil {
				return fmt.Errorf("add new edge: %w", err)
			}

			if err := r.createEdge(ctx, j, i, &graph); err != nil {
				return fmt.Errorf("add new edge: %w", err)
			}
		}
	}

	return nil
}

func (r *Repository) createEdge(ctx context.Context, i, j int, graph *entities.Graph) error {
	_, err := neo4j.ExecuteQuery(
		ctx,
		r.driver.DB(),
		queries.CreateEdge,
		map[string]any{
			"firstNumber":  i,
			"secondNumber": j,
			"graphID":      graph.ID,
			"weight":       graph.AdjencyMaxtrix[i][j],
		},
		neo4j.EagerResultTransformer,
	)

	return err
}
