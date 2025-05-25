//	@title			Ticket Reservation API
//	@version		1.0
//	@description	This is a ticket reservation system API.

//	@contact.name	Kittipat Poonyakariyakorn
//	@contact.email	k.poonyakariyakorn@gmail.com

//	@host		localhost:8080
//	@BasePath	/
//	@schemes	https http

// @securityDefinitions.basic	BasicAuth
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
package httproute

import (
	concertHandler "ticket-reservation/internal/api/http/handler/concert"
	healthHandler "ticket-reservation/internal/api/http/handler/healthcheck"
	"ticket-reservation/internal/api/http/middleware"
	"ticket-reservation/internal/config"

	"github.com/gin-gonic/gin"
)

type Router interface {
	RegisterRoutes(router *gin.Engine) // Register the routes for the application
}

type router struct {
	cfg                config.AppConfig                 // Configuration for the application
	Middleware         middleware.Middleware            // Middleware for handling requests
	HealthCheckHandler healthHandler.HealthCheckHandler // Handler for health check routes
	ConcertHandler     concertHandler.ConcertHandler    // Handler for concert routes
}

type Dependency struct {
	Middleware         middleware.Middleware
	HealthCheckHandler healthHandler.HealthCheckHandler
	ConcertHandler     concertHandler.ConcertHandler
}

// NewHTTPRoutes creates a new instance of Router with the provided configuration and dependencies
func NewHTTPRoutes(cfg config.AppConfig, dep Dependency) Router {
	return &router{
		cfg:                cfg,
		Middleware:         dep.Middleware,
		HealthCheckHandler: dep.HealthCheckHandler,
		ConcertHandler:     dep.ConcertHandler,
	}
}

// RegisterRoutes registers the routes for the application
func (r *router) RegisterRoutes(router *gin.Engine) {
	r.applyHealthCheckRoutes(router)
	r.applyConcertRoutes(router)
}

// applyHealthCheckRoutes applies the health check routes to the provided router
func (r *router) applyHealthCheckRoutes(router *gin.Engine) {
	healthRoute := router.Group("/health")
	{
		healthRoute.GET("/liveness", r.Middleware.BasicAuth(r.cfg.AdminAPIKey, r.cfg.AdminAPISecret), r.HealthCheckHandler.Liveness)
		healthRoute.GET("/readiness", r.Middleware.BasicAuth(r.cfg.AdminAPIKey, r.cfg.AdminAPISecret), r.HealthCheckHandler.Readiness)
	}
}

// applyConcertRoutes applies the concert routes to the provided router
func (r *router) applyConcertRoutes(router *gin.Engine) {
	concertRoute := router.Group("/concerts")
	{
		concertRoute.POST("/", r.ConcertHandler.CreateConcert)
	}
}
