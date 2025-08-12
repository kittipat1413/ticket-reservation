package server

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kittipat1413/go-common/framework/logger"
	"github.com/kittipat1413/go-common/framework/retry"
	"github.com/redis/go-redis/v9"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	infraDB "ticket-reservation/internal/infra/db"
	seatRedisRepo "ticket-reservation/internal/infra/redis/repository/seat"

	redsyncLocker "github.com/kittipat1413/go-common/framework/lockmanager/redsync"

	concertRepo "ticket-reservation/internal/infra/db/repository/concert"
	healthcheckRepo "ticket-reservation/internal/infra/db/repository/healthcheck"
	reservationRepo "ticket-reservation/internal/infra/db/repository/reservation"
	seatRepo "ticket-reservation/internal/infra/db/repository/seat"
	zonerepo "ticket-reservation/internal/infra/db/repository/zone"

	concertUsecase "ticket-reservation/internal/usecase/concert"
	healthcheckUsecase "ticket-reservation/internal/usecase/healthcheck"
	seatUsecase "ticket-reservation/internal/usecase/seat"

	"ticket-reservation/internal/api/http/middleware"
	httproute "ticket-reservation/internal/api/http/route"

	concertHandler "ticket-reservation/internal/api/http/handler/concert"
	healthcheckHandler "ticket-reservation/internal/api/http/handler/healthcheck"
	seatHandler "ticket-reservation/internal/api/http/handler/seat"
)

//nolint:unparam
func (s *Server) setupRouteDependencies(ctx context.Context, tracerProvider *sdktrace.TracerProvider, appLogger logger.Logger, dbConn *sqlx.DB, redisClient redis.UniversalClient) (httproute.Dependency, error) {
	// Redis lock manager
	lockmanager := redsyncLocker.NewRedsyncLockManager(redisClient)
	// Transactor factory
	transactorFactory := infraDB.NewSqlxTransactorFactory(dbConn)

	// Redis Repositories
	seatLockerRepo := seatRedisRepo.NewSeatLockerRepository(lockmanager)
	seatMapRepo := seatRedisRepo.NewSeatMapRepository(redisClient)

	// DB Repositories
	healthRepo := healthcheckRepo.NewHealthCheckRepository(dbConn)
	concertRepo := concertRepo.NewConcertRepository(dbConn)
	zoneRepo := zonerepo.NewZoneRepository(dbConn)
	seatRepo := seatRepo.NewSeatRepository(dbConn)
	reservationRepo := reservationRepo.NewReservationRepository(dbConn)

	// Query retrier
	queryBackoff, _ := retry.NewExponentialBackoffStrategy(500*time.Millisecond, 2.0, 5*time.Second)
	queryRetrier, _ := retry.NewRetrier(retry.Config{
		MaxAttempts: 3,
		Backoff:     queryBackoff,
	})

	// Usecases
	healthcheckUsecase := healthcheckUsecase.NewHealthCheckUsecase(queryRetrier, healthRepo)
	concertUsecase := concertUsecase.NewConcertUsecase(s.cfg.App, transactorFactory, concertRepo)
	seatUsecase := seatUsecase.NewSeatUsecase(s.cfg.App, concertRepo, zoneRepo, seatRepo, reservationRepo, transactorFactory, seatLockerRepo, seatMapRepo)

	// Application middleware
	appMiddleware := middleware.New()

	// Handlers
	healthHandler := healthcheckHandler.NewHealthCheckHandler(healthcheckUsecase)
	concertHandler := concertHandler.NewConcertHandler(s.cfg.App, concertUsecase)
	seatHandler := seatHandler.NewSeatHandler(s.cfg.App, seatUsecase)

	return httproute.Dependency{
		Middleware:         appMiddleware,
		HealthCheckHandler: healthHandler,
		ConcertHandler:     concertHandler,
		SeatHandler:        seatHandler,
	}, nil
}
