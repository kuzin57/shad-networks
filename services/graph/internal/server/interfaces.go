package server

import (
	"context"

	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"github.com/kuzin57/shad-networks/services/graph/internal/services/graph"
)

type GraphService interface {
	AddGraph(ctx context.Context, request entities.GraphParams) (entities.Graph, error)
	GetGraph(ctx context.Context, graphID string) (entities.Graph, []byte, string, error)
	FindPaths(ctx context.Context, graphID string, source, target, amount int) (entities.Graph, []entities.Path, graph.ImageChunk, error)
	ScrollGraphImage(ctx context.Context, scrollID string) ([]byte, bool, error)
}
