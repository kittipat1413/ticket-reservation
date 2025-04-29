package handler

import (
	healthcheckUsecase "ticket-reservation/internal/usecase/healthcheck"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler interface {
	Liveness(c *gin.Context)
	Readiness(c *gin.Context)
}

type healthCheckHandler struct {
	healthcheckUsecase healthcheckUsecase.HealthCheckUsecase
}

func NewHealthCheckHandler(healthcheckUsecase healthcheckUsecase.HealthCheckUsecase) HealthCheckHandler {
	return &healthCheckHandler{
		healthcheckUsecase: healthcheckUsecase,
	}
}
