package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kittipat1413/go-common/framework/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/pkg/testhelper"
)

func TestHealthCheckHandler_Readiness(t *testing.T) {
	tests := []struct {
		name             string
		setupMocks       func(h *testHelper)
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name: "successful readiness check",
			setupMocks: func(h *testHelper) {
				h.mockHealthcheckUsecase.EXPECT().
					CheckReadiness(gomock.Any()).
					Return(true, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"code": "ERR-200000",
				"data": map[string]interface{}{
					"status": "OK",
				},
			},
		},
		{
			name: "readiness check returns false",
			setupMocks: func(h *testHelper) {
				h.mockHealthcheckUsecase.EXPECT().
					CheckReadiness(gomock.Any()).
					Return(false, nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-500000",
				"message": "An unexpected error occurred. Please try again later.",
			},
		},
		{
			name: "readiness check returns error",
			setupMocks: func(h *testHelper) {
				h.mockHealthcheckUsecase.EXPECT().
					CheckReadiness(gomock.Any()).
					Return(false, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-500000",
				"message": "An unexpected error occurred. Please try again later.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			// Setup mocks for this test case
			tt.setupMocks(h)

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context using testhelper
			c := testhelper.NewGinCtx(w).
				Method(http.MethodGet).
				Path("/health/readiness").
				WithContext(logger.NewContext(context.Background(), logger.NewNoopLogger())).
				MustBuild(t)

			// Execute the handler
			h.healthcheckHandler.Readiness(c)

			// Assert HTTP status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			// Assert expected fields
			for key, expectedValue := range tt.expectedResponse {
				actualValue, exists := responseBody[key]
				assert.True(t, exists, "Expected key '%s' to exist in response", key)
				assert.Equal(t, expectedValue, actualValue, "Mismatch for key '%s'", key)
			}
		})
	}
}
