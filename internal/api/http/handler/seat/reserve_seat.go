package handler

import (
	"ticket-reservation/internal/domain/entity"
	seatUsecase "ticket-reservation/internal/usecase/seat"
	"ticket-reservation/internal/util/httpresponse"
	"time"

	"github.com/gin-gonic/gin"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

type ReserveSeatRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

type ReserveSeatResponse struct {
	ReservationID string `json:"reservation_id"`
	SeatID        string `json:"seat_id"`
	Status        string `json:"status"`
	ReservedAt    string `json:"reserved_at"`
	ExpiresAt     string `json:"expires_at"`
}

// @Summary		Reserve a Seat
// @Description	Reserves a seat for a concert by locking it for the current session
// @Tags			Seat
// @Accept			json
// @Produce		json
// @Param			id			path		string																true	"Concert ID"
// @Param			zone_id		path		string																true	"Zone ID"
// @Param			seat_number	path		string																true	"Seat Number"
// @Param			request		body		ReserveSeatRequest													true	"Reservation Request"
// @Success		200			{object}	httpresponse.SuccessResponse{data=ReserveSeatResponse,metadata=nil}	"Seat reserved successfully"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}								"Bad Request - Invalid input"
// @Failure		409			{object}	httpresponse.ErrorResponse{data=nil}								"Conflict - Seat already reserved"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}								"Internal Server Error - Unexpected error occurred"
// @Router			/concerts/{id}/zones/{zone_id}/seats/{seat_number}/reserve [post]
func (h *seatHandler) ReserveSeat(c *gin.Context) {
	var request ReserveSeatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		err = errsFramework.WrapError(err, errsFramework.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	result, err := h.seatUsecase.ReserveSeat(c.Request.Context(), seatUsecase.ReserveSeatInput{
		ConcertID: c.Param("id"),
		ZoneID:    c.Param("zone_id"),
		SeatID:    c.Param("seat_id"),
		SessionID: request.SessionID,
	})

	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newReserveSeatResponse(result))
}

func (h *seatHandler) newReserveSeatResponse(reservation *entity.Reservation) ReserveSeatResponse {
	loc, _ := time.LoadLocation(h.appConfig.Timezone)
	return ReserveSeatResponse{
		ReservationID: reservation.ID.String(),
		SeatID:        reservation.SeatID.String(),
		Status:        reservation.Status.String(),
		ReservedAt:    reservation.ReservedAt.In(loc).Format(time.RFC3339),
		ExpiresAt:     reservation.ExpiresAt.In(loc).Format(time.RFC3339),
	}
}
