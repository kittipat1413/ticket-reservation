package handler

import (
	"ticket-reservation/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type livenessResponse struct {
	Status string `json:"status" example:"OK"`
}

// @Summary		Liveness
// @Description	Check the liveness of the service
// @Tags			HealthCheck
// @security		BasicAuth
// @Produce		json
// @Success		200		{object}	httpresponse.SuccessResponse{data=livenessResponse,metadata=nil}	"Success response"
// @Failure		default	{object}	httpresponse.ErrorResponse{data=nil}								"Default error response"
// @Router			/health/liveness [get]
func (h *healthCheckHandler) Liveness(c *gin.Context) {
	httpresponse.Success(c, h.newLivenessResponse())
}

func (h *healthCheckHandler) newLivenessResponse() livenessResponse {
	return livenessResponse{
		Status: "OK",
	}
}
