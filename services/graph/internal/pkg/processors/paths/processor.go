package paths

import (
	"context"

	"github.com/kuzin57/shad-networks/services/graph/internal/kafka"
)

type PathsProcessor struct {
}

func (p *PathsProcessor) Process(ctx context.Context, request kafka.Request) error {
	return nil
}
