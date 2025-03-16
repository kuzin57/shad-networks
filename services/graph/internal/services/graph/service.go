package graph

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"go.uber.org/zap"
)

type Service struct {
	log        *zap.Logger
	generator  GraphGenerator
	repo       GraphRepository
	cache      GraphCache
	visualizer Visualizer
}

func NewService(
	log *zap.Logger,
	generator GraphGenerator,
	repo GraphRepository,
	cache GraphCache,
	visualizer Visualizer,
) *Service {
	return &Service{
		log:        log,
		generator:  generator,
		repo:       repo,
		cache:      cache,
		visualizer: visualizer,
	}
}

func (s *Service) AddGraph(ctx context.Context, params entities.GraphParams) (entities.Graph, error) {
	graph := s.generator.Generate(params)
	graph.ID = uuid.NewString()

	if err := s.cache.PutGraph(ctx, graph.ID); err != nil {
		s.log.Error("put graph to cache", zap.Error(err))
	}

	if err := s.repo.CreateGraph(ctx, graph); err != nil {
		return entities.Graph{}, fmt.Errorf("create graph: %w", err)
	}

	return graph, nil
}

func (s *Service) GetGraph(ctx context.Context, graphID string) (entities.Graph, []byte, string, error) {
	graph, err := s.repo.GetGraph(ctx, graphID)
	if err != nil {
		return entities.Graph{}, nil, "", fmt.Errorf("get graph: %w", err)
	}

	s.log.Info("got graph, starting visualizing...")

	b64Image, err := s.visualizer.Visualize(ctx, graph, nil)
	if err != nil {
		s.log.Error("graph vizualization failed", zap.Error(err))

		return entities.Graph{}, nil, "", fmt.Errorf("visualize image: %w", err)
	}

	scrollID := uuid.NewString()

	firstChunk, err := s.cache.PutFileData(ctx, scrollID, b64Image)
	if err != nil {
		return entities.Graph{}, nil, "", fmt.Errorf("put file data: %w", err)
	}

	return graph, firstChunk, scrollID, nil
}

func (s *Service) ScrollGraphImage(ctx context.Context, scrollID string) ([]byte, bool, error) {
	chunk, isOver, err := s.cache.GetFileChunk(ctx, scrollID)
	if err != nil {
		return nil, false, fmt.Errorf("get file chunk: %w", err)
	}

	return chunk, isOver, nil
}

func (s *Service) FindPaths(
	ctx context.Context,
	graphID string,
	source,
	target,
	amount int,
) (
	entities.Graph,
	[]entities.Path,
	ImageChunk,
	error,
) {
	exists, err := s.cache.CheckGraphExists(ctx, graphID)
	switch {
	case err != nil:
		return entities.Graph{}, nil, ImageChunk{}, fmt.Errorf("check graph exists: %w", err)
	case !exists:
		if err := s.repo.ProjectIfNotExists(ctx, graphID); err != nil {
			return entities.Graph{}, nil, ImageChunk{}, fmt.Errorf("project graph: %w", err)
		}

		if err := s.cache.PutGraph(ctx, graphID); err != nil {
			return entities.Graph{}, nil, ImageChunk{}, fmt.Errorf("put graph to cache: %w", err)
		}
	}

	graph, err := s.repo.GetGraph(ctx, graphID)
	if err != nil {
		return entities.Graph{}, nil, ImageChunk{}, fmt.Errorf("get graph: %w", err)
	}

	paths, err := s.repo.FindPaths(ctx, graphID, source, target, amount)
	if err != nil {
		return entities.Graph{}, nil, ImageChunk{}, fmt.Errorf("find paths: %w", err)
	}

	b64Image, err := s.visualizer.Visualize(ctx, graph, paths)
	if err != nil {
		s.log.Error("graph vizualization failed", zap.Error(err))

		return entities.Graph{}, nil, ImageChunk{}, fmt.Errorf("visualize image: %w", err)
	}

	scrollID := uuid.NewString()

	firstChunk, err := s.cache.PutFileData(ctx, scrollID, b64Image)
	if err != nil {
		return entities.Graph{}, nil, ImageChunk{}, fmt.Errorf("put file data: %w", err)
	}

	return graph, paths, ImageChunk{Content: firstChunk, ScrollID: scrollID}, nil
}
