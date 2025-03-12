package repositories

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=repositories

type GraphDbDriver interface {
	NewSession(ctx context.Context, config neo4j.SessionConfig) Session
	Unwrap() neo4j.DriverWithContext
}

type Session interface {
	ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(*neo4j.TransactionConfig)) (any, error)
	Close(ctx context.Context) error
}

type CacheDriver interface {
	Unwrap() *redis.Client
}
