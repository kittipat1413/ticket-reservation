package entity

import (
	"time"

	"github.com/google/uuid"
)

type Concert struct {
	ID        uuid.UUID
	Name      string
	Venue     string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
