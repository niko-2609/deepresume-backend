package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/nikolai/ai-resume-builder/backend/internal/interfaces"
)

// DB wraps sql.DB to implement interfaces.DB
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

// BeginTx starts a new transaction
func (db *DB) BeginTx(ctx context.Context) (interfaces.Tx, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

// QueryRowContext executes a query that returns at most one row
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) interfaces.Row {
	return &Row{db.DB.QueryRowContext(ctx, query, args...)}
}

// ExecContext executes a query without returning any rows
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (interfaces.Result, error) {
	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Result{result}, nil
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
	sql.Result
}

// RowsAffected returns the number of rows affected by the query
func (r *Result) RowsAffected() int64 {
	n, _ := r.Result.RowsAffected()
	return n
}

// Tx wraps sql.Tx
type Tx struct {
	*sql.Tx
}

// QueryRowContext executes a query that returns at most one row within a transaction
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) interfaces.Row {
	return &Row{tx.Tx.QueryRowContext(ctx, query, args...)}
}

// ExecContext executes a query without returning any rows within a transaction
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (interfaces.Result, error) {
	result, err := tx.Tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Result{result}, nil
}
