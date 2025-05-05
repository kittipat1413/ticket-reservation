package concert

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	customvalidator "ticket-reservation/pkg/validator"
	"time"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	traceFramework "github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/framework/validator"
)

type CreateConcertInput struct {
	Name  string    `json:"name" validate:"required,gt=0"`
	Venue string    `json:"venue" validate:"required,gt=0"`
	Date  time.Time `json:"date" validate:"required,thaitimezone"`
}

func (u *concertUsecase) CreateConcert(ctx context.Context, input CreateConcertInput) (concert *entity.Concert, err error) {
	const errLocation = "[usecase concert/create_concert CreateConcert] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	return traceFramework.TraceFunc(ctx, traceFramework.GetTracer("concert.usecase"), func(ctx context.Context) (*entity.Concert, error) {
		// Create a new validator instance
		vInstance, err := validator.NewValidator(
			validator.WithTagNameFunc(validator.JSONTagNameFunc),
			validator.WithCustomValidator(new(customvalidator.ThaiTimezoneValidator)),
		)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to create validator", nil))
		}

		// Validate input
		err = vInstance.Struct(input)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewBadRequestError("invalid input", map[string]string{"details": err.Error()}))
		}

		concert := &entity.Concert{
			Name:  input.Name,
			Venue: input.Venue,
			Date:  input.Date,
		}

		created, err := u.concertRepository.CreateOne(ctx, concert)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to create concert", nil))
		}

		return created, nil
	})
}
