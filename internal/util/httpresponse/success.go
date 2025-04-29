package httpresponse

import (
	"net/http"
	"ticket-reservation/internal/domain/errs"

	"github.com/gin-gonic/gin"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

type SuccessResponse struct {
	Code     string `json:"code"`
	Data     any    `json:"data,omitempty"`
	Metadata any    `json:"metadata,omitempty"`
}

func Success(c *gin.Context, data any) {
	successResponse := SuccessResponse{
		Code: errsFramework.GetFullCode(errs.StatusCodeSuccess),
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, successResponse)
}

func SuccessWithCode(c *gin.Context, code string, data any) {
	successResponse := SuccessResponse{
		Code: errsFramework.GetFullCode(code),
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, successResponse)
}
