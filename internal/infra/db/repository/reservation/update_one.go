package reservationrepo

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

func (r *reservationRepositoryImpl) UpdateOne(ctx context.Context, input repository.UpdateReservationInput) (reservation *entity.Reservation, err error) {
	const errLocation = "[repository reservation/update_one UpdateOne] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	reservationsTable := table.Reservations

	var updateModel Reservation
	columns := make(postgres.ColumnList, 0)

	// build the update model
	if input.Status != nil {
		updateModel.Status = input.Status.String()
		columns = append(columns, reservationsTable.Status)
	}
	if input.ExpiresAt != nil {
		updateModel.ExpiresAt = *input.ExpiresAt
		columns = append(columns, reservationsTable.ExpiresAt)
	}
	if len(columns) == 0 {
		return nil, errsFramework.NewBadRequestError("no fields provided to update", nil)
	}

	// SQL statement
	stmt := reservationsTable.
		UPDATE(columns).
		MODEL(updateModel).
		WHERE(reservationsTable.ID.EQ(postgres.UUID(input.ID))).
		RETURNING(reservationsTable.AllColumns)

	query, args := stmt.Sql()

	var model Reservation
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errsFramework.NewNotFoundError("reservation not found", nil)
		}
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while updating reservation", err.Error()))
	}

	reservation = model.ToEntity()
	if reservation == nil {
		return nil, errsFramework.NewInternalServerError("failed to convert reservation model to entity", nil)
	}

	return
}
