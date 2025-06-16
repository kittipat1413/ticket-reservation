package errs

import (
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

type SeatAlreadyBookedError struct {
	*errsFramework.BaseError
}

// NewSeatAlreadyBookedError creates a new SeatAlreadyBookedError instance using the seat booked error code.
func NewSeatAlreadyBookedError() error {
	baseErr, err := errsFramework.NewBaseError(
		StatusCodeSeatBooked,
		"the seat is already booked.",
		nil,
	)
	if err != nil {
		return err
	}
	return &SeatAlreadyBookedError{
		BaseError: baseErr,
	}
}

type SeatLockedError struct {
	*errsFramework.BaseError
}

// NewSeatLockedError creates a new SeatLockedError instance using the seat locked error code.
func NewSeatLockedError() error {
	baseErr, err := errsFramework.NewBaseError(
		StatusCodeSeatLocked,
		"the seat is being reserved by another user.",
		nil,
	)
	if err != nil {
		return err
	}
	return &SeatLockedError{
		BaseError: baseErr,
	}
}
