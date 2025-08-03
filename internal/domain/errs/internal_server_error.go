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
		"the service is currently unavailable. please try again later.",
		data,
	)
	if err != nil {
		return err
	}
	return &ServiceCircuitBreakerError{
		BaseError: baseErr,
	}
}

// As implements the error.As interface for ServiceCircuitBreakerError.
func (e *ServiceCircuitBreakerError) As(target interface{}) bool {
	if target == nil {
		return false
	}

	switch t := target.(type) {
	case **ServiceCircuitBreakerError:
		*t = e
		return true
	case *ServiceCircuitBreakerError:
		*t = *e
		return true
	default:
		return false
	}
}
