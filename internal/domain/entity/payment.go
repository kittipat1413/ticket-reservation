package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidPaymentStatus = fmt.Errorf("invalid payment status")
)

type PaymentStatus string

const (
	PaymentStatusInitiated PaymentStatus = "initiated"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
)

var paymentStatusStringMapper = map[PaymentStatus]string{
	PaymentStatusInitiated: "initiated",
	PaymentStatusPaid:      "paid",
	PaymentStatusFailed:    "failed",
}

func (s PaymentStatus) String() string {
	return paymentStatusStringMapper[s]
}

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusInitiated, PaymentStatusPaid, PaymentStatusFailed:
		return true
	default:
		return false
	}
}

// Parse parses a string into a PaymentStatus. It returns an error if the string is not a valid PaymentStatus.
func (s PaymentStatus) Parse(status string) (PaymentStatus, error) {
	paymentStatus := PaymentStatus(status)
	if !paymentStatus.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidPaymentStatus, status)
	}
	return paymentStatus, nil
}

type Payment struct {
	ID            uuid.UUID
	ReservationID uuid.UUID
	Status        PaymentStatus
	Amount        *decimal.Decimal
	PaidAt        *time.Time
	PaymentMethod *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (p *Payment) IsPaid() bool {
	return p.Status == PaymentStatusPaid
}
