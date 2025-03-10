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
}

func NewServer(graphService GraphService) *Server {
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
	return nil, nil
}

func (s *Server) FindPath(ctx context.Context, request *generated.FindPathRequest) (*generated.FindPathResponse, error) {
	return nil, nil
}
