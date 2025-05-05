package handler

import (
	"ticket-reservation/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type ReadinessResponse struct {
	Status string `json:"status" example:"OK"`
}

// @Summary Readiness
// @Description Check the readiness of the service
// @Tags HealthCheck
// @security BasicAuth
// @Produce json
// @Success 200 {object} ReadinessResponse "Success response"
// @Failure default {object} httpresponse.ErrorResponse "Default error response"
// @Router /health/readiness [get]
func (h *healthCheckHandler) Readiness(c *gin.Context) {
	ok, err := h.healthcheckUsecase.CheckReadiness(c.Request.Context())
	if err != nil || !ok {
		httpresponse.Error(c, err)
		return
	}
	httpresponse.Success(c, h.newReadinessResponse())
}

func (h *healthCheckHandler) newReadinessResponse() ReadinessResponse {
	return ReadinessResponse{
		Status: "OK",
	}
}
