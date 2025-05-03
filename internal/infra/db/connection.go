package db

import (
	"context"
	"fmt"
	"log"
	"ticket-reservation/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type dbContextKey struct{}

func FromContext(ctx context.Context) *sqlx.DB {
	return ctx.Value(dbContextKey{}).(*sqlx.DB)
}

func NewContext(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, dbContextKey{}, db)
}

func MustConnect(cfg *config.Config) *sqlx.DB {
	db, err := Connect(cfg, nil)
	if err != nil {
		log.Fatalln("failed to connect to DB:", err)
		return nil
	}
	return db
}

func Connect(cfg *config.Config, tracerProvider *sdktrace.TracerProvider) (db *sqlx.DB, err error) {
	if tracerProvider == nil {
		// Without tracing
		db, err = sqlx.Open("postgres", cfg.DB.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to DB (sqlx.Open): %w", err)
		}
	} else {
		// With tracing
		db, err = otelsqlx.Open("postgres", cfg.DB.URL, otelsql.WithTracerProvider(tracerProvider))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to DB (otelsqlx.Open): %w", err)
		}
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

	return db, nil
}
