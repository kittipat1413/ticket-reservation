package reservationrepo_test

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

func TestReservationRepositoryImpl_CreateOne(t *testing.T) {
	testID := uuid.New()
	testSeatID := uuid.New()
	testSessionID := "session-123"
	testStatus := entity.ReservationStatusPending
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	inputReservation := &entity.Reservation{
		SeatID:     testSeatID,
		SessionID:  testSessionID,
		Status:     testStatus,
		ReservedAt: testReservedAt,
		ExpiresAt:  testExpiresAt,
	}

	tests := []struct {
		name                string
		input               *entity.Reservation
		setupMock           func(mock sqlmock.Sqlmock)
		expectedReservation *entity.Reservation
		expectedError       bool
		errorType           error
	}{
		{
			name:  "successful reservation creation",
			input: inputReservation,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"reservations.id", "reservations.seat_id", "reservations.session_id",
					"reservations.status", "reservations.reserved_at", "reservations.expires_at",
					"reservations.created_at", "reservations.updated_at",
				}).AddRow(
					testID, testSeatID, testSessionID, testStatus.String(),
					testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`INSERT INTO public\.reservations \(seat_id, session_id, status, expires_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(testSeatID, testSessionID, testStatus.String(), testExpiresAt).
					WillReturnRows(rows)
			},
			expectedReservation: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     testStatus,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name:  "database connection error",
			input: inputReservation,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO public\.reservations \(seat_id, session_id, status, expires_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(testSeatID, testSessionID, testStatus.String(), testExpiresAt).
					WillReturnError(sql.ErrConnDone)
			},
			expectedReservation: nil,
			expectedError:       true,
			errorType:           &errsFramework.DatabaseError{},
		},
		{
			name:  "constraint violation error",
			input: inputReservation,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO public\.reservations \(seat_id, session_id, status, expires_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(testSeatID, testSessionID, testStatus.String(), testExpiresAt).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			expectedReservation: nil,
			expectedError:       true,
			errorType:           &errsFramework.DatabaseError{},
		},
		{
			name:  "database timeout error",
			input: inputReservation,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO public\.reservations \(seat_id, session_id, status, expires_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(testSeatID, testSessionID, testStatus.String(), testExpiresAt).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedReservation: nil,
			expectedError:       true,
			errorType:           &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			tt.setupMock(h.Mock)
			reservation, err := h.Repository.CreateOne(context.Background(), tt.input)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository reservation/create_one CreateOne]")

				// Verify it's the expected error type
				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, reservation)
			} else {
				require.NoError(t, err)
				require.NotNil(t, reservation)

				// Compare all fields
				assert.Equal(t, tt.expectedReservation.ID, reservation.ID)
				assert.Equal(t, tt.expectedReservation.SeatID, reservation.SeatID)
				assert.Equal(t, tt.expectedReservation.SessionID, reservation.SessionID)
				assert.Equal(t, tt.expectedReservation.Status, reservation.Status)
				assert.Equal(t, tt.expectedReservation.ReservedAt.UTC(), reservation.ReservedAt.UTC())
				assert.Equal(t, tt.expectedReservation.ExpiresAt.UTC(), reservation.ExpiresAt.UTC())
				assert.Equal(t, tt.expectedReservation.CreatedAt.UTC(), reservation.CreatedAt.UTC())
				assert.Equal(t, tt.expectedReservation.UpdatedAt.UTC(), reservation.UpdatedAt.UTC())
			}

			// Verify all expectations were met
			h.AssertExpectationsMet(t)
		})
	}
}

func TestReservationRepositoryImpl_CreateOne_QueryValidation(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testSeatID := uuid.New()
	testSessionID := "session-123"
	testStatus := entity.ReservationStatusPending
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	inputReservation := &entity.Reservation{
		SeatID:     testSeatID,
		SessionID:  testSessionID,
		Status:     testStatus,
		ReservedAt: testReservedAt,
		ExpiresAt:  testExpiresAt,
	}

	// Setup expectations - verify exact query structure
	rows := sqlmock.NewRows([]string{
		"reservations.id", "reservations.seat_id", "reservations.session_id",
		"reservations.status", "reservations.reserved_at", "reservations.expires_at",
		"reservations.created_at", "reservations.updated_at",
	}).AddRow(
		testID, testSeatID, testSessionID, testStatus.String(),
		testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
	)

	// The query should exclude default columns and return all columns
	expectedQuery := `INSERT INTO public\.reservations \(seat_id, session_id, status, expires_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`

	h.Mock.ExpectQuery(expectedQuery).
		WithArgs(testSeatID, testSessionID, testStatus.String(), testExpiresAt).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	reservation, err := h.Repository.CreateOne(ctx, inputReservation)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, reservation)
	assert.Equal(t, testID, reservation.ID)
	assert.Equal(t, testSeatID, reservation.SeatID)
	assert.Equal(t, testSessionID, reservation.SessionID)
	assert.Equal(t, testStatus, reservation.Status)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}

func TestReservationRepositoryImpl_CreateOne_ModelConversionError(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testSeatID := uuid.New()
	testSessionID := "session-123"
	testStatus := entity.ReservationStatusPending
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	inputReservation := &entity.Reservation{
		SeatID:     testSeatID,
		SessionID:  testSessionID,
		Status:     testStatus,
		ReservedAt: testReservedAt,
		ExpiresAt:  testExpiresAt,
	}

	// Setup mock to return invalid status that would cause ToEntity() to return nil
	rows := sqlmock.NewRows([]string{
		"reservations.id", "reservations.seat_id", "reservations.session_id",
		"reservations.status", "reservations.reserved_at", "reservations.expires_at",
		"reservations.created_at", "reservations.updated_at",
	}).AddRow(
		testID, testSeatID, testSessionID, "invalid_status", // Invalid status
		testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
	)

	h.Mock.ExpectQuery(`INSERT INTO public\.reservations \(seat_id, session_id, status, expires_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
		WithArgs(testSeatID, testSessionID, testStatus.String(), testExpiresAt).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	reservation, err := h.Repository.CreateOne(ctx, inputReservation)

	// Assert - Should return an internal server error when ToEntity() returns nil
	require.Error(t, err)
	assert.Contains(t, err.Error(), "[repository reservation/create_one CreateOne]")
	assert.Contains(t, err.Error(), "failed to convert reservation model to entity")

	var internalErr *errsFramework.InternalServerError
	assert.ErrorAs(t, err, &internalErr)
	assert.Nil(t, reservation)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
