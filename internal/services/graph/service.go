package graph

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kuzin57/shad-networks/internal/entities"
	"go.uber.org/zap"
)

type Service struct {
	log       *zap.Logger
	generator GraphGenerator
	repo      GraphRepository
}

func NewService(log *zap.Logger, generator GraphGenerator, repo GraphRepository) *Service {
	return &Service{
		log:       log,
		generator: generator,
		repo:      repo,
	}
}

func (s *Service) AddGraph(ctx context.Context, params entities.GraphParams) (entities.Graph, error) {
	graph := s.generator.Generate(params)
	graph.ID = uuid.NewString()

	if err := s.repo.CreateGraph(ctx, graph); err != nil {
		return entities.Graph{}, fmt.Errorf("create graph: %w", err)
	}

	return graph, nil
}
