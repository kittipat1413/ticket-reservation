package concert

import (
	"context"
	"ticket-reservation/internal/config"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/domain/repository"
	"ticket-reservation/internal/infra/db"
)

//go:generate mockgen -source=./main.go -destination=./mocks/concert_usecase.go -package=concert_usecasemocks
type ConcertUsecase interface {
	CreateConcert(ctx context.Context, concert CreateConcertInput) (*entity.Concert, error)
	// GetConcertByID(ctx context.Context, id string) (*entity.Concert, error)
}

type concertUsecase struct {
	appConfig         config.AppConfig
	transactorFactory db.SqlxTransactorFactory
	concertRepository repository.ConcertRepository
}

func NewConcertUsecase(
	appConfig config.AppConfig,
	transactorFactory db.SqlxTransactorFactory,
	concertRepository repository.ConcertRepository,
) ConcertUsecase {
	return &concertUsecase{
		appConfig:         appConfig,
		transactorFactory: transactorFactory,
		concertRepository: concertRepository,
	}
}
