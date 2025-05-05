package server

import (
	"context"
	"fmt"
	"net/http"
	"ticket-reservation/internal/api/http/middleware"
	httproute "ticket-reservation/internal/api/http/route"
	"ticket-reservation/internal/config"
	"ticket-reservation/internal/domain/errs"
	"ticket-reservation/internal/util/httpresponse"
	"time"

	"github.com/gin-gonic/gin"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
	middlewareFramework "github.com/kittipat1413/go-common/framework/middleware/gin"
	"github.com/kittipat1413/go-common/framework/serverutils"
	"github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/util/pointer"
	"github.com/sony/gobreaker"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	infraDB "ticket-reservation/internal/infra/db"

	concertHandler "ticket-reservation/internal/api/http/handler/concert"
	healthcheckHandler "ticket-reservation/internal/api/http/handler/healthcheck"
	concertRepoImpl "ticket-reservation/internal/infra/db/concert"
	healthcheckRepoImpl "ticket-reservation/internal/infra/db/healthcheck"
	concertUsecase "ticket-reservation/internal/usecase/concert"
	healthcheckUsecase "ticket-reservation/internal/usecase/healthcheck"
)

type Server struct {
	cfg *config.Config
}

func New() *Server {
	return &Server{
		cfg: config.MustConfigure(),
	}
}

func (s *Server) Start() error {
	// Initialize context
	ctx := context.Background()

	// Initialize error framework (setup error prefix {prefix}-{error code})
	errsFramework.SetServicePrefix(s.cfg.Service.ErrorPrefix)

	// Initialize Tracer Provider
	tracerProvider, err := trace.InitTracerProvider(
		ctx,
		s.cfg.Service.Name,
		pointer.ToPointer(s.cfg.Service.OtelExporter),
		trace.ExporterGRPC,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer provider: %w", err)
	}

	// Initialize logger
	logConfig := logger.Config{
		Level:       logger.INFO,
		ServiceName: s.cfg.Service.Name,
		Environment: s.cfg.Service.Env,
	}
	appLogger, err := logger.NewLogger(logConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	// Set default logger config
	if err = logger.SetDefaultLoggerConfig(logConfig); err != nil {
		return fmt.Errorf("failed to set default logger config: %w", err)
	}

	// Initialize database connection
	db, err := infraDB.Connect(s.cfg, tracerProvider)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// initialize gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Setup middlewares
	middlewares := s.setupMiddlewares(appLogger, tracerProvider)
	// Apply middlewares
	router.Use(middlewares...)

	// Setup not found handler
	router.NoRoute(func(c *gin.Context) {
		httpresponse.Error(c, errsFramework.NewNotFoundError("the requested endpoint is not registered", nil))
	})

	// Repository
	transactorFactory := infraDB.NewSqlxTransactorFactory(db)
	healthCheckRepository := healthcheckRepoImpl.NewHealthCheckRepository(db)
	concertRepository := concertRepoImpl.NewConcertRepository(db)
	// Usecase
	healthcheckUsecase := healthcheckUsecase.NewHealthCheckUsecase(healthCheckRepository)
	concertUsecase := concertUsecase.NewConcertUsecase(s.cfg.App, transactorFactory, concertRepository)
	// Application middleware
	appMiddleware := middleware.New()
	// Handler
	healthcheckHandler := healthcheckHandler.NewHealthCheckHandler(healthcheckUsecase)
	concertHandler := concertHandler.NewConcertHandler(s.cfg.App, concertUsecase)

	// Register application routes
	appRoutes := httproute.NewHTTPRoutes(s.cfg.App, httproute.Dependency{
		Middleware:         appMiddleware,
		HealthCheckHandler: healthcheckHandler,
		ConcertHandler:     concertHandler,
	})
	appRoutes.RegisterRoutes(router)

	// Create http.Server
	httpServer := &http.Server{
		Addr:    s.cfg.Service.Port,
		Handler: router,
	}

	errCh := make(chan error, 1)

	// Run server in goroutine
	go func() {
		appLogger.Info(ctx, "server started", logger.Fields{
			"service_name": s.cfg.Service.Name,
			"service_env":  s.cfg.Service.Env,
			"service_port": s.cfg.Service.Port,
		})
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	// Wait for shutdown signal
	shutdownDoneCh := serverutils.GracefulShutdownSystem(
		ctx,
		appLogger,
		errCh,
		30*time.Second,
		map[string]serverutils.ShutdownOperation{
			"HTTP Server": func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
			"Tracer Provider": func(ctx context.Context) error {
				return tracerProvider.Shutdown(ctx)
			},
			"Database connection": func(ctx context.Context) error {
				return db.Close()
			},
			// Other resources can be added here (e.g., Redis, etc.)
		},
	)

	// Wait for shutdown to complete
	<-shutdownDoneCh
	appLogger.Info(ctx, "server shutdown complete", nil)
	return nil
}

func (s *Server) setupMiddlewares(appLogger logger.Logger, tracerProvider *sdktrace.TracerProvider) []gin.HandlerFunc {
	var middlewares = []gin.HandlerFunc{
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

	return middlewares
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
