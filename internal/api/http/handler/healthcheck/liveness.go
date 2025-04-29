package handler

import (
	"ticket-reservation/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type LivenessResponse struct {
	Status string `json:"status"`
}

// @Summary Liveness
// @Description Check the liveness of the service
// @Tags HealthCheck
// @security BasicAuth
// @Produce json
// @Success 200 {object} LivenessResponse "Success response"
// @Failure default {object} httpresponse.ErrorResponse "Default error response"
// @Router /health/liveness [get]
func (h *healthCheckHandler) Liveness(c *gin.Context) {
	httpresponse.Success(c, newLivenessResponse())
}

func newLivenessResponse() LivenessResponse {
	return LivenessResponse{
		Status: "OK",
	}
}
