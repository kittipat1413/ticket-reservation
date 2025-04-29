package db

import (
	"context"

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
	CreateSqlxTransactor(ctx context.Context) (SqlxTransactor, error)
}

type sqlxTransactorFactory struct {
	db *sqlx.DB
}

func (f *sqlxTransactorFactory) CreateSqlxTransactor(ctx context.Context) (SqlxTransactor, error) {
	tx, err := f.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlxTransactor{tx: tx}, nil
}

func NewSqlxTransactorFactory(db *sqlx.DB) SqlxTransactorFactory {
	return &sqlxTransactorFactory{db: db}
}
