package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/domain/repository"
	concertusecase "ticket-reservation/internal/usecase/concert"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/util/pointer"
)

func TestConcertUsecase_FindAllConcerts(t *testing.T) {
	bangkokTime, _ := time.LoadLocation("Asia/Bangkok")
	testStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, bangkokTime)
	testEndDate := time.Date(2025, 12, 31, 23, 59, 59, 0, bangkokTime)

	createdTime := time.Date(2024, 12, 1, 10, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2024, 12, 15, 15, 30, 0, 0, time.UTC)

	testConcerts := entity.Concerts{
		{
			ID:        uuid.New(),
			Name:      "Rock Concert 2025",
			Venue:     "Stadium A",
			Date:      time.Date(2025, 6, 15, 20, 0, 0, 0, time.UTC),
			CreatedAt: createdTime,
			UpdatedAt: updatedTime,
		},
		{
			ID:        uuid.New(),
			Name:      "Jazz Night",
			Venue:     "Theatre B",
			Date:      time.Date(2025, 8, 20, 19, 30, 0, 0, time.UTC),
			CreatedAt: createdTime,
			UpdatedAt: updatedTime,
		},
	}

	tests := []struct {
		name           string
		input          concertusecase.FindAllConcertsInput
		setupMocks     func(h *testHelper)
		expectedResult *entity.Concerts
		expectedError  bool
		errorType      error
		errorContains  string
		expectedCount  int
	}{
		{
			name: "successful find all concerts with filters and verify content",
			input: concertusecase.FindAllConcertsInput{
				StartDate: &testStartDate,
				EndDate:   &testEndDate,
				Venue:     pointer.ToPointer("Stadium"),
				Limit:     pointer.ToPointer(int64(10)),
				Offset:    pointer.ToPointer(int64(0)),
				SortBy:    pointer.ToPointer("date"),
				SortOrder: pointer.ToPointer(entity.SortOrderAsc),
			},
			setupMocks: func(h *testHelper) {
				expectedFilter := repository.FindAllConcertsFilter{
					StartDate: &testStartDate,
					EndDate:   &testEndDate,
					Venue:     pointer.ToPointer("Stadium"),
					Limit:     pointer.ToPointer(int64(10)),
					Offset:    pointer.ToPointer(int64(0)),
					SortBy:    pointer.ToPointer("date"),
					SortOrder: pointer.ToPointer(entity.SortOrderAsc),
				}
				h.mockConcertRepository.EXPECT().
					FindAll(gomock.Any(), gomock.Eq(expectedFilter)).
					Return(&testConcerts, int64(2), nil)
			},
			expectedResult: &testConcerts,
			expectedError:  false,
			expectedCount:  2,
		},
		{
			name: "successful find all concerts with minimal input",
			input: concertusecase.FindAllConcertsInput{
				Limit:  pointer.ToPointer(int64(50)),
				Offset: pointer.ToPointer(int64(0)),
			},
			setupMocks: func(h *testHelper) {
				expectedFilter := repository.FindAllConcertsFilter{
					StartDate: nil,
					EndDate:   nil,
					Venue:     nil,
					Limit:     pointer.ToPointer(int64(50)),
					Offset:    pointer.ToPointer(int64(0)),
					SortBy:    nil,
					SortOrder: nil,
				}
				h.mockConcertRepository.EXPECT().
					FindAll(gomock.Any(), gomock.Eq(expectedFilter)).
					Return(&testConcerts, int64(2), nil)
			},
			expectedResult: &testConcerts,
			expectedError:  false,
			expectedCount:  2,
		},
		{
			name: "successful find all concerts with nil results",
			input: concertusecase.FindAllConcertsInput{
				Limit:  pointer.ToPointer(int64(10)),
				Offset: pointer.ToPointer(int64(0)),
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					FindAll(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), nil)
			},
			expectedResult: nil,
			expectedError:  false,
			expectedCount:  0,
		},
		{
			name: "validation error - missing limit",
			input: concertusecase.FindAllConcertsInput{
				Offset: pointer.ToPointer(int64(0)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - missing offset",
			input: concertusecase.FindAllConcertsInput{
				Limit: pointer.ToPointer(int64(10)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - limit too low",
			input: concertusecase.FindAllConcertsInput{
				Limit:  pointer.ToPointer(int64(0)),
				Offset: pointer.ToPointer(int64(0)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - limit too high",
			input: concertusecase.FindAllConcertsInput{
				Limit:  pointer.ToPointer(int64(101)),
				Offset: pointer.ToPointer(int64(0)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - negative offset",
			input: concertusecase.FindAllConcertsInput{
				Limit:  pointer.ToPointer(int64(10)),
				Offset: pointer.ToPointer(int64(-1)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - empty venue",
			input: concertusecase.FindAllConcertsInput{
				Venue:  pointer.ToPointer(""),
				Limit:  pointer.ToPointer(int64(10)),
				Offset: pointer.ToPointer(int64(0)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - invalid sort by",
			input: concertusecase.FindAllConcertsInput{
				Limit:     pointer.ToPointer(int64(10)),
				Offset:    pointer.ToPointer(int64(0)),
				SortBy:    pointer.ToPointer("invalid_field"),
				SortOrder: pointer.ToPointer(entity.SortOrderAsc),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - invalid sort order",
			input: concertusecase.FindAllConcertsInput{
				Limit:     pointer.ToPointer(int64(10)),
				Offset:    pointer.ToPointer(int64(0)),
				SortBy:    pointer.ToPointer("date"),
				SortOrder: (*entity.SortOrder)(pointer.ToPointer("invalid")),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - end date before start date",
			input: concertusecase.FindAllConcertsInput{
				StartDate: &testEndDate,
				EndDate:   &testStartDate, // End date before start date
				Limit:     pointer.ToPointer(int64(10)),
				Offset:    pointer.ToPointer(int64(0)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - invalid timezone for start date",
			input: concertusecase.FindAllConcertsInput{
				StartDate: pointer.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)), // UTC instead of Bangkok
				Limit:     pointer.ToPointer(int64(10)),
				Offset:    pointer.ToPointer(int64(0)),
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "repository error - database failure",
			input: concertusecase.FindAllConcertsInput{
				Limit:  pointer.ToPointer(int64(10)),
				Offset: pointer.ToPointer(int64(0)),
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					FindAll(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), errsFramework.NewDatabaseError("connection failed", "error"))
			},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.InternalServerError{},
			errorContains:  "failed to fetch concerts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			// Setup mocks
			tt.setupMocks(h)

			// Execute
			ctx := context.Background()
			result, err := h.concertUsecase.FindAllConcerts(ctx, tt.input)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[usecase concert/find_all_concerts FindAllConcerts]")

				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}

				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)

				// Assert data
				data := result.GetData()
				assert.Equal(t, tt.expectedCount, len(data))

				// Validate expected result content if specified
				if tt.expectedResult != nil && tt.expectedCount > 0 {
					for i, concert := range data {
						assert.Equal(t, pointer.GetValue(tt.expectedResult)[i], concert)
					}
				}

				// Validate individual concert data if present
				if tt.expectedCount > 0 {
					// Assert pagination details
					pagination := result.GetPagination()
					assert.NotNil(t, pagination)
					assert.Equal(t, int64(tt.expectedCount), pagination.Total)
					assert.Equal(t, *tt.input.Limit, pagination.Limit)
					assert.Equal(t, *tt.input.Offset, pagination.Offset)

					// Calculate expected pagination values
					expectedPageCount := pagination.Total / pagination.Limit
					if pagination.Total%pagination.Limit != 0 {
						expectedPageCount++
					}
					expectedCurrentPage := pagination.Offset/pagination.Limit + 1

					assert.Equal(t, expectedPageCount, pagination.PageCount)
					assert.Equal(t, expectedCurrentPage, pagination.CurrentPage)
				} else {
					// For empty results, pagination should be empty
					pagination := result.GetPagination()
					assert.Equal(t, entity.Pagination{}, pagination)
				}
			}
		})
	}
}
