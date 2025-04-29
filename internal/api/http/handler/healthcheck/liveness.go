package handler

import (
	"ticket-reservation/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type LivenessResponse struct {
	Status string `json:"status"`
}

func (h *healthCheckHandler) Liveness(c *gin.Context) {
	httpresponse.Success(c, newLivenessResponse())
}

func newLivenessResponse() LivenessResponse {
	return LivenessResponse{
		Status: "OK",
	}
}
