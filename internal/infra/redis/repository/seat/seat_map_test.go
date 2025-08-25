package seatrepo_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domaincache "ticket-reservation/internal/domain/cache"
	"ticket-reservation/internal/domain/entity"
	seatrepo "ticket-reservation/internal/infra/redis/repository/seat"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
)

func TestNewSeatMapRepository(t *testing.T) {
	client, _ := redismock.NewClientMock()

	// Execute
	repo := seatrepo.NewSeatMapRepository(client)

	// Assert
	assert.NotNil(t, repo)
}

func TestSeatMapRepositoryImpl_SetSeat(t *testing.T) {
	concertID := uuid.New()
	zoneID := uuid.New()
	seatID := uuid.New()
	expectedKey := "seat_map:concert:" + concertID.String() + ":zone:" + zoneID.String()

	seat := entity.Seat{
		ID:                seatID,
		ZoneID:            zoneID,
		SeatNumber:        "A1",
		Status:            entity.SeatStatusPending,
		LockedUntil:       nil,
		LockedBySessionID: nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	seatJSON, _ := json.Marshal(seat)

	tests := []struct {
		name          string
		ttl           time.Duration
		setupMock     func(mock redismock.ClientMock)
		expectedError bool
		errorType     error
	}{
		{
			name: "successful set without TTL",
			ttl:  domaincache.SeatMapNoExpiration,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHSet(expectedKey, "A1", string(seatJSON)).SetVal(1)
			},
			expectedError: false,
		},
		{
			name: "successful set with TTL",
			ttl:  5 * time.Minute,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHSet(expectedKey, "A1", string(seatJSON)).SetVal(1)
				// Using ExpectDo for HExpire since redismock doesn't have ExpectHExpire yet
				// The actual Redis command includes FIELDS parameter
				mock.ExpectDo("HEXPIRE", expectedKey, int64(5*time.Minute/time.Second), "FIELDS", 1, "A1").SetVal([]int64{1})
			},
			expectedError: false,
		},
		{
			name: "hset operation fails",
			ttl:  domaincache.SeatMapNoExpiration,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHSet(expectedKey, "A1", string(seatJSON)).SetErr(errors.New("redis connection failed"))
			},
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "hexpire operation fails",
			ttl:  5 * time.Minute,
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHSet(expectedKey, "A1", string(seatJSON)).SetVal(1)
				// Using ExpectDo for HExpire since redismock doesn't have ExpectHExpire yet
				mock.ExpectDo("HEXPIRE", expectedKey, int64(5*time.Minute/time.Second), "FIELDS", 1, "A1").SetErr(errors.New("expire command failed"))
			},
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			repository := seatrepo.NewSeatMapRepository(client)

			tt.setupMock(mock)

			// Execute
			err := repository.SetSeat(context.Background(), concertID, zoneID, seat, tt.ttl)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[repository seat/seat_map SetSeat]")

				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSeatMapRepositoryImpl_GetSeat(t *testing.T) {
	concertID := uuid.New()
	zoneID := uuid.New()
	seatID := uuid.New()
	expectedKey := "seat_map:concert:" + concertID.String() + ":zone:" + zoneID.String()

	seat := entity.Seat{
		ID:                seatID,
		ZoneID:            zoneID,
		SeatNumber:        "A1",
		Status:            entity.SeatStatusBooked,
		LockedUntil:       nil,
		LockedBySessionID: nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	seatJSON, _ := json.Marshal(seat)

	tests := []struct {
		name          string
		seatNumber    string
		setupMock     func(mock redismock.ClientMock)
		expectedSeat  *entity.Seat
		expectedError bool
		errorType     error
		errorMessage  string
	}{
		{
			name:       "successful get seat",
			seatNumber: "A1",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGet(expectedKey, "A1").SetVal(string(seatJSON))
			},
			expectedSeat:  &seat,
			expectedError: false,
		},
		{
			name:       "seat not found - redis nil",
			seatNumber: "A1",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGet(expectedKey, "A1").RedisNil()
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.NotFoundError{},
			errorMessage:  "seat not found in cache",
		},
		{
			name:       "seat field expired - redis nil",
			seatNumber: "A1",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGet(expectedKey, "A1").SetErr(redis.Nil)
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.NotFoundError{},
			errorMessage:  "seat not found in cache",
		},
		{
			name:       "redis connection error",
			seatNumber: "A1",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGet(expectedKey, "A1").SetErr(errors.New("connection failed"))
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
			errorMessage:  "failed to get seat entity",
		},
		{
			name:       "invalid json in cache",
			seatNumber: "A1",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGet(expectedKey, "A1").SetVal("invalid-json")
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.InternalServerError{},
			errorMessage:  "failed to deserialize seat entity",
		},
		{
			name:       "empty json in cache",
			seatNumber: "A1",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGet(expectedKey, "A1").SetVal("")
			},
			expectedSeat:  nil,
			expectedError: true,
			errorType:     &errsFramework.InternalServerError{},
			errorMessage:  "failed to deserialize seat entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			repository := seatrepo.NewSeatMapRepository(client)

			tt.setupMock(mock)

			// Execute
			result, err := repository.GetSeat(context.Background(), concertID, zoneID, tt.seatNumber)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[repository seat/seat_map GetSeat]")
				assert.Contains(t, err.Error(), tt.errorMessage)

				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				// Compare seat values
				assert.Equal(t, tt.expectedSeat.ID, result.ID)
				assert.Equal(t, tt.expectedSeat.ZoneID, result.ZoneID)
				assert.Equal(t, tt.expectedSeat.SeatNumber, result.SeatNumber)
				assert.Equal(t, tt.expectedSeat.Status, result.Status)
				assert.Equal(t, tt.expectedSeat.LockedUntil, result.LockedUntil)
				assert.Equal(t, tt.expectedSeat.LockedBySessionID, result.LockedBySessionID)
			}
		})
	}
}

func TestSeatMapRepositoryImpl_GetAllSeats(t *testing.T) {
	concertID := uuid.New()
	zoneID := uuid.New()
	expectedKey := "seat_map:concert:" + concertID.String() + ":zone:" + zoneID.String()

	// Create test seats with different statuses
	lockedUntil := time.Now().Add(5 * time.Minute)
	sessionID := "session-123"

	seat1 := entity.Seat{
		ID:                uuid.New(),
		ZoneID:            zoneID,
		SeatNumber:        "A1",
		Status:            entity.SeatStatusAvailable,
		LockedUntil:       nil,
		LockedBySessionID: nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	seat2 := entity.Seat{
		ID:                uuid.New(),
		ZoneID:            zoneID,
		SeatNumber:        "A2",
		Status:            entity.SeatStatusPending,
		LockedUntil:       &lockedUntil,
		LockedBySessionID: &sessionID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	seat3 := entity.Seat{
		ID:                uuid.New(),
		ZoneID:            zoneID,
		SeatNumber:        "A3",
		Status:            entity.SeatStatusBooked,
		LockedUntil:       nil,
		LockedBySessionID: nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	seat1JSON, _ := json.Marshal(seat1)
	seat2JSON, _ := json.Marshal(seat2)
	seat3JSON, _ := json.Marshal(seat3)

	tests := []struct {
		name          string
		setupMock     func(mock redismock.ClientMock)
		expectedSeats *entity.Seats
		expectedError bool
		errorType     error
		errorMessage  string
		expectedCount int
	}{
		{
			name: "successful get all seats - multiple seats",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(expectedKey).SetVal(map[string]string{
					"A1": string(seat1JSON),
					"A2": string(seat2JSON),
					"A3": string(seat3JSON),
				})
			},
			expectedSeats: &entity.Seats{seat1, seat2, seat3},
			expectedError: false,
			expectedCount: 3,
		},
		{
			name: "successful get all seats - single seat",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(expectedKey).SetVal(map[string]string{
					"A1": string(seat1JSON),
				})
			},
			expectedSeats: &entity.Seats{seat1},
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "successful get all seats - empty hash",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(expectedKey).SetVal(map[string]string{})
			},
			expectedSeats: &entity.Seats{},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "redis connection error",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(expectedKey).SetErr(errors.New("connection failed"))
			},
			expectedSeats: nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
			errorMessage:  "failed to get all seat entities",
		},
		{
			name: "partial success - some seats have invalid json",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(expectedKey).SetVal(map[string]string{
					"A1": string(seat1JSON),
					"A2": "invalid-json",
					"A3": string(seat3JSON),
				})
			},
			expectedSeats: &entity.Seats{seat1, seat3}, // seat2 should be skipped due to invalid JSON
			expectedError: false,
			expectedCount: 2, // Only valid seats are returned
		},
		{
			name: "all seats have invalid json - returns empty slice",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(expectedKey).SetVal(map[string]string{
					"A1": "invalid-json-1",
					"A2": "invalid-json-2",
					"A3": "invalid-json-3",
				})
			},
			expectedSeats: &entity.Seats{},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			repository := seatrepo.NewSeatMapRepository(client)

			tt.setupMock(mock)

			// Execute
			result, err := repository.GetAllSeats(logger.NewContext(context.Background(), logger.NewNoopLogger()), concertID, zoneID)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[repository seat/seat_map GetAllSeats]")
				assert.Contains(t, err.Error(), tt.errorMessage)

				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, *result, tt.expectedCount)

				if tt.expectedCount > 0 {
					// Convert to maps for easier comparison
					resultMap := make(map[string]entity.Seat)
					expectedMap := make(map[string]entity.Seat)

					for _, seat := range *result {
						resultMap[seat.SeatNumber] = seat
					}

					for _, seat := range *tt.expectedSeats {
						expectedMap[seat.SeatNumber] = seat
					}

					assert.Equal(t, len(expectedMap), len(resultMap))

					for seatNumber, expectedSeat := range expectedMap {
						actualSeat, exists := resultMap[seatNumber]
						assert.True(t, exists, "Seat %s should exist in result", seatNumber)
						assert.Equal(t, expectedSeat.ID, actualSeat.ID)
						assert.Equal(t, expectedSeat.SeatNumber, actualSeat.SeatNumber)
						assert.Equal(t, expectedSeat.Status, actualSeat.Status)

						if expectedSeat.LockedUntil != nil && actualSeat.LockedUntil != nil {
							assert.True(t, expectedSeat.LockedUntil.Equal(*actualSeat.LockedUntil),
								"LockedUntil times should be equal for seat %s. Expected: %v, Got: %v",
								seatNumber, expectedSeat.LockedUntil, actualSeat.LockedUntil)
						} else {
							assert.Equal(t, expectedSeat.LockedUntil, actualSeat.LockedUntil)
						}
						assert.Equal(t, expectedSeat.LockedBySessionID, actualSeat.LockedBySessionID)
					}
				}
			}
		})
	}
}
