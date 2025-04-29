package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ticket-reservation/internal/api/http/middleware"
	httproute "ticket-reservation/internal/api/http/route"
	"ticket-reservation/internal/config"
	"ticket-reservation/internal/domain/errs"
	"ticket-reservation/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
	middlewareFramework "github.com/kittipat1413/go-common/framework/middleware/gin"
	"github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/util/pointer"
	"github.com/sony/gobreaker"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	healthcheckHandler "ticket-reservation/internal/api/http/handler/healthcheck"
	"ticket-reservation/internal/infra/db"
	healthcheckRepoImpl "ticket-reservation/internal/infra/db/healthcheck"
	healthcheckUsecase "ticket-reservation/internal/usecase/healthcheck"
)

type Server struct {
	cfg *config.Config
}

func New() *Server {
	cfg := config.MustConfigure()
	return &Server{cfg}
}

func (s *Server) Start() error {
	// Initialize context
	ctx := context.Background()
	// Initialize error framework (setup error prefix {prefix}-{error code})
	errsFramework.SetServicePrefix(s.cfg.ServiceErrPrefix())

	// Initialize Tracer Provider
	tracerProvider, err := trace.InitTracerProvider(ctx, s.cfg.ServiceName(), pointer.ToPointer(s.cfg.OtelExporter()), trace.ExporterGRPC)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer provider: %w", err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	// Initialize logger
	logConfig := logger.Config{
		Level:       logger.INFO,
		Environment: s.cfg.ServiceEnv(),
		ServiceName: s.cfg.ServiceName(),
	}
	appLogger, err := logger.NewLogger(logConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database connection
	db, err := db.Connect(s.cfg, tracerProvider)
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
	healthCheckRepository := healthcheckRepoImpl.NewHealthCheckRepository(db)
	// Usecase
	healthcheckUsecase := healthcheckUsecase.NewHealthCheckUsecase(healthCheckRepository)
	// Application middleware
	appMiddleware := middleware.New()
	// Handler
	healthcheckHandler := healthcheckHandler.NewHealthCheckHandler(healthcheckUsecase)

	// Register application routes
	appRoutes := httproute.NewHTTPRoutes(s.cfg, httproute.Dependency{
		Middleware:         appMiddleware,
		HealthCheckHandler: healthcheckHandler,
	})
	appRoutes.RegisterRoutes(router)

	/*
		// !!!!!!!! gracefulshutdown here !!!!!!!!
	*/

	appLogger.Info(ctx, "server started", logger.Fields{"service_name": s.cfg.ServiceName(), "service_env": s.cfg.ServiceEnv(), "service_port": s.cfg.ServicePort()})
	return router.Run(s.cfg.ServicePort())
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
				Name: fmt.Sprintf("%s-circuit-breaker", s.cfg.ServiceName()),
			}),
			middlewareFramework.WithCircuitBreakerFilter(isNotHealthCheck),
			middlewareFramework.WithCircuitBreakerErrorHandler(s.circuitBreakerHandler),
		),
		// Middleware to insert config into request context
		s.insertCfg(),
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

func (s *Server) insertCfg() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = config.NewRequest(c.Request, s.cfg)
	}
}
