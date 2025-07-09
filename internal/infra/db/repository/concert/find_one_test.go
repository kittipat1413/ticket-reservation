package concertrepo_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func TestConcertRepositoryImpl_FindOne(t *testing.T) {
	testID := uuid.New()
	testTime := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	createdAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name            string
		concertID       uuid.UUID
		setupMock       func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectedConcert *entity.Concert
		expectedError   bool
		errorType       error
	}{
		{
			name:      "successful concert retrieval",
			concertID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.venue",
					"concerts.date", "concerts.created_at", "concerts.updated_at",
				}).AddRow(
					id, "Test Concert", "Test Venue",
					testTime, createdAt, updatedAt,
				)

				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE concerts\.id = \$1`).
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedConcert: &entity.Concert{
				ID:        testID,
				Name:      "Test Concert",
				Venue:     "Test Venue",
				Date:      testTime,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			expectedError: false,
		},
		{
			name:      "concert not found",
			concertID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE concerts\.id = \$1`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedConcert: nil,
			expectedError:   true,
			errorType:       &errsFramework.NotFoundError{},
		},
		{
			name:      "database connection error",
			concertID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE concerts\.id = \$1`).
					WithArgs(id).
					WillReturnError(sql.ErrConnDone)
			},
			expectedConcert: nil,
			expectedError:   true,
			errorType:       &errsFramework.DatabaseError{},
		},
		{
			name:      "database timeout error",
			concertID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE concerts\.id = \$1`).
					WithArgs(id).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedConcert: nil,
			expectedError:   true,
			errorType:       &errsFramework.DatabaseError{},
		},
		{
			name:      "generic database error",
			concertID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE concerts\.id = \$1`).
					WithArgs(id).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedConcert: nil,
			expectedError:   true,
			errorType:       &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			tt.setupMock(h.Mock, tt.concertID)
			concert, err := h.Repository.FindOne(context.Background(), tt.concertID)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository concert/find_one FindOne]")

				// Verify it's the expected error type
				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, concert)
			} else {
				require.NoError(t, err)
				require.NotNil(t, concert)

				// Compare all fields
				assert.Equal(t, tt.expectedConcert.ID, concert.ID)
				assert.Equal(t, tt.expectedConcert.Name, concert.Name)
				assert.Equal(t, tt.expectedConcert.Venue, concert.Venue)
				assert.Equal(t, tt.expectedConcert.Date.UTC(), concert.Date.UTC())
				assert.Equal(t, tt.expectedConcert.CreatedAt.UTC(), concert.CreatedAt.UTC())
				assert.Equal(t, tt.expectedConcert.UpdatedAt.UTC(), concert.UpdatedAt.UTC())
			}

			// Verify all expectations were met
			h.AssertExpectationsMet(t)
		})
	}
}

func TestConcertRepositoryImpl_FindOne_QueryValidation(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testTime := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	createdAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)

	// Setup expectations - verify exact query structure
	rows := sqlmock.NewRows([]string{
		"concerts.id", "concerts.name", "concerts.venue",
		"concerts.date", "concerts.created_at", "concerts.updated_at",
	}).AddRow(
		testID, "Test Concert", "Test Venue",
		testTime, createdAt, updatedAt,
	)

	// The query should include all columns and proper WHERE clause
	expectedQuery := `SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE concerts\.id = \$1`

	h.Mock.ExpectQuery(expectedQuery).
		WithArgs(testID).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	concert, err := h.Repository.FindOne(ctx, testID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, concert)
	assert.Equal(t, testID, concert.ID)
	assert.Equal(t, "Test Concert", concert.Name)
	assert.Equal(t, "Test Venue", concert.Venue)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
