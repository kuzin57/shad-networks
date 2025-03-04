package repository

import (
	"github.com/kuzin57/shad-networks/internal/config"
	"github.com/kuzin57/shad-networks/internal/repositories"
	"go.uber.org/zap"
)

type Repository struct {
	log    *zap.Logger
	driver Driver
	conf   *config.DBConfig
}

func NewRepository(driver *repositories.Neo4jDriver, log *zap.Logger, conf *config.Config) *Repository {
	return &Repository{
		driver: driver.DB(),
		conf:   conf.DB,
		log:    log,
	}
}
