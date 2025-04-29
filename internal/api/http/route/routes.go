// @title Ticket Reservation API
// @version 1.0
// @description This is a ticket reservation system API.

// @contact.name Kittipat Poonyakariyakorn
// @contact.email k.poonyakariyakorn@gmail.com

// @host localhost:8080
// @BasePath /
// @schemes https http

// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package httproute

import (
	healthHandler "ticket-reservation/internal/api/http/handler/healthcheck"
	"ticket-reservation/internal/api/http/middleware"
	"ticket-reservation/internal/config"

	"github.com/gin-gonic/gin"
)

type Router interface {
	RegisterRoutes(router *gin.Engine) // Register the routes for the application
}

type router struct {
	cfg                *config.Config                   // Configuration for the application
	Middleware         middleware.Middleware            // Middleware for handling requests
	HealthCheckHandler healthHandler.HealthCheckHandler // Handler for health check routes
}

type Dependency struct {
	Middleware         middleware.Middleware
	HealthCheckHandler healthHandler.HealthCheckHandler
}

// NewHTTPRoutes creates a new instance of Router with the provided configuration and dependencies
func NewHTTPRoutes(cfg *config.Config, dep Dependency) Router {
	return &router{
		cfg:                cfg,
		Middleware:         dep.Middleware,
		HealthCheckHandler: dep.HealthCheckHandler,
	}
}

// RegisterRoutes registers the routes for the application
func (r *router) RegisterRoutes(router *gin.Engine) {
	r.applyHealthCheckRoutes(router)
}

// applyHealthCheckRoutes applies the health check routes to the provided router
func (r *router) applyHealthCheckRoutes(router *gin.Engine) {
	healthRoute := router.Group("/health")
	{
		healthRoute.GET("/liveness", r.Middleware.BasicAuth(r.cfg.AdminApiKey(), r.cfg.AdminApiSecret()), r.HealthCheckHandler.Liveness)
		healthRoute.GET("/readiness", r.Middleware.BasicAuth(r.cfg.AdminApiKey(), r.cfg.AdminApiSecret()), r.HealthCheckHandler.Readiness)
	}
}
