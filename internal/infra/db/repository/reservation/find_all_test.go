package reservationrepo_test

import (
	"context"
	"database/sql"
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

func TestReservationRepositoryImpl_FindAll(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()
	testSeatID1 := uuid.New()
	testSeatID2 := uuid.New()
	testSessionID := "session-123"
	testStatus := entity.ReservationStatusPending
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name                 string
		filter               repository.FindAllReservationsFilter
		setupMock            func(mock sqlmock.Sqlmock)
		expectedReservations *entity.Reservations
		expectedTotal        int64
		expectedError        bool
		errorType            error
	}{
		{
			name:   "successful retrieval with no filter",
			filter: repository.FindAllReservationsFilter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Count query
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(2)
				mock.ExpectQuery(`SELECT COUNT\(reservations\.id\) AS "total" FROM public\.reservations`).
					WillReturnRows(countRows)

				// Data query
				dataRows := sqlmock.NewRows([]string{
					"reservations.id", "reservations.seat_id", "reservations.session_id",
					"reservations.status", "reservations.reserved_at", "reservations.expires_at",
					"reservations.created_at", "reservations.updated_at",
				}).AddRow(
					testID1, testSeatID1, testSessionID, testStatus.String(),
					testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
				).AddRow(
					testID2, testSeatID2, testSessionID, testStatus.String(),
					testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`SELECT reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at" FROM public\.reservations`).
					WillReturnRows(dataRows)
			},
			expectedReservations: &entity.Reservations{
				{
					ID:         testID1,
					SeatID:     testSeatID1,
					SessionID:  testSessionID,
					Status:     testStatus,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
				{
					ID:         testID2,
					SeatID:     testSeatID2,
					SessionID:  testSessionID,
					Status:     testStatus,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedTotal: 2,
			expectedError: false,
		},
		{
			name: "successful retrieval with filters",
			filter: repository.FindAllReservationsFilter{
				SeatID:    &testSeatID1,
				SessionID: &testSessionID,
				Status:    &testStatus,
				Limit:     pointer.ToPointer(int64(10)),
				Offset:    pointer.ToPointer(int64(0)),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Count query
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(1)
				mock.ExpectQuery(`SELECT COUNT\(reservations\.id\) AS "total" FROM public\.reservations WHERE \( \(reservations\.seat_id = \$1\) AND \(reservations\.session_id = \$2::text\) AND \(reservations\.status = \$3::text\) \)`).
					WithArgs(testSeatID1, testSessionID, testStatus.String()).
					WillReturnRows(countRows)

				// Data query
				dataRows := sqlmock.NewRows([]string{
					"reservations.id", "reservations.seat_id", "reservations.session_id",
					"reservations.status", "reservations.reserved_at", "reservations.expires_at",
					"reservations.created_at", "reservations.updated_at",
				}).AddRow(
					testID1, testSeatID1, testSessionID, testStatus.String(),
					testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`SELECT reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at" FROM public\.reservations WHERE \( \(reservations\.seat_id = \$1\) AND \(reservations\.session_id = \$2::text\) AND \(reservations\.status = \$3::text\) \) LIMIT \$4 OFFSET \$5`).
					WithArgs(testSeatID1, testSessionID, testStatus.String(), int64(10), int64(0)).
					WillReturnRows(dataRows)
			},
			expectedReservations: &entity.Reservations{
				{
					ID:         testID1,
					SeatID:     testSeatID1,
					SessionID:  testSessionID,
					Status:     testStatus,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedTotal: 1,
			expectedError: false,
		},
		{
			name:   "database error in count query",
			filter: repository.FindAllReservationsFilter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(reservations\.id\) AS "total" FROM public\.reservations`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedReservations: nil,
			expectedTotal:        0,
			expectedError:        true,
			errorType:            &errsFramework.DatabaseError{},
		},
		{
			name:   "database error in data query",
			filter: repository.FindAllReservationsFilter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Count query succeeds
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(1)
				mock.ExpectQuery(`SELECT COUNT\(reservations\.id\) AS "total" FROM public\.reservations`).
					WillReturnRows(countRows)

				// Data query fails
				mock.ExpectQuery(`SELECT reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at" FROM public\.reservations`).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedReservations: nil,
			expectedTotal:        0,
			expectedError:        true,
			errorType:            &errsFramework.DatabaseError{},
		},
		{
			name:   "empty result set",
			filter: repository.FindAllReservationsFilter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Count query
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(reservations\.id\) AS "total" FROM public\.reservations`).
					WillReturnRows(countRows)

				// Data query
				dataRows := sqlmock.NewRows([]string{
					"reservations.id", "reservations.seat_id", "reservations.session_id",
					"reservations.status", "reservations.reserved_at", "reservations.expires_at",
					"reservations.created_at", "reservations.updated_at",
				})

				mock.ExpectQuery(`SELECT reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at" FROM public\.reservations`).
					WillReturnRows(dataRows)
			},
			expectedReservations: &entity.Reservations{},
			expectedTotal:        0,
			expectedError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			tt.setupMock(h.Mock)
			reservations, total, err := h.Repository.FindAll(context.Background(), tt.filter)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository reservation/find_all FindAll]")

				// Verify it's the expected error type
				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, reservations)
				assert.Equal(t, int64(0), total)
			} else {
				require.NoError(t, err)
				require.NotNil(t, reservations)
				assert.Equal(t, tt.expectedTotal, total)
				assert.Equal(t, len(*tt.expectedReservations), len(*reservations))

				// Compare reservations
				for i, expected := range *tt.expectedReservations {
					actual := (*reservations)[i]
					assert.Equal(t, expected.ID, actual.ID)
					assert.Equal(t, expected.SeatID, actual.SeatID)
					assert.Equal(t, expected.SessionID, actual.SessionID)
					assert.Equal(t, expected.Status, actual.Status)
					assert.Equal(t, expected.ReservedAt.UTC(), actual.ReservedAt.UTC())
					assert.Equal(t, expected.ExpiresAt.UTC(), actual.ExpiresAt.UTC())
					assert.Equal(t, expected.CreatedAt.UTC(), actual.CreatedAt.UTC())
					assert.Equal(t, expected.UpdatedAt.UTC(), actual.UpdatedAt.UTC())
				}
			}

			// Verify all expectations were met
			h.AssertExpectationsMet(t)
		})
	}
}

func TestReservationRepositoryImpl_FindAll_ModelConversionError(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()
	testSeatID := uuid.New()
	testSessionID := "session-123"
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	// Count query
	countRows := sqlmock.NewRows([]string{"total"}).AddRow(1)
	h.Mock.ExpectQuery(`SELECT COUNT\(reservations\.id\) AS "total" FROM public\.reservations`).
		WillReturnRows(countRows)

	// Data query with invalid status
	dataRows := sqlmock.NewRows([]string{
		"reservations.id", "reservations.seat_id", "reservations.session_id",
		"reservations.status", "reservations.reserved_at", "reservations.expires_at",
		"reservations.created_at", "reservations.updated_at",
	}).AddRow(
		testID, testSeatID, testSessionID, "invalid_status", // Invalid status
		testReservedAt, testExpiresAt, testCreatedAt, testUpdatedAt,
	)

	h.Mock.ExpectQuery(`SELECT reservations\.id AS "reservations\.id", reservations\.seat_id AS "reservations\.seat_id", reservations\.session_id AS "reservations\.session_id", reservations\.status AS "reservations\.status", reservations\.reserved_at AS "reservations\.reserved_at", reservations\.expires_at AS "reservations\.expires_at", reservations\.created_at AS "reservations\.created_at", reservations\.updated_at AS "reservations\.updated_at" FROM public\.reservations`).
		WillReturnRows(dataRows)

	ctx := context.Background()
	filter := repository.FindAllReservationsFilter{}

	// Execute
	reservations, total, err := h.Repository.FindAll(ctx, filter)

	// Assert - Should succeed but skip invalid records
	require.NoError(t, err)
	require.NotNil(t, reservations)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 0, len(*reservations)) // No valid reservations since ToEntity() returned nil

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
