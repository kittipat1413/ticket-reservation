package handler

import (
	"ticket-reservation/internal/util/httpresponse"

	"github.com/gin-gonic/gin"
)

type ReadinessResponse struct {
	Status string `json:"status"`
}

func (h *healthCheckHandler) Readiness(c *gin.Context) {
	ok, err := h.healthcheckUsecase.CheckReadiness(c.Request.Context())
	if err != nil || !ok {
		httpresponse.Error(c, err)
		return
	}
	httpresponse.Success(c, newReadinessResponse())
}

func newReadinessResponse() ReadinessResponse {
	return ReadinessResponse{
		Status: "OK",
	}
}
