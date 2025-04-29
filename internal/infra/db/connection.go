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

func Connect(cfg *config.Config, tracerProvider *sdktrace.TracerProvider) (*sqlx.DB, error) {
	dbConfig, err := cfg.DatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	var db *sqlx.DB
	if tracerProvider == nil {
		// Without tracing
		db, err = sqlx.Open("postgres", dbConfig.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to DB (sqlx.Open): %w", err)
		}
	} else {
		// With tracing
		db, err = otelsqlx.Open("postgres", dbConfig.URL, otelsql.WithTracerProvider(tracerProvider))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to DB (otelsqlx.Open): %w", err)
		}
	}

	// Set connection pool settings
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	db.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)

	return db, nil
}
