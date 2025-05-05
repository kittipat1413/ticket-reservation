package concertrepo

import (
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"
)

type Concert struct {
	model.Concerts
}

func (c *Concert) ToEntity() *entity.Concert {
	return &entity.Concert{
		ID:        c.ID,
		Name:      c.Name,
		Venue:     c.Venue,
		Date:      c.Date,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
