package handler

import (
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/usecase/concert"
	"ticket-reservation/internal/util/httpresponse"
	"time"

	"github.com/gin-gonic/gin"
)

type FindAllConcertsQuery struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02" time_utc:"7"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02" time_utc:"7"`
	Venue     *string    `form:"venue"`
}

type findAllConcertsResponse struct {
	ID    string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name  string `json:"name" example:"Concert Name"`
	Venue string `json:"venue" example:"Concert Venue"`
	Date  string `json:"date" example:"2025-01-01T10:00:00+07:00"`
}

// @Summary		List Concerts
// @Description	List all concerts, filterable by date range and venue
// @Tags			Concert
// @Produce		json
// @Param			startDate	query		string																		false	"Start date (format: 2006-01-02) (UTC+7)"
// @Param			endDate		query		string																		false	"End date (format: 2006-01-02) (UTC+7)"
// @Param			venue		query		string																		false	"Venue name (partial match)"
// @Success		200			{object}	httpresponse.SuccessResponse{data=[]findAllConcertsResponse,metadata=nil}	"List of concerts"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}										"Bad request"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}										"Internal server error"
// @Router			/concerts [get]
func (h *concertHandler) FindAllConcerts(c *gin.Context) {
	var query FindAllConcertsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpresponse.Error(c, err)
		return
	}

	concerts, err := h.concertUsecase.FindAllConcerts(c.Request.Context(), concert.FindAllConcertsInput{
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
		Venue:     query.Venue,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, h.newFindAllConcertsResponse(concerts))
}

func (h *concertHandler) newFindAllConcertsResponse(concerts *entity.Concerts) []findAllConcertsResponse {
	if concerts == nil || len(*concerts) == 0 {
		return []findAllConcertsResponse{}
	}

	loc, _ := time.LoadLocation(h.appConfig.Timezone)
	response := make([]findAllConcertsResponse, 0, len(*concerts))
	for _, concert := range *concerts {
		response = append(response, findAllConcertsResponse{
			ID:    concert.ID.String(),
			Name:  concert.Name,
			Venue: concert.Venue,
			Date:  concert.Date.In(loc).Format(time.RFC3339),
		})
	}
	return response
}
