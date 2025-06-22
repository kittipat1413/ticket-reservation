package usecase

import (
	"context"
	"errors"
	"ticket-reservation/internal/domain/entity"

	"github.com/google/uuid"
	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	traceFramework "github.com/kittipat1413/go-common/framework/trace"
	"github.com/kittipat1413/go-common/framework/validator"
)

type FindOneConcertInput struct {
	ID string `json:"id" validate:"required,uuid4"`
}

func (u *concertUsecase) FindOneConcert(ctx context.Context, input FindOneConcertInput) (concert *entity.Concert, err error) {
	const errLocation = "[usecase concert/find_one_concert FindOneConcert] "
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	return traceFramework.TraceFunc(ctx, traceFramework.GetTracer("concert.usecase"), func(ctx context.Context) (*entity.Concert, error) {
		// Create a new validator instance
		vInstance, err := validator.NewValidator(
			validator.WithTagNameFunc(validator.JSONTagNameFunc),
		)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to create validator", nil))
		}

		// Validate Input
		err = vInstance.Struct(input)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewBadRequestError("the request is invalid", map[string]string{"details": err.Error()}))
		}

		concertID, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, errsFramework.WrapError(err, errsFramework.NewBadRequestError("invalid concert ID", nil))
		}

		// Find concert by ID
		concert, err := u.concertRepository.FindOne(ctx, concertID)
		if err != nil {
			var notFoundErr *errsFramework.NotFoundError
			if !errors.As(err, &notFoundErr) { // If the error is not a NotFoundError, wrap it as an internal server error
				return nil, errsFramework.WrapError(err, errsFramework.NewInternalServerError("failed to find concert by ID", nil))
			}
			return nil, err // Return the NotFoundError directly
		}
		return concert, nil
	})
}
