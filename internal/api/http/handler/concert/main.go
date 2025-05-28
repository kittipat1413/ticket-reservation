package handler

import (
	"ticket-reservation/internal/config"
	concertUsecase "ticket-reservation/internal/usecase/concert"

	"github.com/gin-gonic/gin"
)

type ConcertHandler interface {
	CreateConcert(c *gin.Context)
	FindConcertByID(c *gin.Context)
	FindAllConcerts(c *gin.Context)
}

type concertHandler struct {
	appConfig      config.AppConfig
	concertUsecase concertUsecase.ConcertUsecase
}

func NewConcertHandler(appConfig config.AppConfig, concertUsecase concertUsecase.ConcertUsecase) ConcertHandler {
	return &concertHandler{
		appConfig:      appConfig,
		concertUsecase: concertUsecase,
	}
}
