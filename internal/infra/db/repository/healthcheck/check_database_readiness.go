package healthcheckrepo

import (
	"context"

	postgres "github.com/go-jet/jet/v2/postgres"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *healthCheckRepositoryImpl) CheckDatabaseReadiness(ctx context.Context) (ok bool, err error) {
	const errLocation = "[repository healthcheck/check_database_readiness CheckDatabaseReadiness] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	stmt := postgres.RawStatement(`SELECT 1=1`)

	query, args := stmt.Sql()

	err = r.db.GetContext(ctx, &ok, query, args...)
	if err != nil {
		return false, errsFramework.WrapError(err, errsFramework.NewDatabaseError("error while checking database readiness", err.Error()))
	}
	return true, nil
}
