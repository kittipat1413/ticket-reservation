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

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func TestSeatRepositoryImpl_FindOne(t *testing.T) {
	testID := uuid.New()
	testZoneID := uuid.New()
	testSessionID := "session-123"
	testSeatNumber := "A1"
	testLockedUntil := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		seatID        uuid.UUID
		setupMock     func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectedSeat  *entity.Seat
		expectedError bool
		errorType     error
	}{
		{
			name:   "successful seat retrieval",
			seatID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
					"seats.locked_until", "seats.locked_by_session_id",
					"seats.created_at", "seats.updated_at",
				}).AddRow(
					id, testZoneID, testSeatNumber, entity.SeatStatusAvailable.String(),
					nil, nil, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedSeat: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusAvailable,
				LockedUntil:       nil,
				LockedBySessionID: nil,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name:   "successful seat retrieval with locked status",
			seatID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
					"seats.locked_until", "seats.locked_by_session_id",
					"seats.created_at", "seats.updated_at",
				}).AddRow(
					id, testZoneID, testSeatNumber, entity.SeatStatusPending.String(),
					testLockedUntil, testSessionID, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`).
					WithArgs(id).
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
			name:   "seat not found",
			seatID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.NotFoundError{},
		},
		{
			name:   "database connection error",
			seatID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnError(sql.ErrConnDone)
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name:   "database timeout error",
			seatID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name:   "generic database error",
			seatID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`).
					WithArgs(id).
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

			tt.setupMock(h.Mock, tt.seatID)
			seat, err := h.Repository.FindOne(context.Background(), tt.seatID)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository seat/find_one FindOne]")

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

func TestSeatRepositoryImpl_FindOne_QueryValidation(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testZoneID := uuid.New()
	testSeatNumber := "A1"
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	// Setup expectations - verify exact query structure
	rows := sqlmock.NewRows([]string{
		"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
		"seats.locked_until", "seats.locked_by_session_id",
		"seats.created_at", "seats.updated_at",
	}).AddRow(
		testID, testZoneID, testSeatNumber, entity.SeatStatusAvailable.String(),
		nil, nil, testCreatedAt, testUpdatedAt,
	)

	// The query should include all columns, FOR UPDATE clause, and proper WHERE clause
	expectedQuery := `SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`

	h.Mock.ExpectQuery(expectedQuery).
		WithArgs(testID).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	seat, err := h.Repository.FindOne(ctx, testID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, seat)
	assert.Equal(t, testID, seat.ID)
	assert.Equal(t, testZoneID, seat.ZoneID)
	assert.Equal(t, testSeatNumber, seat.SeatNumber)
	assert.Equal(t, entity.SeatStatusAvailable, seat.Status)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}

func TestSeatRepositoryImpl_FindOne_ModelConversionError(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testZoneID := uuid.New()
	testSeatNumber := "A1"
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	// Setup expectations with invalid status to trigger model conversion error
	rows := sqlmock.NewRows([]string{
		"seats.id", "seats.zone_id", "seats.seat_number", "seats.status",
		"seats.locked_until", "seats.locked_by_session_id",
		"seats.created_at", "seats.updated_at",
	}).AddRow(
		testID, testZoneID, testSeatNumber, "invalid_status",
		nil, nil, testCreatedAt, testUpdatedAt,
	)

	h.Mock.ExpectQuery(`SELECT seats\.id AS "seats\.id", seats\.zone_id AS "seats\.zone_id", seats\.seat_number AS "seats\.seat_number", seats\.status AS "seats\.status", seats\.locked_until AS "seats\.locked_until", seats\.locked_by_session_id AS "seats\.locked_by_session_id", seats\.created_at AS "seats\.created_at", seats\.updated_at AS "seats\.updated_at" FROM public\.seats WHERE seats\.id = \$1 FOR UPDATE`).
		WithArgs(testID).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	seat, err := h.Repository.FindOne(ctx, testID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "[repository seat/find_one FindOne]")
	assert.Contains(t, err.Error(), "failed to convert seat model to entity")

	var internalErr *errsFramework.InternalServerError
	assert.ErrorAs(t, err, &internalErr)

	assert.Nil(t, seat)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
