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

func TestConcertUsecase_FindOneConcert(t *testing.T) {
	testID := uuid.New()
	testTime := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)

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
		input          concertusecase.FindOneConcertInput
		setupMocks     func(h *testHelper)
		expectedResult *entity.Concert
		expectedError  bool
		errorType      error
		errorContains  string
	}{
		{
			name: "successful concert retrieval",
			input: concertusecase.FindOneConcertInput{
				ID: testID.String(),
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					FindOne(gomock.Any(), testID).
					Return(expectedConcert, nil)
			},
			expectedResult: expectedConcert,
			expectedError:  false,
		},
		{
			name: "validation error - empty ID",
			input: concertusecase.FindOneConcertInput{
				ID: "",
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "validation error - invalid UUID format",
			input: concertusecase.FindOneConcertInput{
				ID: "invalid-uuid",
			},
			setupMocks:     func(h *testHelper) {},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.BadRequestError{},
			errorContains:  "the request is invalid",
		},
		{
			name: "concert not found",
			input: concertusecase.FindOneConcertInput{
				ID: testID.String(),
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					FindOne(gomock.Any(), testID).
					Return(nil, errsFramework.NewNotFoundError("concert not found", nil))
			},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.NotFoundError{},
			errorContains:  "concert not found",
		},
		{
			name: "repository error - database failure",
			input: concertusecase.FindOneConcertInput{
				ID: testID.String(),
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					FindOne(gomock.Any(), testID).
					Return(nil, errsFramework.NewDatabaseError("connection failed", "error"))
			},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.InternalServerError{},
			errorContains:  "failed to find concert by ID",
		},
		{
			name: "repository error - internal server error",
			input: concertusecase.FindOneConcertInput{
				ID: testID.String(),
			},
			setupMocks: func(h *testHelper) {
				h.mockConcertRepository.EXPECT().
					FindOne(gomock.Any(), testID).
					Return(nil, errsFramework.NewInternalServerError("internal error", nil))
			},
			expectedResult: nil,
			expectedError:  true,
			errorType:      &errsFramework.InternalServerError{},
			errorContains:  "failed to find concert by ID",
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
			result, err := h.concertUsecase.FindOneConcert(ctx, tt.input)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[usecase concert/find_one_concert FindOneConcert]")

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
