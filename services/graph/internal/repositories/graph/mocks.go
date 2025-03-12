package graph

import (
	"github.com/kuzin57/shad-networks/services/graph/internal/repositories"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	driver  *repositories.MockGraphDbDriver
	session *repositories.MockSession
}

func newMocks(ctrl *gomock.Controller) mocks {
	var (
		driver  = repositories.NewMockGraphDbDriver(ctrl)
		session = repositories.NewMockSession(ctrl)
	)

	driver.EXPECT().
		NewSession(gomock.Any(), gomock.Any()).
		Return(session).AnyTimes()

	return mocks{
		driver:  driver,
		session: session,
	}
}
