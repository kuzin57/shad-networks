package cache

import "context"

type ExpiredGraphCleaner interface {
	DropGraph(ctx context.Context, graphID string) error
}
