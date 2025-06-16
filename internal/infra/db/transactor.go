package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	Commit() error
	Rollback() error
	DB() *sqlx.Tx
}

type SqlxTransactor interface {
	Transactor
}

type sqlxTransactor struct {
	tx *sqlx.Tx
}

func (s *sqlxTransactor) Commit() error {
	return s.tx.Commit()
}

func (s *sqlxTransactor) Rollback() error {
	return s.tx.Rollback()
}

func (s *sqlxTransactor) DB() *sqlx.Tx {
	return s.tx
}

// SqlxTransactorFactory creates a new SqlxTransactor (transaction)
type SqlxTransactorFactory interface {
	CreateSqlxTransactor(ctx context.Context, opts ...TxOptions) (SqlxTransactor, error)
}

// sqlxTransactorFactory implements SqlxTransactorFactory
type sqlxTransactorFactory struct {
	db *sqlx.DB
}

// NewSqlxTransactorFactory creates a new instance of SqlxTransactorFactory
func NewSqlxTransactorFactory(db *sqlx.DB) SqlxTransactorFactory {
	return &sqlxTransactorFactory{db: db}
}

// TxOptions defines options for creating a transaction
type TxOptions func(*sql.TxOptions)

// CreateSqlxTransactor creates a new SQL transaction with optional TxOptions
func (f *sqlxTransactorFactory) CreateSqlxTransactor(ctx context.Context, opts ...TxOptions) (SqlxTransactor, error) {
	txOptions := &sql.TxOptions{} // Allocate a new instance
	for _, opt := range opts {
		opt(txOptions) // Apply each option
	}
	tx, err := f.db.BeginTxx(ctx, txOptions)
	if err != nil {
		return nil, err
	}
	return &sqlxTransactor{tx: tx}, nil
}

// WithTxOptions returns a TxOptions function that copies the given options
func WithTxOptions(inputOptions *sql.TxOptions) TxOptions {
	return func(txOptions *sql.TxOptions) {
		*txOptions = *inputOptions
	}
}
