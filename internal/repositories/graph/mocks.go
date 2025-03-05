package graph

import (
	"github.com/kuzin57/shad-networks/internal/repositories"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	driver  *repositories.MockDriver
	session *repositories.MockSession
}

func newMocks(ctrl *gomock.Controller) mocks {
	var (
		driver  = repositories.NewMockDriver(ctrl)
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
