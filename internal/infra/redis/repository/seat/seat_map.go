package seatrepo

import (
	"context"
	"encoding/json"
	"errors"
	domaincache "ticket-reservation/internal/domain/cache"
	"ticket-reservation/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/framework/logger"
	"github.com/kittipat1413/go-common/util/pointer"
	"github.com/redis/go-redis/v9"
)

type seatMap struct {
	redisClient redis.UniversalClient
}

func NewSeatMapRepository(redisClient redis.UniversalClient) domaincache.SeatMapRepository {
	return &seatMap{
		redisClient: redisClient,
	}
}

func (r *seatMap) SetSeat(ctx context.Context, concertID, zoneID uuid.UUID, seat entity.Seat, ttl time.Duration) (err error) {
	const errLocation = "[repository seat/seat_map SetSeat]"
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	key := getSeatMapKey(concertID, zoneID)

	// Serialize seat entity to JSON
	seatJSON, err := json.Marshal(seat)
	if err != nil {
		return errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to serialize seat entity", nil))
	}

	// Set the field value
	err = r.redisClient.HSet(ctx, key, seat.SeatNumber, string(seatJSON)).Err()
	if err != nil {
		return errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to set seat entity", err.Error()))
	}

	// Set TTL for the specific field if provided
	if ttl > domaincache.SeatMapNoExpiration { // ttl > 0
		err = r.redisClient.HExpire(ctx, key, ttl, seat.SeatNumber).Err()
		if err != nil {
			return errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to set TTL for seat entity", err.Error()))
		}
	}

	return nil
}

func (r *seatMap) GetSeat(ctx context.Context, concertID, zoneID uuid.UUID, seatNumber string) (seat *entity.Seat, err error) {
	const errLocation = "[repository seat/seat_map GetSeat]"
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	key := getSeatMapKey(concertID, zoneID)

	// Get the seat from the Redis hash
	result, err := r.redisClient.HGet(ctx, key, seatNumber).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Field not found or expired
			return nil, errsFramework.NewNotFoundError("seat not found in cache", nil)
		}
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to get seat entity", err.Error()))
	}

	// Deserialize JSON to seat entity
	var seatEntity entity.Seat
	err = json.Unmarshal([]byte(result), &seatEntity)
	if err != nil {
		return nil, errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to deserialize seat entity", nil))
	}

	return pointer.ToPointer(seatEntity), nil
}

func (r *seatMap) GetAllSeats(ctx context.Context, concertID, zoneID uuid.UUID) (seats *entity.Seats, err error) {
	const errLocation = "[repository seat/seat_map GetAllSeats]"
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	key := getSeatMapKey(concertID, zoneID)

	result, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to get all seat entities", err.Error()))
	}

	seatsList := make(entity.Seats, 0, len(result))
	for seatNumber, seatJSON := range result {
		// Deserialize JSON to seat entity
		var seatEntity entity.Seat
		if err := json.Unmarshal([]byte(seatJSON), &seatEntity); err != nil {
			// Log error and continue processing other seats
			logger.FromContext(ctx).Error(ctx, "failed to deserialize seat entity", err, logger.Fields{
				"seat_number": seatNumber,
				"seat_json":   seatJSON,
			})
			continue
		}
		// Append the successfully deserialized seat to the slice
		seatsList = append(seatsList, seatEntity)
	}

	return pointer.ToPointer(seatsList), nil
}

func getSeatMapKey(concertID, zoneID uuid.UUID) string {
	return domaincache.GetSeatMapKey(concertID, zoneID)
}
