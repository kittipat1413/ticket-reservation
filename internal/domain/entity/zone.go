package entity

import (
	"time"

	"github.com/google/uuid"
)

type Zone struct {
	ID          uuid.UUID
	ConcertID   uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Zones []Zone
