package handler

import (
	"ticket-reservation/internal/config"
	seatUsecase "ticket-reservation/internal/usecase/seat"

	"github.com/gin-gonic/gin"
)

type SeatHandler interface {
	ReserveSeat(c *gin.Context)
}

type seatHandler struct {
	appConfig   config.AppConfig
	seatUsecase seatUsecase.SeatUsecase
}

func NewSeatHandler(appConfig config.AppConfig, seatUsecase seatUsecase.SeatUsecase) SeatHandler {
	return &seatHandler{
		appConfig:   appConfig,
		seatUsecase: seatUsecase,
	}
}
