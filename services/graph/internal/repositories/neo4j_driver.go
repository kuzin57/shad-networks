package repositories

import (
	"context"
	"fmt"

	"github.com/kuzin57/shad-networks/services/graph/internal/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	_ GraphDbDriver = (*Neo4jDriver)(nil)
)

type Neo4jDriver struct {
	conf   *config.Neo4jConfig
	secret *config.Neo4jSecret
	log    *zap.Logger
	db     neo4j.DriverWithContext
}

func NewNeo4jDriver(
	lc fx.Lifecycle,
	log *zap.Logger,
	config *config.Config,
	secrets *config.Secrets,
) *Neo4jDriver {
	driver := &Neo4jDriver{
		conf:   config.Neo4j,
		secret: secrets.Neo4j,
		log:    log,
	}

	lc.Append(fx.Hook{
		OnStart: driver.Start,
		OnStop:  driver.Stop,
	})

	return driver
}

func (d *Neo4jDriver) Start(ctx context.Context) error {
	dbUri := fmt.Sprintf("neo4j://neo4j:%d", d.conf.Port)

	d.log.Info("creating new driver...")

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(d.secret.User, d.secret.Password, ""))
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

func (d *Neo4jDriver) Unwrap() neo4j.DriverWithContext {
	return d.db
}
