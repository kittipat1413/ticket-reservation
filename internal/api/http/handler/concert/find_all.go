package handler

import (
	"ticket-reservation/internal/domain/entity"
	concertUsecase "ticket-reservation/internal/usecase/concert"
	"ticket-reservation/internal/util/httpresponse"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kittipat1413/go-common/util/pointer"
)

type FindAllConcertsQuery struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02" time_utc:"7"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02" time_utc:"7"`
	Venue     *string    `form:"venue"`
	Limit     *int64     `form:"limit"`
	Offset    *int64     `form:"offset"`
	SortBy    *string    `form:"sortBy"`
	SortOrder *string    `form:"sortOrder"`
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
// @Param			startDate	query		string																									false	"Start date (format: 2006-01-02) (UTC+7)"
// @Param			endDate		query		string																									false	"End date (format: 2006-01-02) (UTC+7)"
// @Param			venue		query		string																									false	"Venue name (partial match)"
// @Param			limit		query		int64																									false	"Number of results to return (default: 100)"
// @Param			offset		query		int64																									false	"Number of results to skip (default: 0)"
// @Param			sortBy		query		string																									false	"Field to sort by (default: date) (options: date, name, venue)"
// @Param			sortOrder	query		string																									false	"Sort order (default: asc) (options: asc, desc)"
// @Success		200			{object}	httpresponse.SuccessResponse{data=[]findAllConcertsResponse,metadata=httpresponse.PaginationMetadata}	"List of concerts with pagination details"
// @Failure		400			{object}	httpresponse.ErrorResponse{data=nil}																	"Bad request"
// @Failure		500			{object}	httpresponse.ErrorResponse{data=nil}																	"Internal server error"
// @Router			/concerts [get]
func (h *concertHandler) FindAllConcerts(c *gin.Context) {
	var query FindAllConcertsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpresponse.Error(c, err)
		return
	}

	var (
		limit     = pointer.ToPointer(int64(100))          // Default limit to 100
		offset    = pointer.ToPointer(int64(0))            // Default offset to 0
		sortBy    = pointer.ToPointer("date")              // Default sort by date
		sortOrder = pointer.ToPointer(entity.SortOrderAsc) // Default sort order
		err       error
	)
	if query.Limit != nil {
		limit = query.Limit
	}
	if query.Offset != nil {
		offset = query.Offset
	}
	if query.SortBy != nil {
		sortBy = query.SortBy
	}
	if query.SortOrder != nil {
		querySortOrder, err := entity.ParseSortOrder(pointer.GetValue(query.SortOrder))
		if err == nil {
			sortOrder = pointer.ToPointer(querySortOrder)
		}
	}

	concerts, err := h.concertUsecase.FindAllConcerts(c.Request.Context(), concertUsecase.FindAllConcertsInput{
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
		Venue:     query.Venue,
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithMetadata(c, h.newFindAllConcertsResponse(concerts.GetData()), httpresponse.PaginationMetadata{Pagination: concerts.GetPagination()})
}

func (h *concertHandler) newFindAllConcertsResponse(concerts entity.Concerts) []findAllConcertsResponse {
	if len(concerts) == 0 {
		return []findAllConcertsResponse{}
	}

	loc, _ := time.LoadLocation(h.appConfig.Timezone)
	response := make([]findAllConcertsResponse, 0, len(concerts))
	for _, concert := range concerts {
		response = append(response, findAllConcertsResponse{
			ID:    concert.ID.String(),
			Name:  concert.Name,
			Venue: concert.Venue,
			Date:  concert.Date.In(loc).Format(time.RFC3339),
		})
	}
	return response
}
