package repositories

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/internal/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Neo4jDriver struct {
	conf *config.DBConfig
	log  *zap.Logger
	db   neo4j.DriverWithContext
}

func NewNeo4jDriver(lc fx.Lifecycle, log *zap.Logger, config *config.Config) *Neo4jDriver {
	driver := &Neo4jDriver{conf: config.DB, log: log}

	lc.Append(fx.Hook{
		OnStart: driver.Start,
		OnStop:  driver.Stop,
	})

	return driver
}

func (d *Neo4jDriver) Start(ctx context.Context) error {
	dbUri := fmt.Sprintf("neo4j://db:%d", d.conf.Port)

	d.log.Info("creating new driver...")

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(d.conf.User, d.conf.Password, ""))
	if err != nil {
		return fmt.Errorf("can not create neo4j driver with context: %w", err)
	}

	d.db = driver

	return nil
}

func (d *Neo4jDriver) Stop(ctx context.Context) error {
	return d.db.Close(ctx)
}

func (d *Neo4jDriver) NewSession(ctx context.Context, config neo4j.SessionConfig) Session {
	return d.db.NewSession(ctx, config)
}
