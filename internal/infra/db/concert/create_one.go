package concertrepo

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"
	table "ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/table"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *concertRepositoryImpl) CreateOne(ctx context.Context, input *entity.Concert) (concert *entity.Concert, err error) {
	errLocation := "[repository concert/create_one CreateOne] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	concertsTable := table.Concerts
	// SQL statement
	stmt := concertsTable.INSERT(
		concertsTable.AllColumns.Except(concertsTable.DefaultColumns), // Exclude columns with default values
	).MODEL(model.Concerts{
		Name:  input.Name,
		Date:  input.Date,
		Venue: input.Venue,
	}).RETURNING(concertsTable.AllColumns)

	query, args := stmt.Sql()

	var model Concert
	err = r.execer.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while creating concert", err.Error()))
	}

	concert = model.ToEntity()
	return
}
