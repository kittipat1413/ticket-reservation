package server

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kittipat1413/go-common/framework/logger"
	"github.com/kittipat1413/go-common/framework/retry"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	concertHandler "ticket-reservation/internal/api/http/handler/concert"
	healthcheckHandler "ticket-reservation/internal/api/http/handler/healthcheck"
	"ticket-reservation/internal/api/http/middleware"
	httproute "ticket-reservation/internal/api/http/route"

	"ticket-reservation/internal/infra/db"
	concertRepo "ticket-reservation/internal/infra/db/concert"
	healthcheckRepo "ticket-reservation/internal/infra/db/healthcheck"
	concertUsecase "ticket-reservation/internal/usecase/concert"
	healthcheckUsecase "ticket-reservation/internal/usecase/healthcheck"
)

//nolint:unparam
func (s *Server) setupRouteDependencies(ctx context.Context, tracerProvider *sdktrace.TracerProvider, appLogger logger.Logger, dbConn *sqlx.DB) (httproute.Dependency, error) {
	// Repositories
	transactorFactory := db.NewSqlxTransactorFactory(dbConn)
	healthRepo := healthcheckRepo.NewHealthCheckRepository(dbConn)
	concertRepo := concertRepo.NewConcertRepository(dbConn)

	// Query retrier
	queryBackoff, _ := retry.NewExponentialBackoffStrategy(500*time.Millisecond, 2.0, 5*time.Second)
	queryRetrier, _ := retry.NewRetrier(retry.Config{
		MaxAttempts: 3,
		Backoff:     queryBackoff,
	})

	// Usecases
	healthcheckUsecase := healthcheckUsecase.NewHealthCheckUsecase(queryRetrier, healthRepo)
	concertUsecase := concertUsecase.NewConcertUsecase(s.cfg.App, transactorFactory, concertRepo)

	// Application middleware
	appMiddleware := middleware.New()

	// Handlers
	healthHandler := healthcheckHandler.NewHealthCheckHandler(healthcheckUsecase)
	concertHandler := concertHandler.NewConcertHandler(s.cfg.App, concertUsecase)

	return httproute.Dependency{
		Middleware:         appMiddleware,
		HealthCheckHandler: healthHandler,
		ConcertHandler:     concertHandler,
	}, nil
}
