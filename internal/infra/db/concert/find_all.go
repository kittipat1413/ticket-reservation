package concertrepo

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	table "ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/table"

	repository "ticket-reservation/internal/domain/repository"

	postgres "github.com/go-jet/jet/v2/postgres"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *concertRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllConcertsFilter) (*entity.Concerts, error) {
	const errLocation = "[repository concert/find_all FindAll] "
	var err error
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	stmt := postgres.SELECT(
		table.Concerts.AllColumns,
	).
		FROM(table.Concerts)

	// Apply optional filters
	whereClauses := []postgres.BoolExpression{}

	if filter.StartDate != nil {
		whereClauses = append(whereClauses, table.Concerts.Date.GT_EQ(postgres.TimestampzT(*filter.StartDate)))
	}
	if filter.EndDate != nil {
		whereClauses = append(whereClauses, table.Concerts.Date.LT_EQ(postgres.TimestampzT(*filter.EndDate)))
	}
	if filter.Venue != nil && *filter.Venue != "" {
		whereClauses = append(whereClauses, table.Concerts.Venue.LIKE(postgres.String("%"+*filter.Venue+"%")))
	}
	if len(whereClauses) > 0 {
		stmt = stmt.WHERE(postgres.AND(whereClauses...))
	}

	query, args := stmt.Sql()

	var models Concerts
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while querying concerts", err.Error()))
	}

	concerts := models.ToEntities()
	return concerts, nil
}
