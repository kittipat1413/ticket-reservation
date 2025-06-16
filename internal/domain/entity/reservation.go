package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidReservationStatus = fmt.Errorf("invalid reservation status")
)

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusExpired   ReservationStatus = "expired"
)

var reservationStatusStringMapper = map[ReservationStatus]string{
	ReservationStatusPending:   "pending",
	ReservationStatusConfirmed: "confirmed",
	ReservationStatusExpired:   "expired",
}

func (s ReservationStatus) String() string {
	return reservationStatusStringMapper[s]
}
func (s ReservationStatus) IsValid() bool {
	switch s {
	case ReservationStatusPending, ReservationStatusConfirmed, ReservationStatusExpired:
		return true
	default:
		return false
	}
}

// Parse parses a string into a ReservationStatus. It returns an error if the string is not a valid ReservationStatus.
func (s ReservationStatus) Parse(status string) (ReservationStatus, error) {
	reservationStatus := ReservationStatus(status)
	if !reservationStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidReservationStatus, status)
	}
	return reservationStatus, nil
}

type Reservation struct {
	ID         uuid.UUID
	SeatID     uuid.UUID
	SessionID  string
	Status     ReservationStatus
	ReservedAt time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewReservation(seatID uuid.UUID, sessionID string, expiresAt time.Time) *Reservation {
	return &Reservation{
		ID:         uuid.New(),
		SeatID:     seatID,
		SessionID:  sessionID,
		Status:     ReservationStatusPending,
		ReservedAt: time.Now(),
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (r *Reservation) IsExpired(now time.Time) bool {
	return r.Status == ReservationStatusPending && now.After(r.ExpiresAt)
}

func (r *Reservation) CanPay(now time.Time) bool {
	return r.Status == ReservationStatusPending && now.Before(r.ExpiresAt)
}

type Reservations []Reservation
