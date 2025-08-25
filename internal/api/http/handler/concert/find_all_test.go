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
	"github.com/kittipat1413/go-common/util/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	concertUsecase "ticket-reservation/internal/usecase/concert"
	"ticket-reservation/pkg/testhelper"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
)

func TestConcertHandler_FindAllConcerts(t *testing.T) {
	// Test data setup
	concertID1 := uuid.New()
	concertID2 := uuid.New()
	bangkokTime, _ := time.LoadLocation("Asia/Bangkok")
	concert1Date := time.Date(2025, 6, 15, 19, 0, 0, 0, bangkokTime)
	concert2Date := time.Date(2025, 8, 20, 20, 0, 0, 0, bangkokTime)
	createdTime := time.Date(2024, 12, 1, 10, 0, 0, 0, bangkokTime)
	updatedTime := time.Date(2024, 12, 15, 15, 30, 0, 0, bangkokTime)

	testConcerts := []entity.Concert{
		{
			ID:        concertID1,
			Name:      "Summer Festival 2025",
			Venue:     "Bangkok Arena",
			Date:      concert1Date,
			CreatedAt: createdTime,
			UpdatedAt: updatedTime,
		},
		{
			ID:        concertID2,
			Name:      "Rock Concert 2025",
			Venue:     "Impact Arena",
			Date:      concert2Date,
			CreatedAt: createdTime,
			UpdatedAt: updatedTime,
		},
	}

	// Create paginated result
	mockPagination := entity.NewPagination(2, 10, 0)
	provider := func() ([]entity.Concert, entity.PageProvider[entity.Concert], entity.Pagination, error) {
		return testConcerts, nil, mockPagination, nil
	}
	paginatedResult, _ := entity.NewPage(provider)

	tests := []struct {
		name             string
		queryParams      map[string]interface{}
		setupMocks       func(h *testHelper)
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:        "successful retrieval with default parameters",
			queryParams: map[string]interface{}{
				// No parameters - should use defaults
			},
			setupMocks: func(h *testHelper) {
				expectedInput := concertUsecase.FindAllConcertsInput{
					StartDate: nil,
					EndDate:   nil,
					Venue:     nil,
					Limit:     pointer.ToPointer(int64(100)),
					Offset:    pointer.ToPointer(int64(0)),
					SortBy:    pointer.ToPointer("date"),
					SortOrder: pointer.ToPointer(entity.SortOrderAsc),
				}
				h.mockConcertUsecase.EXPECT().
					FindAllConcerts(gomock.Any(), expectedInput).
					Return(paginatedResult, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"code": "ERR-200000",
				"data": []interface{}{
					map[string]interface{}{
						"id":    concertID1.String(),
						"name":  "Summer Festival 2025",
						"venue": "Bangkok Arena",
						"date":  "2025-06-15T19:00:00+07:00",
					},
					map[string]interface{}{
						"id":    concertID2.String(),
						"name":  "Rock Concert 2025",
						"venue": "Impact Arena",
						"date":  "2025-08-20T20:00:00+07:00",
					},
				},
				"metadata": map[string]interface{}{
					"pagination": map[string]interface{}{
						"current_page": float64(1),
						"limit":        float64(10),
						"offset":       float64(0),
						"page_count":   float64(1),
						"total":        float64(2),
					},
				},
			},
		},
		{
			name: "successful retrieval with filters",
			queryParams: map[string]interface{}{
				"startDate": "2025-01-01",
				"endDate":   "2025-12-31",
				"venue":     "Bangkok Arena",
				"limit":     5,
				"offset":    0,
				"sortBy":    "name",
				"sortOrder": "desc",
			},
			setupMocks: func(h *testHelper) {
				expectedInput := concertUsecase.FindAllConcertsInput{
					StartDate: pointer.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, bangkokTime)),
					EndDate:   pointer.ToPointer(time.Date(2025, 12, 31, 0, 0, 0, 0, bangkokTime)),
					Venue:     pointer.ToPointer("Bangkok Arena"),
					Limit:     pointer.ToPointer(int64(5)),
					Offset:    pointer.ToPointer(int64(0)),
					SortBy:    pointer.ToPointer("name"),
					SortOrder: pointer.ToPointer(entity.SortOrderDesc),
				}
				h.mockConcertUsecase.EXPECT().
					FindAllConcerts(gomock.Any(), gomock.Eq(expectedInput)).
					Return(paginatedResult, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"code": "ERR-200000",
				"data": []interface{}{
					map[string]interface{}{
						"id":    concertID1.String(),
						"name":  "Summer Festival 2025",
						"venue": "Bangkok Arena",
						"date":  "2025-06-15T19:00:00+07:00",
					},
					map[string]interface{}{
						"id":    concertID2.String(),
						"name":  "Rock Concert 2025",
						"venue": "Impact Arena",
						"date":  "2025-08-20T20:00:00+07:00",
					},
				},
				"metadata": map[string]interface{}{
					"pagination": map[string]interface{}{
						"current_page": float64(1),
						"limit":        float64(10),
						"offset":       float64(0),
						"page_count":   float64(1),
						"total":        float64(2),
					},
				},
			},
		},
		{
			name: "empty result",
			queryParams: map[string]interface{}{
				"venue": "Non-existent Venue",
			},
			setupMocks: func(h *testHelper) {
				emptyPagination := entity.NewPagination(0, 100, 0)
				emptyProvider := func() ([]entity.Concert, entity.PageProvider[entity.Concert], entity.Pagination, error) {
					return []entity.Concert{}, nil, emptyPagination, nil
				}
				emptyResult, _ := entity.NewPage(emptyProvider)
				h.mockConcertUsecase.EXPECT().
					FindAllConcerts(gomock.Any(), gomock.Any()).
					Return(emptyResult, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"code": "ERR-200000",
				"data": []interface{}{},
				"metadata": map[string]interface{}{
					"pagination": map[string]interface{}{
						"total":  float64(0),
						"limit":  float64(0),
						"offset": float64(0),
					},
				},
			},
		},
		{
			name: "invalid query parameters - malformed date",
			queryParams: map[string]interface{}{
				"startDate": "invalid-date",
			},
			setupMocks: func(h *testHelper) {
				// No usecase calls expected for query binding errors
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-401000",
				"message": "unable to parse request",
			},
		},
		{
			name: "usecase validation error",
			queryParams: map[string]interface{}{
				"limit": 1000, // Exceeds max limit
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					FindAllConcerts(gomock.Any(), gomock.Any()).
					Return(nil, errsFramework.NewBadRequestError("the request is invalid", nil))
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"code":    "ERR-401000",
				"message": "the request is invalid",
			},
		},
		{
			name: "usecase internal error",
			queryParams: map[string]interface{}{
				"limit": 10,
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertUsecase.EXPECT().
					FindAllConcerts(gomock.Any(), gomock.Any()).
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

			// Create Gin context with query parameters using testhelper
			c := testhelper.NewGinCtx(w).
				Method(http.MethodGet).
				Path("/concerts").
				Queries(tt.queryParams).
				WithContext(logger.NewContext(context.Background(), logger.NewNoopLogger())).
				MustBuild(t)

			// Execute the handler
			h.concertHandler.FindAllConcerts(c)

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
