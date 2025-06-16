package cache

import "fmt"

const (
	SeatLockCacheKeyFormat = "seat_lock:concert:%s:zone:%s:seat:%s" // Format: seat_lock:concert:<concert_id>:zone:<zone_id>:seat:<seat_id>
)

func GetSeatLockKey(concertID, zoneID, seatID string) string {
	return fmt.Sprintf(SeatLockCacheKeyFormat, concertID, zoneID, seatID)
}
