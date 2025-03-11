package graphcleaner

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	graphsIDsQueueCapacity = 100
)

type Cleaner struct {
	graphsIDs chan string
	done      chan struct{}
	repo      GraphRepository
	log       *zap.Logger
}

func NewCleaner(lc fx.Lifecycle, repo GraphRepository, log *zap.Logger) *Cleaner {
	cleaner := &Cleaner{
		repo:      repo,
		done:      make(chan struct{}),
		graphsIDs: make(chan string, graphsIDsQueueCapacity),
		log:       log,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go cleaner.Run(ctx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			close(cleaner.done)

			return nil
		},
	})

	return cleaner
}

func (c *Cleaner) Run(ctx context.Context) {
	for {
		select {
		case <-c.done:
			return
		case graphID := <-c.graphsIDs:
			if err := c.repo.DropGraph(ctx, graphID); err != nil {
				c.log.Error("drop graph", zap.Error(err))
			}
		}
	}
}

func (c *Cleaner) AddForCleaning(graphID string) {
	c.graphsIDs <- graphID
}
