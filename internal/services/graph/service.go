package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuzin57/shad-networks/internal/entities"
	graphcleaner "github.com/kuzin57/shad-networks/internal/pkg/graph_cleaner"
	"github.com/kuzin57/shad-networks/pkg/lru"
	"go.uber.org/zap"
)

const (
	maxCacheSize       = 10
	graphProjectionTTL = time.Minute * 3
)

type empty struct{}

type Service struct {
	log       *zap.Logger
	generator GraphGenerator
	repo      GraphRepository
	cleaner   *graphcleaner.Cleaner
	cache     *lru.LRUCache[string, empty]
}

func NewService(
	log *zap.Logger,
	generator GraphGenerator,
	repo GraphRepository,
	cleaner *graphcleaner.Cleaner,
) *Service {
	return &Service{
		log:       log,
		generator: generator,
		repo:      repo,
		cache:     lru.NewLRUCache[string, empty](int(maxCacheSize), graphProjectionTTL),
		cleaner:   cleaner,
	}
}

func (s *Service) AddGraph(ctx context.Context, params entities.GraphParams) (entities.Graph, error) {
	graph := s.generator.Generate(params)
	graph.ID = uuid.NewString()

	for _, graphID := range s.cache.Put(graph.ID, empty{}) {
		s.cleaner.AddForCleaning(graphID)
	}

	if err := s.repo.CreateGraph(ctx, graph); err != nil {
		return entities.Graph{}, fmt.Errorf("create graph: %w", err)
	}

	return graph, nil
}

func (s *Service) GetGraph(ctx context.Context, graphID string) (entities.Graph, error) {
	graph, err := s.repo.GetGraph(ctx, graphID)
	if err != nil {
		return entities.Graph{}, fmt.Errorf("get graph: %w", err)
	}

	return graph, nil
}

func (s *Service) FindPath(ctx context.Context, graphID string, source, target int) (entities.Graph, entities.Path, error) {
	_, exists := s.cache.Get(graphID)
	if !exists {
		if err := s.repo.ProjectIfNotExists(ctx, graphID); err != nil {
			return entities.Graph{}, nil, fmt.Errorf("project graph: %w", err)
		}

		for _, id := range s.cache.Put(graphID, empty{}) {
			s.cleaner.AddForCleaning(id)
		}
	}

	graph, err := s.repo.GetGraph(ctx, graphID)
	if err != nil {
		return entities.Graph{}, nil, fmt.Errorf("get graph: %w", err)
	}

	path, err := s.repo.FindPath(ctx, graphID, source, target)
	if err != nil {
		return entities.Graph{}, nil, fmt.Errorf("find path: %w", err)
	}

	return graph, path, nil
}
