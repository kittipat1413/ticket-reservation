package seatrepo_test

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

func TestSeatRepositoryImpl_UpdateOne(t *testing.T) {
	testID := uuid.New()
	testZoneID := uuid.New()
	testSeatNumber := "A1"
	testSessionID := "session-123"
	testLockedUntil := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		input         repository.UpdateSeatInput
		setupMock     func(mock sqlmock.Sqlmock)
		expectedSeat  *entity.Seat
		expectedError bool
		errorType     error
	}{
		{
			name: "successful status update",
			input: repository.UpdateSeatInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.SeatStatusPending),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
					"seats.locked_until", "seats.locked_by_session_id",
					"seats.created_at", "seats.updated_at",
				}).AddRow(
					testID, testZoneID, testSeatNumber, entity.SeatStatusPending.String(),
					nil, nil, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`UPDATE public\.seats SET status = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(entity.SeatStatusPending.String(), testID).
					WillReturnRows(rows)
			},
			expectedSeat: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusPending,
				LockedUntil:       nil,
				LockedBySessionID: nil,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name: "successful locked_until update",
			input: repository.UpdateSeatInput{
				ID:          testID,
				LockedUntil: &testLockedUntil,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
					"seats.locked_until", "seats.locked_by_session_id",
					"seats.created_at", "seats.updated_at",
				}).AddRow(
					testID, testZoneID, testSeatNumber, entity.SeatStatusPending.String(),
					testLockedUntil, nil, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`UPDATE public\.seats SET locked_until = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(testLockedUntil, testID).
					WillReturnRows(rows)
			},
			expectedSeat: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusPending,
				LockedUntil:       &testLockedUntil,
				LockedBySessionID: nil,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name: "successful locked_by_session_id update",
			input: repository.UpdateSeatInput{
				ID:                testID,
				LockedBySessionID: &testSessionID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
					"seats.locked_until", "seats.locked_by_session_id",
					"seats.created_at", "seats.updated_at",
				}).AddRow(
					testID, testZoneID, testSeatNumber, entity.SeatStatusPending.String(),
					nil, testSessionID, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`UPDATE public\.seats SET locked_by_session_id = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(testSessionID, testID).
					WillReturnRows(rows)
			},
			expectedSeat: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusPending,
				LockedUntil:       nil,
				LockedBySessionID: &testSessionID,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name: "successful multiple fields update",
			input: repository.UpdateSeatInput{
				ID:                testID,
				Status:            pointer.ToPointer(entity.SeatStatusPending),
				LockedUntil:       &testLockedUntil,
				LockedBySessionID: &testSessionID,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
					"seats.locked_until", "seats.locked_by_session_id",
					"seats.created_at", "seats.updated_at",
				}).AddRow(
					testID, testZoneID, testSeatNumber, entity.SeatStatusPending.String(),
					testLockedUntil, testSessionID, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`UPDATE public\.seats SET \(status, locked_until, locked_by_session_id\) = \(\$1, \$2, \$3\) WHERE seats\.id = \$4 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(entity.SeatStatusPending.String(), testLockedUntil, testSessionID, testID).
					WillReturnRows(rows)
			},
			expectedSeat: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusPending,
				LockedUntil:       &testLockedUntil,
				LockedBySessionID: &testSessionID,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name: "seat not found",
			input: repository.UpdateSeatInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.SeatStatusPending),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.seats SET status = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(entity.SeatStatusPending.String(), testID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.NotFoundError{},
		},
		{
			name: "database connection error",
			input: repository.UpdateSeatInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.SeatStatusPending),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.seats SET status = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(entity.SeatStatusPending.String(), testID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "database timeout error",
			input: repository.UpdateSeatInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.SeatStatusPending),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.seats SET status = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(entity.SeatStatusPending.String(), testID).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "generic database error",
			input: repository.UpdateSeatInput{
				ID:     testID,
				Status: pointer.ToPointer(entity.SeatStatusPending),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE public\.seats SET status = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
					WithArgs(entity.SeatStatusPending.String(), testID).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			tt.setupMock(h.Mock)
			seat, err := h.Repository.UpdateOne(context.Background(), tt.input)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository seat/update_one UpdateOne]")

				// Verify it's the expected error type
				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, seat)
			} else {
				require.NoError(t, err)
				require.NotNil(t, seat)

				// Compare all fields
				assert.Equal(t, tt.expectedSeat.ID, seat.ID)
				assert.Equal(t, tt.expectedSeat.ZoneID, seat.ZoneID)
				assert.Equal(t, tt.expectedSeat.SeatNumber, seat.SeatNumber)
				assert.Equal(t, tt.expectedSeat.Status, seat.Status)

				if tt.expectedSeat.LockedUntil != nil {
					require.NotNil(t, seat.LockedUntil)
					assert.Equal(t, tt.expectedSeat.LockedUntil.UTC(), seat.LockedUntil.UTC())
				} else {
					assert.Nil(t, seat.LockedUntil)
				}

				if tt.expectedSeat.LockedBySessionID != nil {
					require.NotNil(t, seat.LockedBySessionID)
					assert.Equal(t, *tt.expectedSeat.LockedBySessionID, *seat.LockedBySessionID)
				} else {
					assert.Nil(t, seat.LockedBySessionID)
				}

				assert.Equal(t, tt.expectedSeat.CreatedAt.UTC(), seat.CreatedAt.UTC())
				assert.Equal(t, tt.expectedSeat.UpdatedAt.UTC(), seat.UpdatedAt.UTC())
			}

			// Verify all expectations were met
			h.AssertExpectationsMet(t)
		})
	}
}

func TestSeatRepositoryImpl_UpdateOne_NoFieldsToUpdate(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()

	// No fields provided for update
	input := repository.UpdateSeatInput{
		ID: testID,
	}

	// Execute
	seat, err := h.Repository.UpdateOne(context.Background(), input)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "[repository seat/update_one UpdateOne]")
	assert.Contains(t, err.Error(), "no fields provided to update")

	var badRequestErr *errsFramework.BadRequestError
	assert.ErrorAs(t, err, &badRequestErr)

	assert.Nil(t, seat)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}

func TestSeatRepositoryImpl_UpdateOne_ModelConversionError(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testZoneID := uuid.New()
	testSeatNumber := "A1"
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	input := repository.UpdateSeatInput{
		ID:     testID,
		Status: pointer.ToPointer(entity.SeatStatusPending),
	}

	// Setup expectations with invalid status to trigger model conversion error
	rows := sqlmock.NewRows([]string{
		"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
		"seats.locked_until", "seats.locked_by_session_id",
		"seats.created_at", "seats.updated_at",
	}).AddRow(
		testID, testZoneID, testSeatNumber, "invalid_status",
		nil, nil, testCreatedAt, testUpdatedAt,
	)

	h.Mock.ExpectQuery(`UPDATE public\.seats SET status = \$1 WHERE seats\.id = \$2 RETURNING seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at"`).
		WithArgs(entity.SeatStatusPending.String(), testID).
		WillReturnRows(rows)

	// Execute
	seat, err := h.Repository.UpdateOne(context.Background(), input)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "[repository seat/update_one UpdateOne]")
	assert.Contains(t, err.Error(), "failed to convert seat model to entity")

	var internalErr *errsFramework.InternalServerError
	assert.ErrorAs(t, err, &internalErr)

	assert.Nil(t, seat)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
