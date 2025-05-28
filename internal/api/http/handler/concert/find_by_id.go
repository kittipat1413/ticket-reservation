package handler

import (
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/util/httpresponse"
	"time"

	concertUsecase "ticket-reservation/internal/usecase/concert"

	"github.com/gin-gonic/gin"
)

type findOneConcertResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name  string `json:"name" example:"Concert Name"`
	Venue string `json:"venue" example:"Concert Venue"`
	Date  string `json:"date" example:"2025-01-01T10:00:00+07:00"`
}

// @Summary		Find Concert by ID
// @Description	Retrieve concert details by its ID
// @Tags			Concert
// @Produce		json
// @Param			id	path		string						true	"Concert ID"
// @Success		200	{object}	findOneConcertResponse		"Concert found"
// @Failure		404	{object}	httpresponse.ErrorResponse	"Concert not found"
// @Failure		500	{object}	httpresponse.ErrorResponse	"Internal server error"
// @Router			/concerts/{id} [get]
func (h *concertHandler) FindConcertByID(c *gin.Context) {
	concertID := c.Param("id")
	concert, err := h.concertUsecase.FindOneConcert(c.Request.Context(), concertUsecase.FindOneConcertInput{
		ID: concertID,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindOneConcertResponse(concert))
}

func (h *concertHandler) newFindOneConcertResponse(concert *entity.Concert) findOneConcertResponse {
	if concert == nil {
		return findOneConcertResponse{}
	}

	loc, _ := time.LoadLocation(h.appConfig.Timezone)
	return findOneConcertResponse{
		ID:    concert.ID.String(),
		Name:  concert.Name,
		Venue: concert.Venue,
		Date:  concert.Date.In(loc).Format(time.RFC3339),
	}
}
