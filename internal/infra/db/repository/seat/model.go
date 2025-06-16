package seatrepo

import (
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"

	"github.com/kittipat1413/go-common/util/pointer"
)

type Seat struct {
	model.Seats
}

func (s *Seat) ToEntity() *entity.Seat {
	seatStatus, err := new(entity.SeatStatus).Parse(s.Status)
	if err != nil {
		return nil
	}
	return &entity.Seat{
		ID:                s.ID,
		ZoneID:            s.ZoneID,
		SeatNumber:        s.SeatNumber,
		Status:            seatStatus,
		LockedUntil:       s.LockedUntil,
		LockedBySessionID: s.LockedBySessionID,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

type Seats []Seat

func (ss Seats) ToEntities() *entity.Seats {
	seats := make(entity.Seats, 0, len(ss))
	for _, s := range ss {
		seat := s.ToEntity()
		if seat == nil {
			continue
		}
		seats = append(seats, pointer.GetValue(seat))
	}
	return pointer.ToPointer(seats)
}
