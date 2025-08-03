package usecase

import (
	"context"
	"errors"
	"ticket-reservation/internal/domain/cache"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/domain/repository"
	"time"

	"ticket-reservation/internal/domain/errs"

	"github.com/google/uuid"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
	traceFramework "github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/framework/validator"
	"github.com/kittipat1413/go-common/util/pointer"
)

type ReserveSeatInput struct {
	ConcertID string `json:"concert_id" validate:"required,uuid4"`
	ZoneID    string `json:"zone_id" validate:"required,uuid4"`
	SeatID    string `json:"seat_id" validate:"required,uuid4"`
	SessionID string `json:"session_id" validate:"required"`
}

func (u *seatUsecase) ReserveSeat(ctx context.Context, input ReserveSeatInput) (reservation *entity.Reservation, err error) {
	const errLocation = "[usecase seat/reserve_seat ReserveSeat] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	return traceFramework.TraceFunc(ctx, traceFramework.GetTracer("seat.usecase"), func(ctx context.Context) (*entity.Reservation, error) {
		requestTime := time.Now()

		// Create a new validator instance
		vInstance, err := validator.NewValidator(
			validator.WithTagNameFunc(validator.JSONTagNameFunc),
		)
		if err != nil {
			err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to create validator", nil))
			return nil, err
		}

		// Validate Input
		if err := vInstance.Struct(input); err != nil {
			err = errsFramework.WrapError(err, errsFramework.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
			return nil, err
		}

		var (
			concertID uuid.UUID
			zoneID    uuid.UUID
			seatID    uuid.UUID
		)
		concertID, err = uuid.Parse(input.ConcertID)
		if err != nil {
			err = errsFramework.WrapError(err, errsFramework.NewBadRequestError("invalid concert ID", nil))
			return nil, err
		}
		zoneID, err = uuid.Parse(input.ZoneID)
		if err != nil {
			err = errsFramework.WrapError(err, errsFramework.NewBadRequestError("invalid zone ID", nil))
			return nil, err
		}
		seatID, err = uuid.Parse(input.SeatID)
		if err != nil {
			err = errsFramework.WrapError(err, errsFramework.NewBadRequestError("invalid seat ID", nil))
			return nil, err
		}

		// Find concert by ID and check if it has already passed
		concert, err := u.concertRepository.FindOne(ctx, concertID)
		if err != nil {
			if !errors.As(err, &errsFramework.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
				err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to find concert by ID", nil))
				return nil, err
			}
			return nil, err // Return the NotFoundError directly
		}
		if concert.Date.Before(requestTime) {
			err = errsFramework.WrapError(err, errsFramework.NewConflictError("the concert has already passed", nil))
			return nil, err
		}

		// Find zone by ID and check if it belongs to the concert
		zone, err := u.zoneRepository.FindOne(ctx, zoneID)
		if err != nil {
			if !errors.As(err, &errsFramework.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
				err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to find zone by ID", nil))
				return nil, err
			}
			return nil, err // Return the NotFoundError directly
		}
		if zone.ConcertID != concertID {
			err = errsFramework.NewBadRequestError("the zone does not belong to the specified concert", nil)
			return nil, err
		}

		// Attempt to lock the seat
		err = u.seatLocker.LockSeat(ctx, input.ConcertID, input.ZoneID, input.SeatID, input.SessionID, u.appConfig.SeatLockTTL)
		if err != nil && errors.Is(err, cache.ErrSeatAlreadyLocked) {
			// If the seat is already locked, return an error
			err = errsFramework.WrapError(err, errs.NewSeatLockedError())
			return nil, err
		}
		defer func() {
			if err != nil {
				// If any error occurs, unlock the seat
				// This ensures that the seat lock is released if the operation fails
				unlockErr := u.seatLocker.UnlockSeat(ctx, input.ConcertID, input.ZoneID, input.SeatID, input.SessionID)
				if unlockErr != nil {
					// Log the error but do not return it, as the main error has already been handled
					logger.FromContext(ctx).Error(ctx, "failed to unlock seat after error", unlockErr, logger.Fields{
						"concert_id": input.ConcertID,
						"zone_id":    input.ZoneID,
						"seat_id":    input.SeatID,
						"session_id": input.SessionID,
					})
				}
			}
		}()

		// Start a transaction for database operations
		tx, err := u.transactorFactory.CreateSqlxTransactor(ctx)
		if err != nil {
			err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to create transaction", nil))
			return nil, err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}()

		// Get a seat with explicit row locking
		seat, err := u.seatRepository.WithTx(tx.DB()).FindOne(ctx, seatID)
		if err != nil {
			if !errors.As(err, &errsFramework.NotFoundError{}) { // If the error is not a NotFoundError, wrap it as an internal server error
				err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to find seat by ID", nil))
				return nil, err
			}
			return nil, err // Return the NotFoundError directly
		}
		// Check if the seat belongs to the specified zone
		if seat.ZoneID != zoneID {
			err = errsFramework.NewBadRequestError("the seat does not belong to the specified zone", nil)
			return nil, err
		}
		// Check if the seat is already booked
		if seat.IsBooked() {
			err = errs.NewSeatAlreadyBookedError()
			return nil, err
		}
		// Check if the seat is pending and lock is not expired
		if seat.IsLocked(requestTime) {
			// Check if the seat is locked by another session
			if seat.LockedBySessionID != nil && pointer.GetValue(seat.LockedBySessionID) != input.SessionID {
				err = errs.NewSeatLockedError()
				return nil, err
			}
		}

		// Update seat status in database
		seat, err = u.seatRepository.WithTx(tx.DB()).UpdateOne(ctx, repository.UpdateSeatInput{
			ID:                seat.ID,
			Status:            pointer.ToPointer(entity.SeatStatusPending),
			LockedBySessionID: pointer.ToPointer(input.SessionID),
			LockedUntil:       pointer.ToPointer(requestTime.Add(u.appConfig.SeatLockTTL)),
		})
		if err != nil {
			err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to update seat status", nil))
			return nil, err
		}

		var (
			reservation *entity.Reservation
		)

		// Find existing reservations for the seat
		existingReservations, _, err := u.reservationRepository.WithTx(tx.DB()).FindAll(ctx, repository.FindAllReservationsFilter{
			SeatID:    pointer.ToPointer(seat.ID),
			SessionID: pointer.ToPointer(input.SessionID),
			Status:    pointer.ToPointer(entity.ReservationStatusPending),
		})
		for _, existingReservation := range pointer.GetValue(existingReservations) {
			if existingReservation.CanPay(requestTime) {
				// Extend the expiration time of the existing reservation
				reservation, err = u.reservationRepository.WithTx(tx.DB()).UpdateOne(ctx, repository.UpdateReservationInput{
					ID:        existingReservation.ID,
					ExpiresAt: seat.LockedUntil,
				})
				if err != nil {
					err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to update existing reservation", nil))
					return nil, err
				}
			} else {
				// Mark the existing reservation as expired
				_, err = u.reservationRepository.WithTx(tx.DB()).UpdateOne(ctx, repository.UpdateReservationInput{
					ID:     existingReservation.ID,
					Status: pointer.ToPointer(entity.ReservationStatusExpired),
				})
				if err != nil {
					err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to update existing reservation", nil))
					return nil, err
				}
			}
		}

		// If no existing reservation found, create a new one
		if reservation == nil {
			reservation, err = u.reservationRepository.WithTx(tx.DB()).CreateOne(ctx, entity.NewReservation(seat.ID, input.SessionID, *seat.LockedUntil))
			if err != nil {
				err = errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to create reservation", nil))
				return nil, err
			}
		}

		/*
		 TODO: update redis cache for seat map
		 	HSET seat_map:concert:{cid}:zone:{zid} A5 "pending"
		*/

		return reservation, nil
	})
}
