package usecase

import (
	"context"
	"errors"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	traceFramework "github.com/kittipat1413/go-common/framework/trace"
)

func (u *healthCheckUsecase) CheckReadiness(ctx context.Context) (ok bool, err error) {
	const errLocation = "[usecase healthcheck/check_readiness CheckReadiness] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	return traceFramework.TraceFunc(ctx, traceFramework.GetTracer("healthcheck.usecase"), func(ctx context.Context) (bool, error) {
		var ok bool
		err := u.retrier.ExecuteWithRetry(ctx, func(ctx context.Context) error {
			ok, err = u.healthcheckRepository.CheckDatabaseReadiness(ctx)
			return err
		}, func(attempt int, err error) bool {
			if errors.As(err, &errsFramework.DatabaseError{}) {
				return true // retry if database is not ready
			}
			return false
		})
		if err != nil {
			return false, errsFramework.WrapError(err, errsFramework.NewInternalServerError("service is not ready", nil))
		}
		return ok, nil
	})
}
