package cache

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrSeatAlreadyLocked indicates that the seat is already locked by another process.
	ErrSeatAlreadyLocked = errors.New("seat already locked")
	// ErrSeatUnlockDenied indicates that the unlock operation was denied, likely due to a token mismatch.
	ErrSeatUnlockDenied = errors.New("seat unlock denied")
)

//go:generate mockgen -source=./seat_lock_repository.go -destination=./mocks/seat_lock_repository.go -package=cache_mocks
type SeatLocker interface {
	// LockSeat attempts to lock a specific seat for a concert in a given zone.
	// Returns ErrSeatAlreadyLocked if the seat is already locked by another process.
	LockSeat(ctx context.Context, concertID, zoneID, seatID, token string, ttl time.Duration) error
	// UnlockSeat releases the lock on a specific seat for a concert in a given zone.
	// Returns ErrSeatUnlockDenied if the unlock operation is denied, likely due to a token mismatch.
	UnlockSeat(ctx context.Context, concertID, zoneID, seatID, token string) error
}
