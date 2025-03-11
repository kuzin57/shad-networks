package graphcleaner

import "context"

type GraphRepository interface {
	DropGraph(ctx context.Context, graphID string) error
}
