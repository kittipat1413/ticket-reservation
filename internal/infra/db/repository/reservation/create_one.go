package reservationrepo

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"
	table "ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/table"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *reservationRepositoryImpl) CreateOne(ctx context.Context, input *entity.Reservation) (reservation *entity.Reservation, err error) {
	const errLocation = "[repository reservation/create_one CreateOne] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	reservationsTable := table.Reservations
	// SQL statement
	stmt := reservationsTable.INSERT(
		reservationsTable.AllColumns.Except(reservationsTable.DefaultColumns), // Exclude columns with default values
	).MODEL(model.Reservations{
		SeatID:     input.SeatID,
		SessionID:  input.SessionID,
		Status:     input.Status.String(),
		ReservedAt: input.ReservedAt,
		ExpiresAt:  input.ExpiresAt,
	}).RETURNING(reservationsTable.AllColumns)

	query, args := stmt.Sql()

	var model Reservation
	if err := r.execer.GetContext(ctx, &model, query, args...); err != nil {
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while creating reservation", err.Error()))
	}

	res := model.ToEntity()
	if res == nil {
		return nil, errsFramework.NewInternalServerError("failed to convert reservation model to entity", nil)
	}

	return res, nil
}
