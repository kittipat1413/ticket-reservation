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
	"ticket-reservation/internal/domain/repository"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/util/pointer"
)

func TestReservationRepositoryImpl_UpdateOne(t *testing.T) {
	testID := uuid.New()
	testSeatID := uuid.New()
	testSessionID := "session-123"
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name                string
		input               repository.UpdateReservationInput
		setupMock           func(mock sqlmock.Sqlmock)
		expectedReservation *entity.Reservation
		expectedError       bool
		errorType           error
	}{
		{
			name: "successful status update",
			input: repository.UpdateReservationInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.ReservationStatusConfirmed),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"reservations.id", "reservations.seat_id", "reservations.session_id",
					"reservations.status", "reservations.reserved_at", "reservations.expires_at",
					"reservations.created_at", "reservations.updated_at",
				}).AddRow(
					testID, testSeatID, testSessionID, entity.ReservationStatusConfirmed.String(),
					testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`UPDATE public\.reservations SET status = \$1 WHERE reservations\.id = \$2 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(entity.ReservationStatusConfirmed.String(), testID).
					WillReturnRows(rows)
			},
			expectedReservation: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     entity.ReservationStatusConfirmed,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name: "successful expires_at update",
			input: repository.UpdateReservationInput{
				ID:        testID,
				ExpiresAt: pointer.ToPointer(testExpiresAt),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"reservations.id", "reservations.seat_id", "reservations.session_id",
					"reservations.status", "reservations.reserved_at", "reservations.expires_at",
					"reservations.created_at", "reservations.updated_at",
				}).AddRow(
					testID, testSeatID, testSessionID, entity.ReservationStatusPending.String(),
					testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`UPDATE public\.reservations SET expires_at = \$1 WHERE reservations\.id = \$2 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(testExpiresAt, testID).
					WillReturnRows(rows)
			},
			expectedReservation: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     entity.ReservationStatusPending,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name: "successful multiple fields update",
			input: repository.UpdateReservationInput{
				ID:        testID,
				Status:    pointer.ToPointer(entity.ReservationStatusConfirmed),
				ExpiresAt: pointer.ToPointer(testExpiresAt),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"reservations.id", "reservations.seat_id", "reservations.session_id",
					"reservations.status", "reservations.reserved_at", "reservations.expires_at",
					"reservations.created_at", "reservations.updated_at",
				}).AddRow(
					testID, testSeatID, testSessionID, entity.ReservationStatusConfirmed.String(),
					testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`UPDATE public\.reservations SET \(status, expires_at\) = \(\$1, \$2\) WHERE reservations\.id = \$3 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(entity.ReservationStatusConfirmed.String(), testExpiresAt, testID).
					WillReturnRows(rows)
			},
			expectedReservation: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     entity.ReservationStatusConfirmed,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name: "reservation not found",
			input: repository.UpdateReservationInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.ReservationStatusConfirmed),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.reservations SET status = \$1 WHERE reservations\.id = \$2 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(entity.ReservationStatusConfirmed.String(), testID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedReservation: nil,
			expectedError:       true,
			errorType:           &errsFramework.NotFoundError{},
		},
		{
			name: "database connection error",
			input: repository.UpdateReservationInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.ReservationStatusConfirmed),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.reservations SET status = \$1 WHERE reservations\.id = \$2 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(entity.ReservationStatusConfirmed.String(), testID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedReservation: nil,
			expectedError:       true,
			errorType:           &errsFramework.DatabaseError{},
		},
		{
			name: "database timeout error",
			input: repository.UpdateReservationInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.ReservationStatusConfirmed),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.reservations SET status = \$1 WHERE reservations\.id = \$2 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(entity.ReservationStatusConfirmed.String(), testID).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedReservation: nil,
			expectedError:       true,
			errorType:           &errsFramework.DatabaseError{},
		},
		{
			name: "generic database error",
			input: repository.UpdateReservationInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.ReservationStatusConfirmed),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.reservations SET status = \$1 WHERE reservations\.id = \$2 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
					WithArgs(entity.ReservationStatusConfirmed.String(), testID).
					WillReturnError(errors.New("database connection failed"))
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
			reservation, err := h.Repository.UpdateOne(context.Background(), tt.input)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository reservation/update_one UpdateOne]")

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

func TestReservationRepositoryImpl_UpdateOne_NoFieldsToUpdate(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()

	// Input with no fields to update
	input := repository.UpdateReservationInput{
		ID: testID,
		// No Status or ExpiresAt provided
	}

	// No mock expectations since no query should be executed
	ctx := context.Background()

	// Execute
	reservation, err := h.Repository.UpdateOne(ctx, input)

	// Assert - Should return an error since no fields to update
	require.Error(t, err)
	assert.Contains(t, err.Error(), "[repository reservation/update_one UpdateOne]")
	assert.Contains(t, err.Error(), "no fields provided to update")

	// The actual error type depends on the implementation
	// It might be an InternalServerError or a ValidationError
	assert.NotNil(t, err)
	assert.Nil(t, reservation)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}

func TestReservationRepositoryImpl_UpdateOne_ModelConversionError(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testSeatID := uuid.New()
	testSessionID := "session-123"
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	input := repository.UpdateReservationInput{
		ID:     testID,
		Status: pointer.ToPointer(entity.ReservationStatusConfirmed),
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

	h.Mock.ExpectQuery(`UPDATE public\.reservations SET status = \$1 WHERE reservations\.id = \$2 RETURNING reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at"`).
		WithArgs(entity.ReservationStatusConfirmed.String(), testID).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	reservation, err := h.Repository.UpdateOne(ctx, input)

	// Assert - Should return an internal server error when ToEntity() returns nil
	require.Error(t, err)
	assert.Contains(t, err.Error(), "[repository reservation/update_one UpdateOne]")
	assert.Contains(t, err.Error(), "failed to convert reservation model to entity")

	var internalErr *errsFramework.InternalServerError
	assert.ErrorAs(t, err, &internalErr)
	assert.Nil(t, reservation)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
