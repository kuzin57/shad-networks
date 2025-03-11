package graph

import (
	"context"

	"github.com/kuzin57/shad-networks/internal/entities"
)

type GraphGenerator interface {
	Generate(entities.GraphParams) entities.Graph
}

type GraphRepository interface {
	CreateGraph(ctx context.Context, graph entities.Graph) error
	DropGraph(ctx context.Context, graphID string) error
	GetGraph(ctx context.Context, graphID string) (entities.Graph, error)
	ProjectIfNotExists(ctx context.Context, graphID string) error
	FindPath(ctx context.Context, graphID string, source, target int) (entities.Path, error)
}
