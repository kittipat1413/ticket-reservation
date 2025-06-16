package seatrepo

import (
	"context"
	"database/sql"
	"errors"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/domain/repository"
	table "ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/table"

	postgres "github.com/go-jet/jet/v2/postgres"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *seatRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateSeatInput) (seat *entity.Seat, err error) {
	const errLocation = "[repository seat/update_one UpdateOne] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	seatsTable := table.Seats

	var updateModel Seat
	columns := make(postgres.ColumnList, 0)

	// build the update model
	if input.Status != nil {
		updateModel.Status = input.Status.String()
		columns = append(columns, seatsTable.Status)
	}
	if input.LockedUntil != nil {
		updateModel.LockedUntil = input.LockedUntil
		columns = append(columns, seatsTable.LockedUntil)
	}
	if input.LockedBySessionID != nil {
		updateModel.LockedBySessionID = input.LockedBySessionID
		columns = append(columns, seatsTable.LockedBySessionID)
	}
	if len(columns) == 0 {
		return nil, errsFramework.NewBadRequestError("no fields provided to update", nil)
	}

	// SQL statement
	stmt := seatsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(seatsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(seatsTable.AllColumns)

	query, args := stmt.Sql()

	var model Seat
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errsFramework.NewNotFoundError("seat not found", nil)
		}
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while updating seat", err.Error()))
	}

	seat = model.ToEntity()
	if seat == nil {
		return nil, errsFramework.NewInternalServerError("failed to convert seat model to entity", nil)
	}

	return seat, nil
}
