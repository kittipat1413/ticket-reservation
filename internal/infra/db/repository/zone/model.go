package zonerepo

import (
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"

	"github.com/kittipat1413/go-common/util/pointer"
)

type Zone struct {
	model.Zones
}

func (z *Zone) ToEntity() *entity.Zone {
	return &entity.Zone{
		ID:          z.ID,
		ConcertID:   z.ConcertID,
		Name:        z.Name,
		Description: z.Description,
		CreatedAt:   z.CreatedAt,
		UpdatedAt:   z.UpdatedAt,
	}
}

type Zones []Zone

func (zs Zones) ToEntities() *entity.Zones {
	zones := make(entity.Zones, 0, len(zs))
	for _, z := range zs {
		zone := z.ToEntity()
		if zone == nil {
			continue
		}
		zones = append(zones, pointer.GetValue(zone))
	}
	return pointer.ToPointer(zones)
}
