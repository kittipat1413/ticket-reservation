package seat

import (
	"context"
	"ticket-reservation/internal/config"
	"ticket-reservation/internal/domain/cache"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/domain/repository"
	"ticket-reservation/internal/infra/db"
)

//go:generate mockgen -source=./main.go -destination=./mocks/seat_usecase.go -package=seat_usecasemocks
type SeatUsecase interface {
	ReserveSeat(ctx context.Context, input ReserveSeatInput) (*entity.Reservation, error)
}

type seatUsecase struct {
	appConfig             config.AppConfig
	concertRepository     repository.ConcertRepository
	zoneRepository        repository.ZoneRepository
	seatRepository        repository.SeatRepository
	reservationRepository repository.ReservationRepository
	transactorFactory     db.SqlxTransactorFactory
	seatLocker            cache.SeatLocker
}

func NewSeatUsecase(
	appConfig config.AppConfig,
	concertRepository repository.ConcertRepository,
	zoneRepository repository.ZoneRepository,
	seatRepository repository.SeatRepository,
	reservationRepository repository.ReservationRepository,
	transactorFactory db.SqlxTransactorFactory,
	seatLocker cache.SeatLocker,
) SeatUsecase {
	return &seatUsecase{
		appConfig:             appConfig,
		concertRepository:     concertRepository,
		zoneRepository:        zoneRepository,
		seatRepository:        seatRepository,
		reservationRepository: reservationRepository,
		transactorFactory:     transactorFactory,
		seatLocker:            seatLocker,
	}
}
