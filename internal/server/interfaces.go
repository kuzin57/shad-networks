package server

import (
	"context"

	"github.com/kuzin57/shad-networks/internal/entities"
)

type GraphService interface {
	AddGraph(ctx context.Context, request entities.GraphParams) (entities.Graph, error)
	GetGraph(ctx context.Context, graphID string) (entities.Graph, error)
	FindPath(ctx context.Context, graphID string, source, target int) (entities.Graph, entities.Path, error)
}

type Visualizer interface {
	Visualize(ctx context.Context, graph entities.Graph, path entities.Path) ([]byte, error)
}
