package repository

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./concert_repository.go -destination=./mocks/concert_repository.go -package=repository_mocks
type ConcertRepository interface {
	CreateOne(ctx context.Context, concert *entity.Concert) (*entity.Concert, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Concert, error)
	FindAll(ctx context.Context, filter FindAllConcertsFilter) (*entity.Concerts, int64, error)
	WithTx(tx db.SqlExecer) ConcertRepository // Optional: WithTx if you want to use a transaction
}

type FindAllConcertsFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Venue     *string
	Limit     *int64
	Offset    *int64
	SortBy    *string
	SortOrder *entity.SortOrder
}
