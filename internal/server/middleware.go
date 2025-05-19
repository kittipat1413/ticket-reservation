package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"

	"ticket-reservation/internal/domain/errs"
	"ticket-reservation/internal/util/httpresponse"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
	middlewareFramework "github.com/kittipat1413/go-common/framework/middleware/gin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func (s *Server) setupMiddlewares(appLogger logger.Logger, tracerProvider *sdktrace.TracerProvider) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// Recovery middleware
		middlewareFramework.Recovery(
			middlewareFramework.WithRecoveryLogger(appLogger),
			middlewareFramework.WithRecoveryHandler(s.recoveryHandler),
		),
		// Trace middleware
		middlewareFramework.Trace(
			middlewareFramework.WithTracerProvider(tracerProvider),
			middlewareFramework.WithTraceFilter(isNotHealthCheck),
		),
		// RequestID middleware
		middlewareFramework.RequestID(),
		// RequestLogger middleware
		middlewareFramework.RequestLogger(
			middlewareFramework.WithRequestLogger(appLogger),
			middlewareFramework.WithRequestLoggerFilter(isNotHealthCheck),
		),
		// CircuitBreaker middleware
		middlewareFramework.CircuitBreaker(
			middlewareFramework.WithCircuitBreakerSettings(gobreaker.Settings{
				Name: fmt.Sprintf("%s-circuit-breaker", s.cfg.Service.Name),
			}),
			middlewareFramework.WithCircuitBreakerFilter(isNotHealthCheck),
			middlewareFramework.WithCircuitBreakerErrorHandler(s.circuitBreakerHandler),
		),
	}
}

func isNotHealthCheck(req *http.Request) bool {
	path := req.URL.Path
	return path != "/health/liveness" && path != "/health/readiness"
}

func (s *Server) recoveryHandler(c *gin.Context, err interface{}) {
	httpresponse.Error(c, errsFramework.NewInternalServerError("unexpected server error occurred", map[string]interface{}{
		"panic": err,
	}))
}

func (s *Server) circuitBreakerHandler(c *gin.Context) {
	httpresponse.Error(c, errs.NewServiceCircuitBreakerError(nil))
}
