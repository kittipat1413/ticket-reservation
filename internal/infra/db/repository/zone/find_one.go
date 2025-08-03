package zonerepo

import (
	"context"
	"database/sql"
	"errors"
	"ticket-reservation/internal/domain/entity"
	table "ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/table"

	postgres "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *zoneRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (zone *entity.Zone, err error) {
	const errLocation = "[repository zone/find_one FindOne] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	zonesTable := table.Zones
	// SQL statement
	stmt := postgres.SELECT(
		zonesTable.AllColumns,
	).FROM(
		zonesTable,
	).WHERE(
		zonesTable.ID.EQ(postgres.UUID(id)),
	).FOR(
		postgres.UPDATE(),
	)

	query, args := stmt.Sql()

	var model Zone
	if err := r.execer.GetContext(ctx, &model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errsFramework.NewNotFoundError("zone not found", nil)
		}
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while getting zone", err.Error()))
	}

	zone = model.ToEntity()
	if zone == nil {
		return nil, errsFramework.NewInternalServerError("failed to convert zone model to entity", nil)
	}

	return
}
