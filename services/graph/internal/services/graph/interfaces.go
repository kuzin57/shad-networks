package graph

import (
	"context"

	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
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
	FindPaths(ctx context.Context, graphID string, source, target, amount int) ([]entities.Path, error)
}

type GraphCache interface {
	CheckGraphExists(ctx context.Context, graphID string) (bool, error)
	PutGraph(ctx context.Context, graphID string) error
	PutFileData(ctx context.Context, scrollID string, content []byte) ([]byte, error)
	GetFileChunk(ctx context.Context, scrollID string) ([]byte, bool, error)
}

type Visualizer interface {
	Visualize(ctx context.Context, graph entities.Graph, paths []entities.Path) ([]byte, error)
}
