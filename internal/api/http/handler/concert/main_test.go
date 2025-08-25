package handler_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	handler "ticket-reservation/internal/api/http/handler/concert"
	"ticket-reservation/internal/config"
	concert_mocks "ticket-reservation/internal/usecase/concert/mocks"
)

type testHelper struct {
	ctrl               *gomock.Controller
	appConfig          config.AppConfig
	mockConcertUsecase *concert_mocks.MockConcertUsecase
	concertHandler     handler.ConcertHandler
}

func initTest(t *testing.T) *testHelper {
	ctrl := gomock.NewController(t)

	appConfig := config.AppConfig{
		AdminAPIKey:    "test-api-key",
		AdminAPISecret: "test-api-secret",
		Timezone:       "Asia/Bangkok",
		SeatLockTTL:    5 * time.Minute,
	}

	mockConcertUsecase := concert_mocks.NewMockConcertUsecase(ctrl)

	concertHandler := handler.NewConcertHandler(appConfig, mockConcertUsecase)

	return &testHelper{
		ctrl:               ctrl,
		appConfig:          appConfig,
		mockConcertUsecase: mockConcertUsecase,
		concertHandler:     concertHandler,
	}
}

func (h *testHelper) Done() {
	h.ctrl.Finish()
}

func TestNewConcertHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appConfig := config.AppConfig{
		AdminAPIKey:    "test-api-key",
		AdminAPISecret: "test-api-secret",
		Timezone:       "Asia/Bangkok",
		SeatLockTTL:    5 * time.Minute,
	}
	mockConcertUsecase := concert_mocks.NewMockConcertUsecase(ctrl)

	// Execute
	handler := handler.NewConcertHandler(appConfig, mockConcertUsecase)

	// Assert
	assert.NotNil(t, handler)
}
