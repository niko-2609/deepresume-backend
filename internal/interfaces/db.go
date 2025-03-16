package interfaces

import "context"

// Row represents a database row
type Row interface {
	Scan(dest ...interface{}) error
}

// Result represents the result of a database operation
type Result interface {
	RowsAffected() int64
}

// DB represents a database that can handle transactions
type DB interface {
	BeginTx(ctx context.Context) (Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
}

// Tx represents a database transaction
type Tx interface {
	Commit() error
	Rollback() error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
}
