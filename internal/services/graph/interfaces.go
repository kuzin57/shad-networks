package graph

import (
	"context"

	"github.com/kuzin57/shad-networks/internal/entities"
)

type GraphGenerator interface {
	Generate(entities.GraphParams) entities.Graph
}

type GraphRepository interface {
	CreateGraph(context.Context, entities.Graph) error
}
