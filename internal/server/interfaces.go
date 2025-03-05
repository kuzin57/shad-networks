package server

import (
	"context"

	"github.com/kuzin57/shad-networks/internal/entities"
)

type GraphService interface {
	AddGraph(ctx context.Context, request entities.GraphParams) (entities.Graph, error)
	GetGraph(ctx context.Context, graphID string) (entities.Graph, error)
	FindPath(ctx context.Context, graphID string, from, to int) (entities.Graph, []int, error)
}
