package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/kuzin57/shad-networks/internal/entities"
	entitiesmocks "github.com/kuzin57/shad-networks/internal/mocks/entities"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateGraph(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(createGraphTestSuite))
}

type createGraphTestSuite struct {
	suite.Suite
	mocks

	graph entities.Graph

	repository *Repository
}

func (s *createGraphTestSuite) SetupSuite() {
	s.graph = entitiesmocks.GetMockGraph()
}

func (s *createGraphTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.mocks = newMocks(ctrl)
	s.repository = NewRepository(s.driver, zap.NewExample())
}

func (s *createGraphTestSuite) TestSuccessPath() {
	ctx := context.Background()

	gomock.InOrder(
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			Close(ctx).
			Return(nil),
	)

	err := s.repository.CreateGraph(ctx, s.graph)
	s.Require().NoError(err)
}

func (s *createGraphTestSuite) TestExecuteWriteError() {
	ctx := context.Background()

	gomock.InOrder(
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, nil),
		s.session.EXPECT().
			ExecuteWrite(ctx, gomock.Any()).
			Return(nil, errors.New("some error")),
		s.session.EXPECT().
			Close(ctx).
			Return(nil),
	)

	err := s.repository.CreateGraph(ctx, s.graph)
	s.Require().Error(err)
}
