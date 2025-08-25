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

func TestConcertHandler_FindConcertByID(t *testing.T) {
	// Test data setup
	concertID := uuid.New()
	bangkokTime, _ := time.LoadLocation("Asia/Bangkok")
	expectedConcert := &entity.Concert{
		ID:        concertID,
		Name:      "New Year Concert 2025",
		Venue:     "Bangkok Arena",
		Date:      time.Date(2025, 12, 25, 20, 0, 0, 0, bangkokTime),
		CreatedAt: time.Date(2024, 12, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 12, 15, 15, 30, 0, 0, time.UTC),
	}

	tests := []struct {
		name             string
		concertID        string
		setupMocks       func(h *testHelper)
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:      "successful concert retrieval",
			concertID: concertID.String(),
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					FindOneConcert(gomock.Any(), concertUsecase.FindOneConcertInput{
						ID: concertID.String(),
					}).
					Return(expectedConcert, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"code": "ERR-200000",
				"data": map[string]interface{}{
					"id":    concertID.String(),
					"name":  "New Year Concert 2025",
					"venue": "Bangkok Arena",
					"date":  "2025-12-25T20:00:00+07:00",
				},
			},
		},
		{
			name:      "concert not found",
			concertID: "550e8400-e29b-41d4-a716-446655440999",
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					FindOneConcert(gomock.Any(), concertUsecase.FindOneConcertInput{
						ID: "550e8400-e29b-41d4-a716-446655440999",
					}).
					Return(nil, errsFramework.NewNotFoundError("concert not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-402000",
				"message": "concert not found",
			},
		},
		{
			name:      "invalid UUID format",
			concertID: "invalid-uuid",
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					FindOneConcert(gomock.Any(), concertUsecase.FindOneConcertInput{
						ID: "invalid-uuid",
					}).
					Return(nil, errsFramework.NewBadRequestError("the request is invalid", nil))
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-401000",
				"message": "the request is invalid",
			},
		},
		{
			name:      "usecase internal error",
			concertID: concertID.String(),
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					FindOneConcert(gomock.Any(), concertUsecase.FindOneConcertInput{
						ID: concertID.String(),
					}).
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

			// Create Gin context with path parameter using testhelper
			c := testhelper.NewGinCtx(w).
				Method(http.MethodGet).
				Path("/concerts/:id").
				Param("id", tt.concertID).
				WithContext(logger.NewContext(context.Background(), logger.NewNoopLogger())).
				MustBuild(t)

			// Execute the handler
			h.concertHandler.FindConcertByID(c)

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
