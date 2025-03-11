package repositories

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=repositories

type Driver interface {
	NewSession(ctx context.Context, config neo4j.SessionConfig) Session
	Get() neo4j.DriverWithContext
}

type Session interface {
	ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(*neo4j.TransactionConfig)) (any, error)
	Close(ctx context.Context) error
}
