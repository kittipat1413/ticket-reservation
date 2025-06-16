package repository

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./zone_repository.go -destination=./mocks/zone_repository.go -package=repository_mocks
type ZoneRepository interface {
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Zone, error)
	WithTx(tx db.SqlExecer) ZoneRepository // Optional: WithTx if you want to use a transaction
}
