package server

import (
	"context"

	"github.com/kuzin57/shad-networks/internal/entities"
)

type GraphService interface {
	AddGraph(ctx context.Context, request entities.GraphParams) (entities.Graph, error)
}
