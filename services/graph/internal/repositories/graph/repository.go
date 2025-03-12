package graph

import (
	"github.com/kuzin57/shad-networks/services/graph/internal/repositories"
	"go.uber.org/zap"
)

type Repository struct {
	log    *zap.Logger
	driver repositories.GraphDbDriver
}

func NewRepository(driver repositories.GraphDbDriver, log *zap.Logger) *Repository {
	return &Repository{
		driver: driver,
		log:    log,
	}
}
