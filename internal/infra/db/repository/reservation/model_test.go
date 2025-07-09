package reservationrepo_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"
	reservationrepo "ticket-reservation/internal/infra/db/repository/reservation"

	"github.com/kittipat1413/go-common/util/pointer"
)

func TestReservation_ToEntity(t *testing.T) {
	testID := uuid.New()
	testSeatID := uuid.New()
	testSessionID := "session-123"
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          reservationrepo.Reservation
		expectedEntity *entity.Reservation
		expectedNil    bool
	}{
		{
			name: "successful conversion with pending status",
			input: reservationrepo.Reservation{
				Reservations: model.Reservations{
					ID:         testID,
					SeatID:     testSeatID,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusPending.String(),
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedEntity: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     entity.ReservationStatusPending,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "successful conversion with confirmed status",
			input: reservationrepo.Reservation{
				Reservations: model.Reservations{
					ID:         testID,
					SeatID:     testSeatID,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusConfirmed.String(),
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedEntity: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     entity.ReservationStatusConfirmed,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "successful conversion with expired status",
			input: reservationrepo.Reservation{
				Reservations: model.Reservations{
					ID:         testID,
					SeatID:     testSeatID,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusExpired.String(),
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedEntity: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     entity.ReservationStatusExpired,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "invalid status returns nil",
			input: reservationrepo.Reservation{
				Reservations: model.Reservations{
					ID:         testID,
					SeatID:     testSeatID,
					SessionID:  testSessionID,
					Status:     "invalid_status",
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedEntity: nil,
			expectedNil:    true,
		},
		{
			name: "empty status returns nil",
			input: reservationrepo.Reservation{
				Reservations: model.Reservations{
					ID:         testID,
					SeatID:     testSeatID,
					SessionID:  testSessionID,
					Status:     "",
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedEntity: nil,
			expectedNil:    true,
		},
		{
			name: "conversion with zero time values",
			input: reservationrepo.Reservation{
				Reservations: model.Reservations{
					ID:         testID,
					SeatID:     testSeatID,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusPending.String(),
					ReservedAt: time.Time{},
					ExpiresAt:  time.Time{},
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
				},
			},
			expectedEntity: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  testSessionID,
				Status:     entity.ReservationStatusPending,
				ReservedAt: time.Time{},
				ExpiresAt:  time.Time{},
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
			},
			expectedNil: false,
		},
		{
			name: "conversion with empty session ID",
			input: reservationrepo.Reservation{
				Reservations: model.Reservations{
					ID:         testID,
					SeatID:     testSeatID,
					SessionID:  "",
					Status:     entity.ReservationStatusPending.String(),
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedEntity: &entity.Reservation{
				ID:         testID,
				SeatID:     testSeatID,
				SessionID:  "",
				Status:     entity.ReservationStatusPending,
				ReservedAt: testReservedAt,
				ExpiresAt:  testExpiresAt,
				CreatedAt:  testCreatedAt,
				UpdatedAt:  testUpdatedAt,
			},
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			result := tt.input.ToEntity()

			// Assert
			if tt.expectedNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedEntity.ID, result.ID)
				assert.Equal(t, tt.expectedEntity.SeatID, result.SeatID)
				assert.Equal(t, tt.expectedEntity.SessionID, result.SessionID)
				assert.Equal(t, tt.expectedEntity.Status, result.Status)
				assert.Equal(t, tt.expectedEntity.ReservedAt.UTC(), result.ReservedAt.UTC())
				assert.Equal(t, tt.expectedEntity.ExpiresAt.UTC(), result.ExpiresAt.UTC())
				assert.Equal(t, tt.expectedEntity.CreatedAt.UTC(), result.CreatedAt.UTC())
				assert.Equal(t, tt.expectedEntity.UpdatedAt.UTC(), result.UpdatedAt.UTC())
			}
		})
	}
}

func TestReservations_ToEntities(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()
	testID3 := uuid.New()
	testSeatID1 := uuid.New()
	testSeatID2 := uuid.New()
	testSeatID3 := uuid.New()
	testSessionID := "session-123"
	testReservedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testExpiresAt := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name               string
		input              reservationrepo.Reservations
		expectedEntities   *entity.Reservations
		expectedLength     int
		shouldSkipInvalids bool
	}{
		{
			name:             "empty slice",
			input:            reservationrepo.Reservations{},
			expectedEntities: &entity.Reservations{},
			expectedLength:   0,
		},
		{
			name: "single reservation",
			input: reservationrepo.Reservations{
				{
					Reservations: model.Reservations{
						ID:         testID1,
						SeatID:     testSeatID1,
						SessionID:  testSessionID,
						Status:     entity.ReservationStatusPending.String(),
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Reservations{
				{
					ID:         testID1,
					SeatID:     testSeatID1,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusPending,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedLength: 1,
		},
		{
			name: "multiple reservations",
			input: reservationrepo.Reservations{
				{
					Reservations: model.Reservations{
						ID:         testID1,
						SeatID:     testSeatID1,
						SessionID:  testSessionID,
						Status:     entity.ReservationStatusPending.String(),
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
				{
					Reservations: model.Reservations{
						ID:         testID2,
						SeatID:     testSeatID2,
						SessionID:  testSessionID,
						Status:     entity.ReservationStatusConfirmed.String(),
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Reservations{
				{
					ID:         testID1,
					SeatID:     testSeatID1,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusPending,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
				{
					ID:         testID2,
					SeatID:     testSeatID2,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusConfirmed,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedLength: 2,
		},
		{
			name: "mixed valid and invalid reservations",
			input: reservationrepo.Reservations{
				{
					Reservations: model.Reservations{
						ID:         testID1,
						SeatID:     testSeatID1,
						SessionID:  testSessionID,
						Status:     entity.ReservationStatusPending.String(),
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
				{
					Reservations: model.Reservations{
						ID:         testID2,
						SeatID:     testSeatID2,
						SessionID:  testSessionID,
						Status:     "invalid_status", // Invalid status
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
				{
					Reservations: model.Reservations{
						ID:         testID3,
						SeatID:     testSeatID3,
						SessionID:  testSessionID,
						Status:     entity.ReservationStatusConfirmed.String(),
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Reservations{
				{
					ID:         testID1,
					SeatID:     testSeatID1,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusPending,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
				{
					ID:         testID3,
					SeatID:     testSeatID3,
					SessionID:  testSessionID,
					Status:     entity.ReservationStatusConfirmed,
					ReservedAt: testReservedAt,
					ExpiresAt:  testExpiresAt,
					CreatedAt:  testCreatedAt,
					UpdatedAt:  testUpdatedAt,
				},
			},
			expectedLength:     2, // Only valid reservations should be included
			shouldSkipInvalids: true,
		},
		{
			name: "all invalid reservations",
			input: reservationrepo.Reservations{
				{
					Reservations: model.Reservations{
						ID:         testID1,
						SeatID:     testSeatID1,
						SessionID:  testSessionID,
						Status:     "invalid_status_1",
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
				{
					Reservations: model.Reservations{
						ID:         testID2,
						SeatID:     testSeatID2,
						SessionID:  testSessionID,
						Status:     "invalid_status_2",
						ReservedAt: testReservedAt,
						ExpiresAt:  testExpiresAt,
						CreatedAt:  testCreatedAt,
						UpdatedAt:  testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Reservations{},
			expectedLength:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			result := tt.input.ToEntities()

			// Assert
			require.NotNil(t, result)
			assert.Equal(t, tt.expectedLength, len(*result))

			// Compare each reservation
			for i, expected := range *tt.expectedEntities {
				actual := (*result)[i]
				assert.Equal(t, expected.ID, actual.ID)
				assert.Equal(t, expected.SeatID, actual.SeatID)
				assert.Equal(t, expected.SessionID, actual.SessionID)
				assert.Equal(t, expected.Status, actual.Status)
				assert.Equal(t, expected.ReservedAt.UTC(), actual.ReservedAt.UTC())
				assert.Equal(t, expected.ExpiresAt.UTC(), actual.ExpiresAt.UTC())
				assert.Equal(t, expected.CreatedAt.UTC(), actual.CreatedAt.UTC())
				assert.Equal(t, expected.UpdatedAt.UTC(), actual.UpdatedAt.UTC())
			}
		})
	}
}

func TestReservations_ToEntities_EmptyAndNilChecks(t *testing.T) {
	tests := []struct {
		name   string
		input  reservationrepo.Reservations
		assert func(t *testing.T, result *entity.Reservations)
	}{
		{
			name:  "nil input",
			input: nil,
			assert: func(t *testing.T, result *entity.Reservations) {
				require.NotNil(t, result)
				assert.Equal(t, 0, len(*result))
			},
		},
		{
			name:  "empty slice",
			input: reservationrepo.Reservations{},
			assert: func(t *testing.T, result *entity.Reservations) {
				require.NotNil(t, result)
				assert.Equal(t, 0, len(*result))
			},
		},
		{
			name: "slice with nil ToEntity results",
			input: reservationrepo.Reservations{
				{
					Reservations: model.Reservations{
						ID:         uuid.New(),
						SeatID:     uuid.New(),
						SessionID:  "session-123",
						Status:     "", // Empty status will cause ToEntity to return nil
						ReservedAt: time.Now(),
						ExpiresAt:  time.Now(),
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					},
				},
			},
			assert: func(t *testing.T, result *entity.Reservations) {
				require.NotNil(t, result)
				assert.Equal(t, 0, len(*result)) // Should skip nil results
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToEntities()
			tt.assert(t, result)
		})
	}
}

func TestReservations_ToEntities_ReturnType(t *testing.T) {
	// Test that ToEntities always returns a pointer to entity.Reservations
	input := reservationrepo.Reservations{}
	result := input.ToEntities()

	// Check that it's a pointer
	assert.IsType(t, &entity.Reservations{}, result)
	assert.NotNil(t, result)

	// Check that it's not a nil pointer
	assert.NotNil(t, pointer.ToPointer(entity.Reservations{}))
}
