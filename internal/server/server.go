package server

import (
	"context"
	"fmt"
	"net/http"
	httproute "ticket-reservation/internal/api/http/route"
	"ticket-reservation/internal/config"
	"ticket-reservation/internal/util/httpresponse"
	"time"

	"github.com/gin-gonic/gin"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
	middlewareFramework "github.com/kittipat1413/go-common/framework/middleware/gin"
	"github.com/kittipat1413/go-common/framework/serverutils"
	"github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/util/pointer"

	infraDB "ticket-reservation/internal/infra/db"
	redis "ticket-reservation/internal/infra/redis"
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

	// Initialize Redis client
	redisClient := redis.NewClient(s.cfg)

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

	// Setup route dependencies
	deps, err := s.setupRouteDependencies(ctx, tracerProvider, appLogger, db, redisClient)
	if err != nil {
		return fmt.Errorf("failed to setup route dependencies: %w", err)
	}
	// Register application routes
	appRoutes := httproute.NewHTTPRoutes(s.cfg.App, deps)
	appRoutes.RegisterRoutes(router)

	// Prometheus metrics
	router.GET("/metrics", middlewareFramework.MetricsHandler())

	// Create http.Server
	httpServer := &http.Server{
		Addr:              s.cfg.Service.Port,
		Handler:           router,
		ReadHeaderTimeout: 15 * time.Second,
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
			"Redis client": func(ctx context.Context) error {
				return redisClient.Close()
			},
			// Other resources can be added here
		},
	)

	// Wait for shutdown to complete
	<-shutdownDoneCh
	appLogger.Info(ctx, "server shutdown complete", nil)
	return nil
}
