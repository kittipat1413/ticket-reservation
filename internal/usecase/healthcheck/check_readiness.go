package usecase

import (
	"context"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	traceFramework "github.com/kittipat1413/go-common/framework/trace"
)

func (u *healthCheckUsecase) CheckReadiness(ctx context.Context) (ok bool, err error) {
	const errLocation = "[usecase healthcheck/check_readiness CheckReadiness] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	return traceFramework.TraceFunc(ctx, traceFramework.GetTracer("healthcheck.usecase"), func(ctx context.Context) (bool, error) {
		ok, err := u.healthcheckRepository.CheckDatabaseReadiness(ctx)
		if err != nil {
			return false, errsFramework.WrapError(err, errsFramework.NewDatabaseError("database is not ready", nil))
		}
		return ok, nil
	})
}
