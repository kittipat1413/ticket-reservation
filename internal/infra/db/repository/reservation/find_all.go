package reservationrepo

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/domain/repository"
	table "ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/table"

	"github.com/go-jet/jet/v2/postgres"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *reservationRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllReservationsFilter) (reservations *entity.Reservations, total int64, err error) {
	const errLocation = "[repository reservation/find_all FindAll] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	// Build WHERE conditions for filtering
	whereClauses := []postgres.BoolExpression{}
	if filter.SeatID != nil {
		whereClauses = append(whereClauses, table.Reservations.SeatID.EQ(postgres.UUID(*filter.SeatID)))
	}
	if filter.SessionID != nil {
		whereClauses = append(whereClauses, table.Reservations.SessionID.EQ(postgres.String(*filter.SessionID)))
	}
	if filter.Status != nil {
		whereClauses = append(whereClauses, table.Reservations.Status.EQ(postgres.String(filter.Status.String())))
	}

	// Get total count of reservations matching the filter
	countStmt := postgres.SELECT(
		postgres.COUNT(table.Reservations.ID).AS("total"),
	).FROM(table.Reservations)
	if len(whereClauses) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(whereClauses...))
	}

	countQuery, countArgs := countStmt.Sql()

	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while counting reservations", err.Error()))
	}

	// Get reservations with the same filter
	stmt := postgres.SELECT(
		table.Reservations.AllColumns,
	).FROM(table.Reservations)

	if len(whereClauses) > 0 {
		stmt = stmt.WHERE(postgres.AND(whereClauses...))
	}
	// Apply pagination
	if filter.Limit != nil {
		stmt = stmt.LIMIT(*filter.Limit)
	}
	if filter.Offset != nil {
		stmt = stmt.OFFSET(*filter.Offset)
	}

	query, args := stmt.Sql()

	var models Reservations
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while getting reservations", err.Error()))
	}

	reservations = models.ToEntities()
	return reservations, total, nil
}
