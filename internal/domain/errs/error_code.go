package errs

import (
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

const (
	StatusCodeSuccess                      = errsFramework.StatusCodeSuccess // 200000
	StatusCodeServiceCircuitBreakerTripped = "503001"                        // Client request blocked
)
