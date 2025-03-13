package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kuzin57/shad-networks/services/graph/internal/repositories"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CacheRepository struct {
	done      chan struct{}
	log       *zap.Logger
	pubsub    *redis.PubSub
	driver    repositories.CacheDriver
	processor ExpiredGraphCleaner
}

func NewCacheRepository(
	lc fx.Lifecycle,
	driver repositories.CacheDriver,
	processor ExpiredGraphCleaner,
	log *zap.Logger,
) *CacheRepository {
	repo := &CacheRepository{
		driver:    driver,
		log:       log,
		processor: processor,
		done:      make(chan struct{}),
	}

	lc.Append(fx.Hook{
		OnStart: repo.Start,
		OnStop:  repo.Stop,
	})

	return repo
}

func (r *CacheRepository) Start(ctx context.Context) error {
	r.log.Info("starting cache...")

	r.pubsub = r.driver.Unwrap().PSubscribe(ctx, expirationEvent)

	go r.Run(ctx)

	return nil
}

func (r *CacheRepository) Stop(ctx context.Context) error {
	close(r.done)

	return nil
}

func (r *CacheRepository) Run(ctx context.Context) {
	r.log.Info("listening for messages...")

	for {
		message, err := r.pubsub.ReceiveMessage(ctx)
		if err != nil {
			r.log.Error("received error", zap.Error(err))

			break
		}

		if strings.HasPrefix(message.Payload, defaultGraphIDKeyPrefix) {
			r.log.Sugar().Infof("key expired: %s", message.Payload)

			if err := r.processor.DropGraph(ctx, strings.TrimPrefix(message.Payload, defaultGraphIDKeyPrefix)); err != nil {
				r.log.Error("drop graph", zap.Error(err))
			}
		}
	}
}

func (r *CacheRepository) PutFileData(ctx context.Context, scrollID string, content []byte) ([]byte, error) {
	if len(content) <= defaultFileChunkSize {
		return content, nil
	}

	if err := r.driver.Unwrap().SetNX(ctx, scrollID, content[defaultFileChunkSize:], defaultScrollTTL).Err(); err != nil {
		return nil, fmt.Errorf("set key %s: %w", scrollID, err)
	}

	return content[:defaultFileChunkSize], nil
}

func (r *CacheRepository) GetFileChunk(ctx context.Context, scrollID string) ([]byte, bool, error) {
	value, err := r.driver.Unwrap().Get(ctx, scrollID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, fmt.Errorf("%w: key: %s", ErrNotExists, scrollID)
		}

		return nil, false, fmt.Errorf("get by key %s: %w", scrollID, err)
	}

	if len(value) <= defaultFileChunkSize {
		return []byte(value), true, nil
	}

	if err := r.driver.Unwrap().Set(ctx, scrollID, value[defaultFileChunkSize:], redis.KeepTTL).Err(); err != nil {
		return nil, false, fmt.Errorf("set key %s: %w", scrollID, err)
	}

	return []byte(value[:defaultFileChunkSize]), false, nil
}

func (r *CacheRepository) PutGraph(ctx context.Context, graphID string) error {
	if err := r.driver.Unwrap().SetNX(ctx, r.buildGraphIDKey(graphID), "value", defaultGraphTTL).Err(); err != nil {
		return fmt.Errorf("put graph to cache: %w", err)
	}

	return nil
}

func (r *CacheRepository) CheckGraphExists(ctx context.Context, graphID string) (bool, error) {
	key := r.buildGraphIDKey(graphID)

	_, err := r.driver.Unwrap().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("get graph by key %s: %w", key, err)
	}

	return true, nil
}

func (r *CacheRepository) buildGraphIDKey(graphID string) string {
	return defaultGraphIDKeyPrefix + graphID
}
