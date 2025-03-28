package interfaces

import (
	"context"

	"gorm.io/gorm"
)

// Row represents a database row
type Row interface {
	Scan(dest ...interface{}) error
}

// Result represents the result of a database operation
type Result interface {
	RowsAffected() int64
}

// DB represents a database that can handle transactions and GORM operations
type DB interface {
	BeginTx(ctx context.Context) (Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	WithContext(ctx context.Context) *gorm.DB
}

// Tx represents a database transaction
type Tx interface {
	Commit() error
	Rollback() error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Create(value interface{}) error
	Save(value interface{}) error
	Model(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) error
	Find(dest interface{}, conds ...interface{}) error
	Association(column string) *gorm.Association
	WithContext(ctx context.Context) *gorm.DB
	CreateInBatches(value interface{}, batchSize int) *gorm.DB
}
