package seatrepo_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"
	seatrepo "ticket-reservation/internal/infra/db/repository/seat"

	"github.com/kittipat1413/go-common/util/pointer"
)

func TestSeat_ToEntity(t *testing.T) {
	testID := uuid.New()
	testZoneID := uuid.New()
	testSeatNumber := "A1"
	testSessionID := "session-123"
	testLockedUntil := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          seatrepo.Seat
		expectedEntity *entity.Seat
		expectedNil    bool
	}{
		{
			name: "successful conversion with available status",
			input: seatrepo.Seat{
				Seats: model.Seats{
					ID:                testID,
					ZoneID:            testZoneID,
					SeatNumber:        testSeatNumber,
					Status:            entity.SeatStatusAvailable.String(),
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedEntity: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusAvailable,
				LockedUntil:       nil,
				LockedBySessionID: nil,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "successful conversion with pending status",
			input: seatrepo.Seat{
				Seats: model.Seats{
					ID:                testID,
					ZoneID:            testZoneID,
					SeatNumber:        testSeatNumber,
					Status:            entity.SeatStatusPending.String(),
					LockedUntil:       &testLockedUntil,
					LockedBySessionID: &testSessionID,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedEntity: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusPending,
				LockedUntil:       &testLockedUntil,
				LockedBySessionID: &testSessionID,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "successful conversion with booked status",
			input: seatrepo.Seat{
				Seats: model.Seats{
					ID:                testID,
					ZoneID:            testZoneID,
					SeatNumber:        testSeatNumber,
					Status:            entity.SeatStatusBooked.String(),
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedEntity: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusBooked,
				LockedUntil:       nil,
				LockedBySessionID: nil,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "invalid status returns nil",
			input: seatrepo.Seat{
				Seats: model.Seats{
					ID:                testID,
					ZoneID:            testZoneID,
					SeatNumber:        testSeatNumber,
					Status:            "invalid_status",
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedEntity: nil,
			expectedNil:    true,
		},
		{
			name: "empty status returns nil",
			input: seatrepo.Seat{
				Seats: model.Seats{
					ID:                testID,
					ZoneID:            testZoneID,
					SeatNumber:        testSeatNumber,
					Status:            "",
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedEntity: nil,
			expectedNil:    true,
		},
		{
			name: "conversion with zero time values",
			input: seatrepo.Seat{
				Seats: model.Seats{
					ID:                testID,
					ZoneID:            testZoneID,
					SeatNumber:        testSeatNumber,
					Status:            entity.SeatStatusAvailable.String(),
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				},
			},
			expectedEntity: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        testSeatNumber,
				Status:            entity.SeatStatusAvailable,
				LockedUntil:       nil,
				LockedBySessionID: nil,
				CreatedAt:         time.Time{},
				UpdatedAt:         time.Time{},
			},
			expectedNil: false,
		},
		{
			name: "conversion with empty seat number",
			input: seatrepo.Seat{
				Seats: model.Seats{
					ID:                testID,
					ZoneID:            testZoneID,
					SeatNumber:        "",
					Status:            entity.SeatStatusAvailable.String(),
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedEntity: &entity.Seat{
				ID:                testID,
				ZoneID:            testZoneID,
				SeatNumber:        "",
				Status:            entity.SeatStatusAvailable,
				LockedUntil:       nil,
				LockedBySessionID: nil,
				CreatedAt:         testCreatedAt,
				UpdatedAt:         testUpdatedAt,
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
				assert.Equal(t, tt.expectedEntity.ZoneID, result.ZoneID)
				assert.Equal(t, tt.expectedEntity.SeatNumber, result.SeatNumber)
				assert.Equal(t, tt.expectedEntity.Status, result.Status)

				if tt.expectedEntity.LockedUntil != nil {
					require.NotNil(t, result.LockedUntil)
					assert.Equal(t, tt.expectedEntity.LockedUntil.UTC(), result.LockedUntil.UTC())
				} else {
					assert.Nil(t, result.LockedUntil)
				}

				if tt.expectedEntity.LockedBySessionID != nil {
					require.NotNil(t, result.LockedBySessionID)
					assert.Equal(t, *tt.expectedEntity.LockedBySessionID, *result.LockedBySessionID)
				} else {
					assert.Nil(t, result.LockedBySessionID)
				}

				assert.Equal(t, tt.expectedEntity.CreatedAt.UTC(), result.CreatedAt.UTC())
				assert.Equal(t, tt.expectedEntity.UpdatedAt.UTC(), result.UpdatedAt.UTC())
			}
		})
	}
}

func TestSeats_ToEntities(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()
	testID3 := uuid.New()
	testZoneID1 := uuid.New()
	testZoneID2 := uuid.New()
	testZoneID3 := uuid.New()
	testSeatNumber1 := "A1"
	testSeatNumber2 := "A2"
	testSeatNumber3 := "A3"
	testSessionID := "session-123"
	testLockedUntil := time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		input            seatrepo.Seats
		expectedEntities *entity.Seats
		expectedLength   int
	}{
		{
			name:             "empty slice",
			input:            seatrepo.Seats{},
			expectedEntities: &entity.Seats{},
			expectedLength:   0,
		},
		{
			name: "single seat",
			input: seatrepo.Seats{
				{
					Seats: model.Seats{
						ID:                testID1,
						ZoneID:            testZoneID1,
						SeatNumber:        testSeatNumber1,
						Status:            entity.SeatStatusAvailable.String(),
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Seats{
				{
					ID:                testID1,
					ZoneID:            testZoneID1,
					SeatNumber:        testSeatNumber1,
					Status:            entity.SeatStatusAvailable,
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedLength: 1,
		},
		{
			name: "multiple seats",
			input: seatrepo.Seats{
				{
					Seats: model.Seats{
						ID:                testID1,
						ZoneID:            testZoneID1,
						SeatNumber:        testSeatNumber1,
						Status:            entity.SeatStatusAvailable.String(),
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
				{
					Seats: model.Seats{
						ID:                testID2,
						ZoneID:            testZoneID2,
						SeatNumber:        testSeatNumber2,
						Status:            entity.SeatStatusPending.String(),
						LockedUntil:       &testLockedUntil,
						LockedBySessionID: &testSessionID,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Seats{
				{
					ID:                testID1,
					ZoneID:            testZoneID1,
					SeatNumber:        testSeatNumber1,
					Status:            entity.SeatStatusAvailable,
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
				{
					ID:                testID2,
					ZoneID:            testZoneID2,
					SeatNumber:        testSeatNumber2,
					Status:            entity.SeatStatusPending,
					LockedUntil:       &testLockedUntil,
					LockedBySessionID: &testSessionID,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedLength: 2,
		},
		{
			name: "mixed valid and invalid seats",
			input: seatrepo.Seats{
				{
					Seats: model.Seats{
						ID:                testID1,
						ZoneID:            testZoneID1,
						SeatNumber:        testSeatNumber1,
						Status:            entity.SeatStatusAvailable.String(),
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
				{
					Seats: model.Seats{
						ID:                testID2,
						ZoneID:            testZoneID2,
						SeatNumber:        testSeatNumber2,
						Status:            "invalid_status", // Invalid status
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
				{
					Seats: model.Seats{
						ID:                testID3,
						ZoneID:            testZoneID3,
						SeatNumber:        testSeatNumber3,
						Status:            entity.SeatStatusBooked.String(),
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Seats{
				{
					ID:                testID1,
					ZoneID:            testZoneID1,
					SeatNumber:        testSeatNumber1,
					Status:            entity.SeatStatusAvailable,
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
				{
					ID:                testID3,
					ZoneID:            testZoneID3,
					SeatNumber:        testSeatNumber3,
					Status:            entity.SeatStatusBooked,
					LockedUntil:       nil,
					LockedBySessionID: nil,
					CreatedAt:         testCreatedAt,
					UpdatedAt:         testUpdatedAt,
				},
			},
			expectedLength: 2, // Only valid seats should be included
		},
		{
			name: "all invalid seats",
			input: seatrepo.Seats{
				{
					Seats: model.Seats{
						ID:                testID1,
						ZoneID:            testZoneID1,
						SeatNumber:        testSeatNumber1,
						Status:            "invalid_status_1",
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
				{
					Seats: model.Seats{
						ID:                testID2,
						ZoneID:            testZoneID2,
						SeatNumber:        testSeatNumber2,
						Status:            "invalid_status_2",
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         testCreatedAt,
						UpdatedAt:         testUpdatedAt,
					},
				},
			},
			expectedEntities: &entity.Seats{},
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

			// Compare each seat
			for i, expected := range *tt.expectedEntities {
				actual := (*result)[i]
				assert.Equal(t, expected.ID, actual.ID)
				assert.Equal(t, expected.ZoneID, actual.ZoneID)
				assert.Equal(t, expected.SeatNumber, actual.SeatNumber)
				assert.Equal(t, expected.Status, actual.Status)

				if expected.LockedUntil != nil {
					require.NotNil(t, actual.LockedUntil)
					assert.Equal(t, expected.LockedUntil.UTC(), actual.LockedUntil.UTC())
				} else {
					assert.Nil(t, actual.LockedUntil)
				}

				if expected.LockedBySessionID != nil {
					require.NotNil(t, actual.LockedBySessionID)
					assert.Equal(t, *expected.LockedBySessionID, *actual.LockedBySessionID)
				} else {
					assert.Nil(t, actual.LockedBySessionID)
				}

				assert.Equal(t, expected.CreatedAt.UTC(), actual.CreatedAt.UTC())
				assert.Equal(t, expected.UpdatedAt.UTC(), actual.UpdatedAt.UTC())
			}
		})
	}
}

func TestSeats_ToEntities_EmptyAndNilChecks(t *testing.T) {
	tests := []struct {
		name   string
		input  seatrepo.Seats
		assert func(t *testing.T, result *entity.Seats)
	}{
		{
			name:  "nil input",
			input: nil,
			assert: func(t *testing.T, result *entity.Seats) {
				require.NotNil(t, result)
				assert.Equal(t, 0, len(*result))
			},
		},
		{
			name:  "empty slice",
			input: seatrepo.Seats{},
			assert: func(t *testing.T, result *entity.Seats) {
				require.NotNil(t, result)
				assert.Equal(t, 0, len(*result))
			},
		},
		{
			name: "slice with nil ToEntity results",
			input: seatrepo.Seats{
				{
					Seats: model.Seats{
						ID:                uuid.New(),
						ZoneID:            uuid.New(),
						SeatNumber:        "A1",
						Status:            "", // Empty status will cause ToEntity to return nil
						LockedUntil:       nil,
						LockedBySessionID: nil,
						CreatedAt:         time.Now(),
						UpdatedAt:         time.Now(),
					},
				},
			},
			assert: func(t *testing.T, result *entity.Seats) {
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

func TestSeats_ToEntities_ReturnType(t *testing.T) {
	// Test that ToEntities always returns a pointer to entity.Seats
	input := seatrepo.Seats{}
	result := input.ToEntities()

	// Check that it's a pointer
	assert.IsType(t, &entity.Seats{}, result)
	assert.NotNil(t, result)

	// Check that it's not a nil pointer
	assert.NotNil(t, pointer.ToPointer(entity.Seats{}))
}
