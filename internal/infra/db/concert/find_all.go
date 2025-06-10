package concertrepo

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	table "ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/table"

	repository "ticket-reservation/internal/domain/repository"

	postgres "github.com/go-jet/jet/v2/postgres"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/util/pointer"
)

func (r *concertRepositoryImpl) FindAll(ctx context.Context, filter repository.FindAllConcertsFilter) (*entity.Concerts, int64, error) {
	const errLocation = "[repository concert/find_all FindAll] "
	var err error
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	// Build WHERE conditions for filtering
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

	// Get total count of concerts matching the filter
	countStmt := postgres.SELECT(
		postgres.COUNT(table.Concerts.ID).AS("total"),
	).FROM(table.Concerts)

	if len(whereClauses) > 0 {
		countStmt = countStmt.WHERE(postgres.AND(whereClauses...))
	}

	countQuery, countArgs := countStmt.Sql()

	var total int64
	if err := r.execer.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while counting concerts", err.Error()))
	}

	// Get concerts with the same filter
	stmt := postgres.SELECT(
		table.Concerts.AllColumns,
	).FROM(table.Concerts)

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
	// Apply sorting
	if filter.SortBy != nil {
		if filter.SortOrder == nil {
			filter.SortOrder = pointer.ToPointer(entity.SortOrderAsc) // Default to ascending if not provided
		}
		switch *filter.SortBy {
		case "name":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Concerts.Name.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Concerts.Name.ASC())
			}
		case "venue":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Concerts.Venue.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Concerts.Venue.ASC())
			}
		case "date":
			if *filter.SortOrder == entity.SortOrderDesc {
				stmt = stmt.ORDER_BY(table.Concerts.Date.DESC())
			} else {
				stmt = stmt.ORDER_BY(table.Concerts.Date.ASC())
			}
		}
	}

	query, args := stmt.Sql()

	var models Concerts
	if err := r.execer.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, 0, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while querying concerts", err.Error()))
	}

	concerts := models.ToEntities()
	return concerts, total, nil
}
