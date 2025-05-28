package concertrepo

import (
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"

	"github.com/kittipat1413/go-common/util/pointer"
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

type Concerts []Concert

func (cs Concerts) ToEntities() *entity.Concerts {
	var entities entity.Concerts
	for _, c := range cs {
		entities = append(entities, pointer.GetValue(c.ToEntity()))
	}
	return pointer.ToPointer(entities)
}
