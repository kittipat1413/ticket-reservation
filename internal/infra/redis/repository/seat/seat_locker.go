package seatrepo

import (
	"context"
	domaincache "ticket-reservation/internal/domain/cache"
	"time"

	"github.com/google/uuid"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	lockmanager "github.com/kittipat1413/go-common/framework/lockmanager"
)

type seatLocker struct {
	lockmanager lockmanager.LockManager
}

func NewSeatLockerRepository(lockmanager lockmanager.LockManager) domaincache.SeatLockerRepository {
	return &seatLocker{
		lockmanager: lockmanager,
	}
}

func (s *seatLocker) LockSeat(ctx context.Context, concertID, zoneID, seatID uuid.UUID, token string, ttl time.Duration) (err error) {
	const errLocation = "[repository seat/seat_locker LockSeat]"
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	key := getSeatLockKey(concertID, zoneID, seatID)

	_, err = s.lockmanager.Acquire(ctx, key, ttl, token)
	if err == lockmanager.ErrLockAlreadyTaken {
		return domaincache.ErrSeatAlreadyLocked
	}
	if err != nil {
		return errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to lock seat", err.Error()))
	}
	return nil
}

func (s *seatLocker) UnlockSeat(ctx context.Context, concertID, zoneID, seatID uuid.UUID, token string) (err error) {
	const errLocation = "[repository seat/seat_locker UnlockSeat]"
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	key := getSeatLockKey(concertID, zoneID, seatID)

	err = s.lockmanager.Release(ctx, key, token)
	if err == lockmanager.ErrUnlockNotPermitted {
		return domaincache.ErrSeatUnlockDenied
	}
	if err != nil {
		return errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to unlock seat", err.Error()))
	}
	return nil
}

func getSeatLockKey(concertID, zoneID, seatID uuid.UUID) string {
	return domaincache.GetSeatLockKey(concertID, zoneID, seatID)
}
