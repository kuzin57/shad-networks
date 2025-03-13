package server

import (
	"context"

	"github.com/kuzin57/shad-networks/generated"
	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"go.uber.org/zap"
)

var (
	_ generated.GraphServer = (*Server)(nil)
)

type Server struct {
	generated.UnimplementedGraphServer
	graphService GraphService
	log          *zap.Logger
}

func NewServer(graphService GraphService, log *zap.Logger) *Server {
	return &Server{
		graphService: graphService,
		log:          log,
	}
}

func (s *Server) Add(ctx context.Context, request *generated.AddGraphRequest) (*generated.AddGraphResponse, error) {
	req := entities.GraphParams{
		VerticesCount:    request.VerticesCount,
		Weights:          request.Weights,
		Degrees:          request.Degrees,
		MaxMultipleEdges: request.MaxMultipleEdges,
	}

	graph, err := s.graphService.AddGraph(ctx, req)
	if err != nil {
		s.log.Error("add graph failed", zap.Error(err))

		return nil, err
	}

	return &generated.AddGraphResponse{
		GraphId: graph.ID,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *generated.GetGraphRequest) (*generated.GetGraphResponse, error) {
	_, b64Image, scrollID, err := s.graphService.GetGraph(ctx, request.GraphId)
	if err != nil {
		s.log.Error("get graph failed", zap.Error(err))

		return nil, err
	}

	return &generated.GetGraphResponse{
		B64Image: b64Image,
		ScrollId: scrollID,
	}, nil
}

func (s *Server) FindPath(ctx context.Context, request *generated.FindPathRequest) (*generated.FindPathResponse, error) {
	_, _, imageChunk, err := s.graphService.FindPath(ctx, request.GraphId, int(request.From), int(request.To))
	if err != nil {
		s.log.Error("find path failed", zap.Error(err))

		return nil, err
	}

	return &generated.FindPathResponse{
		B64Image: imageChunk.Content,
		ScrollId: imageChunk.ScrollID,
	}, nil
}

func (s *Server) Scroll(ctx context.Context, request *generated.ScrollRequest) (*generated.ScrollResponse, error) {
	chunk, isOver, err := s.graphService.ScrollGraphImage(ctx, request.ScrollId)
	if err != nil {
		s.log.Error("scroll graph image", zap.Error(err))

		return nil, err
	}

	return &generated.ScrollResponse{
		ScrollId: request.ScrollId,
		IsOver:   isOver,
		B64Image: chunk,
	}, nil
}
