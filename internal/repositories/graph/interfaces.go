package repository

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

//go:generate mockgen
type DriverWrapper interface {
	DB() neo4j.DriverWithContext
}

type Driver interface {
	NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext
}

type Session interface {
	ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork, configurers ...func(*neo4j.TransactionConfig)) (any, error)
}
