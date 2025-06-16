package reservationrepo

import (
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"

	"github.com/kittipat1413/go-common/util/pointer"
)

type Reservation struct {
	model.Reservations
}

func (r *Reservation) ToEntity() *entity.Reservation {
	reservationStatus, err := new(entity.ReservationStatus).Parse(r.Status)
	if err != nil {
		return nil
	}
	return &entity.Reservation{
		ID:         r.ID,
		SeatID:     r.SeatID,
		SessionID:  r.SessionID,
		Status:     reservationStatus,
		ReservedAt: r.ReservedAt,
		ExpiresAt:  r.ExpiresAt,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

type Reservations []Reservation

func (rs Reservations) ToEntities() *entity.Reservations {
	reservations := make(entity.Reservations, 0, len(rs))
	for _, r := range rs {
		reservation := r.ToEntity()
		if reservation == nil {
			continue
		}
		reservations = append(reservations, pointer.GetValue(reservation))
	}
	return pointer.ToPointer(reservations)
}
