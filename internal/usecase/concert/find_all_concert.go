package concert

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	customvalidator "ticket-reservation/pkg/validator"
	"time"

	repository "ticket-reservation/internal/domain/repository"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	traceFramework "github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/framework/validator"
)

type FindAllConcertsInput struct {
	StartDate *time.Time `json:"start_date" validate:"omitempty,thaitimezone"`
	EndDate   *time.Time `json:"end_date" validate:"omitempty,thaitimezone,gtfield=StartDate"`
	Venue     *string    `json:"venue" validate:"omitempty,gt=0"`
}

func (u *concertUsecase) FindAllConcerts(ctx context.Context, input FindAllConcertsInput) (concerts *entity.Concerts, err error) {
	const errLocation = "[usecase concert/find_all_concerts FindAllConcerts] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	return traceFramework.TraceFunc(ctx, traceFramework.GetTracer("concert.usecase"), func(ctx context.Context) (*entity.Concerts, error) {
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

		// Fetch all concerts with optional filters
		concerts, err := u.concertRepository.FindAll(ctx, repository.FindAllConcertsFilter{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Venue:     input.Venue,
		})
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to fetch concerts", nil))
		}

		return concerts, nil
	})
}
