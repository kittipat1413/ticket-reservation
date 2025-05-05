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

// Success returns HTTP 200 with default internal success code.
func Success(c *gin.Context, data any) {
	successResponse := SuccessResponse{
		Code: errsFramework.GetFullCode(errs.StatusCodeSuccess),
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, successResponse)
}

// SuccessWithStatus returns custom HTTP status with default internal success code.
func SuccessWithStatus(c *gin.Context, httpStatus int, data any) {
	resp := SuccessResponse{
		Code: errsFramework.GetFullCode(errs.StatusCodeSuccess),
		Data: data,
	}
	c.AbortWithStatusJSON(httpStatus, resp)
}

// SuccessCustom returns custom HTTP status and custom internal success code.
func SuccessCustom(c *gin.Context, httpStatus int, code string, data any) {
	successResponse := SuccessResponse{
		Code: errsFramework.GetFullCode(code),
		Data: data,
	}
	c.AbortWithStatusJSON(httpStatus, successResponse)
}
