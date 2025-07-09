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

func TestConcertRepositoryImpl_CreateOne(t *testing.T) {
	testDate := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	createdAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)
	testID := uuid.New()

	tests := []struct {
		name            string
		input           *entity.Concert
		setupMock       func(mock sqlmock.Sqlmock, input *entity.Concert)
		expectedConcert *entity.Concert
		expectedError   bool
		errorType       error
	}{
		{
			name: "successful concert creation",
			input: &entity.Concert{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  testDate,
			},
			setupMock: func(mock sqlmock.Sqlmock, input *entity.Concert) {
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.date",
					"concerts.venue", "concerts.created_at", "concerts.updated_at",
				}).AddRow(
					testID, input.Name, input.Date,
					input.Venue, createdAt, updatedAt,
				)

				mock.ExpectQuery(`INSERT INTO public\.concerts \(name, date, venue\) VALUES \(\$1, \$2, \$3\) RETURNING concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at"`).
					WithArgs(input.Name, input.Date, input.Venue).
					WillReturnRows(rows)
			},
			expectedConcert: &entity.Concert{
				ID:        testID,
				Name:      "Test Concert",
				Venue:     "Test Venue",
				Date:      testDate,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			expectedError: false,
		},
		{
			name: "successful creation with empty strings",
			input: &entity.Concert{
				Name:  "",
				Venue: "",
				Date:  testDate,
			},
			setupMock: func(mock sqlmock.Sqlmock, input *entity.Concert) {
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.date",
					"concerts.venue", "concerts.created_at", "concerts.updated_at",
				}).AddRow(
					testID, input.Name, input.Date,
					input.Venue, createdAt, updatedAt,
				)

				mock.ExpectQuery(`INSERT INTO public\.concerts \(name, date, venue\) VALUES \(\$1, \$2, \$3\) RETURNING concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at"`).
					WithArgs(input.Name, input.Date, input.Venue).
					WillReturnRows(rows)
			},
			expectedConcert: &entity.Concert{
				ID:        testID,
				Name:      "",
				Venue:     "",
				Date:      testDate,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			expectedError: false,
		},
		{
			name: "database constraint violation",
			input: &entity.Concert{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  testDate,
			},
			setupMock: func(mock sqlmock.Sqlmock, input *entity.Concert) {
				mock.ExpectQuery(`INSERT INTO public\.concerts \(name, date, venue\) VALUES \(\$1, \$2, \$3\) RETURNING concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at"`).
					WithArgs(input.Name, input.Date, input.Venue).
					WillReturnError(errors.New("pq: duplicate key value violates unique constraint"))
			},
			expectedConcert: nil,
			expectedError:   true,
			errorType:       &errsFramework.DatabaseError{},
		},
		{
			name: "database connection error",
			input: &entity.Concert{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  testDate,
			},
			setupMock: func(mock sqlmock.Sqlmock, input *entity.Concert) {
				mock.ExpectQuery(`INSERT INTO public\.concerts \(name, date, venue\) VALUES \(\$1, \$2, \$3\) RETURNING concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at"`).
					WithArgs(input.Name, input.Date, input.Venue).
					WillReturnError(sql.ErrConnDone)
			},
			expectedConcert: nil,
			expectedError:   true,
			errorType:       &errsFramework.DatabaseError{},
		},
		{
			name: "database timeout error",
			input: &entity.Concert{
				Name:  "Test Concert",
				Venue: "Test Venue",
				Date:  testDate,
			},
			setupMock: func(mock sqlmock.Sqlmock, input *entity.Concert) {
				mock.ExpectQuery(`INSERT INTO public\.concerts \(name, date, venue\) VALUES \(\$1, \$2, \$3\) RETURNING concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at"`).
					WithArgs(input.Name, input.Date, input.Venue).
					WillReturnError(context.DeadlineExceeded)
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

			tt.setupMock(h.Mock, tt.input)
			concert, err := h.Repository.CreateOne(context.Background(), tt.input)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository concert/create_one CreateOne]")

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

func TestConcertRepositoryImpl_CreateOne_QueryValidation(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testDate := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	createdAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)
	testID := uuid.New()

	input := &entity.Concert{
		Name:  "Test Concert",
		Venue: "Test Venue",
		Date:  testDate,
	}

	// Setup expectations - verify exact query structure
	rows := sqlmock.NewRows([]string{
		"concerts.id", "concerts.name", "concerts.date",
		"concerts.venue", "concerts.created_at", "concerts.updated_at",
	}).AddRow(
		testID, input.Name, input.Date,
		input.Venue, createdAt, updatedAt,
	)

	// The query should be an INSERT with RETURNING clause
	expectedQuery := `INSERT INTO public\.concerts \(name, date, venue\) VALUES \(\$1, \$2, \$3\) RETURNING concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at"`

	h.Mock.ExpectQuery(expectedQuery).
		WithArgs(input.Name, input.Date, input.Venue).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	concert, err := h.Repository.CreateOne(ctx, input)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, concert)
	assert.Equal(t, testID, concert.ID)
	assert.Equal(t, "Test Concert", concert.Name)
	assert.Equal(t, "Test Venue", concert.Venue)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
