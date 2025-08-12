package cache

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	SeatLockCacheKeyFormat = "seat_lock:concert:%s:zone:%s:seat:%s" // Format: seat_lock:concert:<concert_id>:zone:<zone_id>:seat:<seat_id>
	SeatMapCacheKeyFormat  = "seat_map:concert:%s:zone:%s"          // Format: seat_map:concert:<concert_id>:zone:<zone_id>
)

func GetSeatLockKey(concertID, zoneID, seatID uuid.UUID) string {
	return fmt.Sprintf(SeatLockCacheKeyFormat, concertID.String(), zoneID.String(), seatID.String())
}

func GetSeatMapKey(concertID, zoneID uuid.UUID) string {
	return fmt.Sprintf(SeatMapCacheKeyFormat, concertID.String(), zoneID.String())
}
