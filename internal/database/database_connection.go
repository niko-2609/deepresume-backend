package database

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB wraps gorm.DB to implement interfaces.DB
type DB struct {
	*gorm.DB
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// NewDB creates a new database connection
func NewDB(dataSourceName string) (*DB, error) {
	// Remove any quotes from the connection string
	dataSourceName = strings.Trim(dataSourceName, "'\"")

	gormDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &DB{DB: gormDB}, nil
}

// BeginTx starts a new transaction
func (db *DB) BeginTx(ctx context.Context) (interfaces.Tx, error) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &Tx{tx}, nil
}

// QueryRowContext executes a query that returns at most one row
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) interfaces.Row {
	// This is a placeholder since GORM doesn't have direct equivalent
	// We'll need to implement this differently or use GORM's methods
	return nil
}

// ExecContext executes a query without returning any rows
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (interfaces.Result, error) {
	result := db.DB.Exec(query, args...)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Result{result}, nil
}

// WithContext returns a new DB instance with the given context
func (db *DB) WithContext(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx)
}

// Row wraps sql.Row
type Row struct {
	*sql.Row
}

// Scan copies the columns from the matched row into the values pointed at by dest
func (r *Row) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}

// Result wraps sql.Result
type Result struct {
	*gorm.DB
}

// RowsAffected returns the number of rows affected by the query
func (r *Result) RowsAffected() int64 {
	return r.DB.RowsAffected
}

// Tx wraps gorm.DB
type Tx struct {
	*gorm.DB
}

// Commit commits the transaction
func (tx *Tx) Commit() error {
	return tx.DB.Commit().Error
}

// Rollback rolls back the transaction
func (tx *Tx) Rollback() error {
	return tx.DB.Rollback().Error
}

// QueryRowContext executes a query that returns at most one row within a transaction
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) interfaces.Row {
	// This is a placeholder since GORM doesn't have direct equivalent
	// We'll need to implement this differently or use GORM's methods
	return nil
}

// ExecContext executes a query without returning any rows within a transaction
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (interfaces.Result, error) {
	result := tx.DB.Exec(query, args...)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Result{result}, nil
}

// Create creates a new record in the database
func (tx *Tx) Create(value interface{}) error {
	return tx.DB.Create(value).Error
}

// Save saves a record to the database
func (tx *Tx) Save(value interface{}) error {
	return tx.DB.Save(value).Error
}

// Model specify the model you would like to run db operations
func (tx *Tx) Model(value interface{}) *gorm.DB {
	return tx.DB.Model(value)
}

// Where specify the where conditions
func (tx *Tx) Where(query interface{}, args ...interface{}) *gorm.DB {
	return tx.DB.Where(query, args...)
}

// Preload preload associations
func (tx *Tx) Preload(query string, args ...interface{}) *gorm.DB {
	return tx.DB.Preload(query, args...)
}

// First find first record that match given conditions
func (tx *Tx) First(dest interface{}, conds ...interface{}) error {
	return tx.DB.First(dest, conds...).Error
}

// Find find records that match given conditions
func (tx *Tx) Find(dest interface{}, conds ...interface{}) error {
	return tx.DB.Find(dest, conds...).Error
}

// Association returns association handler
func (tx *Tx) Association(column string) *gorm.Association {
	return tx.DB.Association(column)
}

// WithContext set context for transaction
func (tx *Tx) WithContext(ctx context.Context) *gorm.DB {
	return tx.DB.WithContext(ctx)
}

// CreateInBatches creates records in batches
func (tx *Tx) CreateInBatches(value interface{}, batchSize int) *gorm.DB {
	return tx.DB.CreateInBatches(value, batchSize)
}
