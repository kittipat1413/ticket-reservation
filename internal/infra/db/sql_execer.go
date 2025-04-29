package db

import (
	"context"
	"database/sql"
)

// SqlExecer is an interface that generalizes the methods of *sql.DB and *sql.Tx
type SqlExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
