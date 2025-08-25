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
	concertusecase "ticket-reservation/internal/usecase/concert"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func TestConcertUsecase_CreateConcert(t *testing.T) {
	bangkokTime, _ := time.LoadLocation("Asia/Bangkok")
	testTime := time.Date(2025, 12, 25, 20, 0, 0, 0, bangkokTime)
	testID := uuid.New()

	expectedConcert := &entity.Concert{
		ID:        testID,
		Name:      "Test Concert",
		Venue:     "Test Venue",
		Date:      testTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		input          concertusecase.CreateConcertInput
		setupMocks     func(h *testHelper)
		expectedResult *entity.Concert
		expectedError  bool
		errorType      error
		errorContains  string
	}{
		{
			name: "successful concert creation",
			input: concertusecase.CreateConcertInput{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  testTime,
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					CreateOne(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, concert *entity.Concert) (*entity.Concert, error) {
						// Verify the input concert has the expected fields
						assert.Equal(h.ctrl.T, "Test Concert", concert.Name)
						assert.Equal(h.ctrl.T, "Test Venue", concert.Venue)
						assert.Equal(h.ctrl.T, testTime, concert.Date)
						return expectedConcert, nil
					})
			},
			expectedResult: expectedConcert,
			expectedError:  false,
		},
		{
			name: "validation error - empty name",
			input: concertusecase.CreateConcertInput{
				Name:  "",
				Venue: "Test Venue",
				Date:  testTime,
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - empty venue",
			input: concertusecase.CreateConcertInput{
				Name:  "Test Concert",
				Venue: "",
				Date:  testTime,
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - invalid timezone",
			input: concertusecase.CreateConcertInput{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC), // UTC instead of Bangkok timezone
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "repository error - database failure",
			input: concertusecase.CreateConcertInput{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  testTime,
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					CreateOne(gomock.Any(), gomock.Any()).
					Return(nil, errsFramework.NewDatabaseError("connection failed", "error"))
			},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.InternalServerError{},
			errorContains:  "failed to create concert",
		},
		{
			name: "repository error - conflict error",
			input: concertusecase.CreateConcertInput{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  testTime,
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					CreateOne(gomock.Any(), gomock.Any()).
					Return(nil, errsFramework.NewConflictError("concert already exists", nil))
			},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.InternalServerError{},
			errorContains:  "failed to create concert",
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
			result, err := h.concertUsecase.CreateConcert(ctx, tt.input)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[usecase concert/create_concert CreateConcert]")

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

				// Assert expected result if specified
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ID, result.ID)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.Venue, result.Venue)
					assert.Equal(t, tt.expectedResult.Date.UTC(), result.Date.UTC())
					assert.Equal(t, tt.expectedResult.CreatedAt.UTC(), result.CreatedAt.UTC())
					assert.Equal(t, tt.expectedResult.UpdatedAt.UTC(), result.UpdatedAt.UTC())
				}
			}
		})
	}
}
