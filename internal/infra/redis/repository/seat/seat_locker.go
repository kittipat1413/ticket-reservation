package seatrepo

import (
	"context"
	domaincache "ticket-reservation/internal/domain/cache"
	"time"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	lockmanager "github.com/kittipat1413/go-common/framework/lockmanager"
)

type seatLocker struct {
	lockmanager lockmanager.LockManager
}

func NewSeatLocker(lm lockmanager.LockManager) domaincache.SeatLocker {
	return &seatLocker{lockmanager: lm}
}

func (s *seatLocker) LockSeat(ctx context.Context, concertID, zoneID, seatID, token string, ttl time.Duration) error {
	key := getSeatLockKey(concertID, zoneID, seatID)

	_, err := s.lockmanager.Acquire(ctx, key, ttl, token)
	if err == lockmanager.ErrLockAlreadyTaken {
		return domaincache.ErrSeatAlreadyLocked
	}
	if err != nil {
		return errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to lock seat", err.Error()))
	}
	return nil
}

func (s *seatLocker) UnlockSeat(ctx context.Context, concertID, zoneID, seatID, token string) error {
	key := getSeatLockKey(concertID, zoneID, seatID)

	err := s.lockmanager.Release(ctx, key, token)
	if err == lockmanager.ErrUnlockNotPermitted {
		return domaincache.ErrSeatUnlockDenied
	}
	if err != nil {
		return errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to unlock seat", err.Error()))
	}
	return nil
}

func getSeatLockKey(concertID, zoneID, seatID string) string {
	return domaincache.GetSeatLockKey(concertID, zoneID, seatID)
}
