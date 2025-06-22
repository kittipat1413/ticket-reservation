package usecase

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	customvalidator "ticket-reservation/pkg/validator"
	"time"

	repository "ticket-reservation/internal/domain/repository"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	traceFramework "github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/framework/validator"
	"github.com/kittipat1413/go-common/util/pointer"
)

type FindAllConcertsInput struct {
	StartDate *time.Time        `json:"start_date" validate:"omitempty,thaitimezone"`
	EndDate   *time.Time        `json:"end_date" validate:"omitempty,thaitimezone,gtfield=StartDate"`
	Venue     *string           `json:"venue" validate:"omitempty,gt=0"`
	Limit     *int64            `json:"limit" validate:"required,gte=1,lte=100"`
	Offset    *int64            `json:"offset" validate:"required,gte=0"`
	SortBy    *string           `json:"sort_by" validate:"required_with=SortOrder,omitempty,oneof=date name venue"`
	SortOrder *entity.SortOrder `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

func (u *concertUsecase) FindAllConcerts(ctx context.Context, input FindAllConcertsInput) (concerts entity.Page[entity.Concert], err error) {
	const errLocation = "[usecase concert/find_all_concerts FindAllConcerts] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	return traceFramework.TraceFunc(ctx, traceFramework.GetTracer("concert.usecase"), func(ctx context.Context) (entity.Page[entity.Concert], error) {
		// Create a new validator instance
		vInstance, err := validator.NewValidator(
			validator.WithTagNameFunc(validator.JSONTagNameFunc),
			validator.WithCustomValidator(new(customvalidator.ThaiTimezoneValidator)),
		)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to create validator", nil))
		}

		// Validate Input
		err = vInstance.Struct(input)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
		}

		return entity.NewPage(u.findAllConcerts(ctx, input))
	})
}

func (u *concertUsecase) findAllConcerts(ctx context.Context, input FindAllConcertsInput) entity.PageProvider[entity.Concert] {
	return func() ([]entity.Concert, entity.PageProvider[entity.Concert], entity.Pagination, error) {
		// Fetch all concerts with optional filters
		concerts, count, err := u.concertRepository.FindAll(ctx, repository.FindAllConcertsFilter{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Venue:     input.Venue,
			Limit:     input.Limit,
			Offset:    input.Offset,
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		})
		if err != nil {
			return entity.Concerts{}, nil, entity.Pagination{}, errsFramework.NewInternalServerError("failed to fetch concerts", nil)
		}
		if concerts == nil || len(pointer.GetValue(concerts)) == 0 {
			return entity.Concerts{}, nil, entity.Pagination{}, nil
		}

		// Create pagination and next search criteria
		pagination := entity.NewPagination(count, pointer.GetValue(input.Limit), pointer.GetValue(input.Offset))
		nextSearchCriteria := FindAllConcertsInput{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Venue:     input.Venue,
			Limit:     input.Limit,
			Offset:    pointer.ToPointer((*input.Limit) + (*input.Offset)),
			SortBy:    input.SortBy,
			SortOrder: input.SortOrder,
		}
		return pointer.GetValue(concerts), u.findAllConcerts(ctx, nextSearchCriteria), pagination, nil
	}
}
