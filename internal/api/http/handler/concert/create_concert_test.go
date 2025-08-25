package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	concertUsecase "ticket-reservation/internal/usecase/concert"
	"ticket-reservation/pkg/testhelper"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
)

func TestConcertHandler_CreateConcert(t *testing.T) {
	bangkokTime, _ := time.LoadLocation("Asia/Bangkok")
	expectedConcert := &entity.Concert{
		ID:        uuid.New(),
		Name:      "New Year Concert 2025",
		Venue:     "Bangkok Arena",
		Date:      time.Date(2025, 12, 25, 20, 0, 0, 0, bangkokTime),
		CreatedAt: time.Date(2024, 12, 1, 10, 0, 0, 0, bangkokTime),
		UpdatedAt: time.Date(2024, 12, 15, 15, 30, 0, 0, bangkokTime),
	}

	validRequestBody := map[string]interface{}{
		"name":  "New Year Concert 2025",
		"venue": "Bangkok Arena",
		"date":  "2025-12-25T20:00:00+07:00",
	}

	tests := []struct {
		name             string
		requestBody      interface{}
		setupMocks       func(h *testHelper)
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:        "successful concert creation",
			requestBody: validRequestBody,
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					CreateConcert(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, input concertUsecase.CreateConcertInput) (*entity.Concert, error) {
						// Validate input
						assert.Equal(t, "New Year Concert 2025", input.Name)
						assert.Equal(t, "Bangkok Arena", input.Venue)
						assert.Equal(t, time.Date(2025, 12, 25, 20, 0, 0, 0, bangkokTime).UTC(), input.Date.UTC())
						return expectedConcert, nil
					})
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: map[string]interface{}{
				"code": "ERR-200000",
				"data": map[string]interface{}{
					"id":    expectedConcert.ID.String(),
					"name":  "New Year Concert 2025",
					"venue": "Bangkok Arena",
					"date":  "2025-12-25T20:00:00+07:00",
				},
			},
		},
		{
			name: "invalid JSON body - missing required fields",
			requestBody: map[string]interface{}{
				"name": "Concert Without Venue",
				// Missing venue and date
			},
			setupMocks: func(h *testHelper) {
				// No usecase calls expected for validation errors
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-401000",
				"message": "unable to parse request",
			},
		},
		{
			name: "invalid JSON body - malformed date",
			requestBody: map[string]interface{}{
				"name":  "New Year Concert 2025",
				"venue": "Bangkok Arena",
				"date":  "invalid-date-format",
			},
			setupMocks: func(h *testHelper) {
				// No usecase calls expected for validation errors
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-401000",
				"message": "unable to parse request",
			},
		},
		{
			name:        "usecase validation error",
			requestBody: validRequestBody,
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					CreateConcert(gomock.Any(), gomock.Any()).
					Return(nil, errsFramework.NewBadRequestError("the request is invalid", nil))
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-401000",
				"message": "the request is invalid",
			},
		},
		{
			name:        "usecase internal error",
			requestBody: validRequestBody,
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					CreateConcert(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database connection failed"))
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

			// Create Gin context with JSON body using testhelper
			c := testhelper.NewGinCtx(w).
				Method(http.MethodPost).
				Path("/concerts").
				JSONBody(tt.requestBody).
				WithContext(logger.NewContext(context.Background(), logger.NewNoopLogger())).
				MustBuild(t)

			// Execute the handler
			h.concertHandler.CreateConcert(c)

			// Assert HTTP status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			for key, expectedValue := range tt.expectedResponse {
				actualValue, exists := responseBody[key]
				assert.True(t, exists, "Expected key '%s' to exist in response", key)
				assert.Equal(t, expectedValue, actualValue, "Mismatch for key '%s'", key)
			}
		})
	}
}
