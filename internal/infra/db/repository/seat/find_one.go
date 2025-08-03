package seatrepo

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

func (r *seatRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (seat *entity.Seat, err error) {
	const errLocation = "[repository seat/find_one FindOne] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	seatsTable := table.Seats
	// SQL statement
	stmt := postgres.SELECT(
		seatsTable.AllColumns,
	).FROM(
		seatsTable,
	).WHERE(
		seatsTable.ID.EQ(postgres.UUID(id)),
	).FOR(
		postgres.UPDATE(),
	)

	query, args := stmt.Sql()

	var model Seat
	if err := r.execer.GetContext(ctx, &model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errsFramework.NewNotFoundError("seat not found", nil)
		}
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while getting seat", err.Error()))
	}

	seat = model.ToEntity()
	if seat == nil {
		return nil, errsFramework.NewInternalServerError("failed to convert seat model to entity", nil)
	}

	return
}
