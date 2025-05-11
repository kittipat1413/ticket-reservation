package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidSeatStatus = fmt.Errorf("invalid seat status")
)

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusPending   SeatStatus = "pending"
	SeatStatusBooked    SeatStatus = "booked"
)

var seatStatusStringMapper = map[SeatStatus]string{
	SeatStatusAvailable: "available",
	SeatStatusPending:   "pending",
	SeatStatusBooked:    "booked",
}

func (s SeatStatus) String() string {
	return seatStatusStringMapper[s]
}

func (s SeatStatus) IsValid() bool {
	switch s {
	case SeatStatusAvailable, SeatStatusPending, SeatStatusBooked:
		return true
	default:
		return false
	}
}

// Parse parses a string into a SeatStatus. It returns an error if the string is not a valid SeatStatus.
func (s SeatStatus) Parse(status string) (SeatStatus, error) {
	seatStatus := SeatStatus(status)
	if !seatStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidSeatStatus, status)
	}
	return seatStatus, nil
}

type Seat struct {
	ID                uuid.UUID
	ZoneID            uuid.UUID
	SeatNumber        string
	Status            SeatStatus
	LockedUntil       *time.Time
	LockedBySessionID *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (s *Seat) IsLocked(now time.Time) bool {
	return s.Status == SeatStatusPending && s.LockedUntil != nil && now.Before(*s.LockedUntil)
}

func (s *Seat) IsAvailable(now time.Time) bool {
	return s.Status == SeatStatusAvailable || (s.Status == SeatStatusPending && s.LockedUntil != nil && now.After(*s.LockedUntil))
}
