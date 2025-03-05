package graph

import (
	"github.com/kuzin57/shad-networks/internal/repositories"
	"go.uber.org/zap"
)

type Repository struct {
	log    *zap.Logger
	driver repositories.Driver
}

func NewRepository(driver repositories.Driver, log *zap.Logger) *Repository {
	return &Repository{
		driver: driver,
		log:    log,
	}
}
