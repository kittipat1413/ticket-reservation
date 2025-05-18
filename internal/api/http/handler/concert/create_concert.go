package handler

import (
	"net/http"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/usecase/concert"
	"ticket-reservation/internal/util/httpresponse"
	"time"

	"github.com/gin-gonic/gin"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

type createConcertRequest struct {
	Name  string    `json:"name" example:"Concert Name" binding:"required"`
	Venue string    `json:"venue" example:"Concert Venue" binding:"required"`
	Date  time.Time `json:"date" example:"2025-01-01T10:00:00+07:00" binding:"required"`
}

type createConcertResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name  string `json:"name" example:"Concert Name"`
	Venue string `json:"venue" example:"Concert Venue"`
	Date  string `json:"date" example:"2025-01-01T10:00:00+07:00"`
}

//	@Summary		Create Concert
//	@Description	Create a new concert
//	@Tags			Concert
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createConcertRequest		true	"Concert creation input"
//	@Success		201		{object}	createConcertResponse		"Concert created"
//	@Failure		400		{object}	httpresponse.ErrorResponse	"Bad request"
//	@Failure		500		{object}	httpresponse.ErrorResponse	"Internal server error"
//	@Router			/concerts [post]
func (h *concertHandler) CreateConcert(c *gin.Context) {
	var input createConcertRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		err = errsFramework.WrapError(err, errsFramework.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	createdConcert, err := h.concertUsecase.CreateConcert(c.Request.Context(), concert.CreateConcertInput{
		Name:  input.Name,
		Venue: input.Venue,
		Date:  input.Date,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateConcertResponse(createdConcert))
}

func (h *concertHandler) newCreateConcertResponse(concert *entity.Concert) createConcertResponse {
	if concert == nil {
		return createConcertResponse{}
	}

	loc, _ := time.LoadLocation(h.appConfig.Timezone)
	return createConcertResponse{
		ID:    concert.ID.String(),
		Name:  concert.Name,
		Venue: concert.Venue,
		Date:  concert.Date.In(loc).Format(time.RFC3339),
	}
}
