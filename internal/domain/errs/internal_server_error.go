package errs

import (
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

type ServiceCircuitBreakerError struct {
	*errsFramework.BaseError
}

// NewServiceCircuitBreakerError creates a new ServiceCircuitBreakerError instance using the service circuit breaker error code.
func NewServiceCircuitBreakerError(data interface{}) error {
	baseErr, err := errsFramework.NewBaseError(
		StatusCodeServiceCircuitBreakerTripped,
		"The service is currently unavailable. Please try again later.",
		data,
	)
	if err != nil {
		return err
	}
	return &ServiceCircuitBreakerError{
		BaseError: baseErr,
	}
}
