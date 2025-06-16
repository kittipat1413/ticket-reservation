package concertrepo

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

func (r *concertRepositoryImpl) FindOne(ctx context.Context, id uuid.UUID) (concert *entity.Concert, err error) {
	errLocation := "[repository concert/find_one FindOne] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	concertsTable := table.Concerts
	// SQL statement
	stmt := postgres.SELECT(
		concertsTable.AllColumns,
	).
		FROM(table.Concerts).
		WHERE(table.Concerts.ID.EQ(postgres.UUID(id)))

	query, args := stmt.Sql()

	var model Concert
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errsFramework.NewNotFoundError("concert not found", nil)
		}
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while querying concert by id", err.Error()))
	}

	concert = model.ToEntity()
	if concert == nil {
		return nil, errsFramework.NewInternalServerError("failed to convert concert model to entity", nil)
	}

	return
}
