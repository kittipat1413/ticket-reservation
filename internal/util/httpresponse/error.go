package httpresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
)

type ErrorResponse struct {
	Code    string `json:"code" example:"TR-XXXXXX"`
	Message string `json:"message" example:"Error message"`
	Data    any    `json:"data,omitempty"`
}

// Error sends an error response with the appropriate status code.
func Error(c *gin.Context, err error) {
	// Unwrap the error and extract the error response.
	httpCode, errResp := unwrapError(err)
	// Log the error with the appropriate context.
	appLogger := logger.FromContext(c.Request.Context())
	appLogger.Error(c.Request.Context(), errResp.Message, err, nil)
	// Send the error response.
	c.AbortWithStatusJSON(httpCode, errResp)
}

// unwrapError processes the error and extracts information for the response.
func unwrapError(err error) (httpCode int, errResp ErrorResponse) {
	// Default error response for non-domain errors.
	// This will be used if the error is not a errsFramework.DomainError.
	httpCode = http.StatusInternalServerError
	errResp = ErrorResponse{
		Code:    errsFramework.GetFullCode(errsFramework.StatusCodeGenericInternalServerError),
		Message: "An unexpected error occurred. Please try again later.",
	}

	// Try to unwrap the error and find a first valid errsFramework.DomainError in the chain.
	if domainErr := errsFramework.UnwrapDomainError(err); domainErr != nil {
		httpCode = domainErr.GetHTTPCode()
		errResp.Code = domainErr.Code()
		errResp.Message = domainErr.GetMessage()
		errResp.Data = domainErr.GetData()
	}

	return httpCode, errResp
}
