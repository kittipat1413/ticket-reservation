package cache

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

const (
	// SeatMapNoExpiration represents no TTL (permanent storage)
	SeatMapNoExpiration time.Duration = 0
)

//go:generate mockgen -source=./seat_map_repository.go -destination=./mocks/seat_map_repository.go -package=cache_mocks
type SeatMapRepository interface {
	// SetSeat updates a complete seat entity with field-level TTL
	SetSeat(ctx context.Context, concertID, zoneID uuid.UUID, seat entity.Seat, ttl time.Duration) error
	// GetSeat retrieves a complete seat entity
	GetSeat(ctx context.Context, concertID, zoneID uuid.UUID, seatNumber string) (*entity.Seat, error)
	// GetAllSeats retrieves all seat entities for a zone
	GetAllSeats(ctx context.Context, concertID, zoneID uuid.UUID) (*entity.Seats, error)
}
