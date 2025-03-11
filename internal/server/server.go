package server

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/internal/entities"
	"github.com/kuzin57/shad-networks/internal/generated"
)

var (
	_ generated.GraphServer = (*Server)(nil)
)

type Server struct {
	generated.UnimplementedGraphServer
	graphService GraphService
	visualizer   Visualizer
}

func NewServer(graphService GraphService, visualizer Visualizer) *Server {
	return &Server{
		graphService: graphService,
	}
}

func (s *Server) Add(ctx context.Context, request *generated.AddGraphRequest) (*generated.AddGraphResponse, error) {
	req := entities.GraphParams{
		VerticesCount: request.VerticesCount,
		Weights:       request.Weights,
		Degrees:       request.Degrees,
	}

	graph, err := s.graphService.AddGraph(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("add graph: %w", err)
	}

	return &generated.AddGraphResponse{
		GraphId: graph.ID,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *generated.GetGraphRequest) (*generated.GetGraphResponse, error) {
	graph, err := s.graphService.GetGraph(ctx, request.GraphId)
	if err != nil {
		return nil, fmt.Errorf("get graph: %w", err)
	}

	b64Image, err := s.visualizer.Visualize(ctx, graph, nil)
	if err != nil {
		return nil, fmt.Errorf("visualize image: %w", err)
	}

	return &generated.GetGraphResponse{
		B64Image: b64Image,
	}, nil
}

func (s *Server) FindPath(ctx context.Context, request *generated.FindPathRequest) (*generated.FindPathResponse, error) {
	graph, path, err := s.graphService.FindPath(ctx, request.GraphId, int(request.From), int(request.To))
	if err != nil {
		return nil, fmt.Errorf("find path: %w", err)
	}

	b64Image, err := s.visualizer.Visualize(ctx, graph, path)
	if err != nil {
		return nil, fmt.Errorf("visualize image: %w", err)
	}

	return &generated.FindPathResponse{
		B64Image: b64Image,
	}, nil
}
