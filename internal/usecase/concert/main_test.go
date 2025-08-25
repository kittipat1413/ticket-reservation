package usecase_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"ticket-reservation/internal/config"
	repository_mocks "ticket-reservation/internal/domain/repository/mocks"
	db_mocks "ticket-reservation/internal/infra/db/mocks"
	concertusecase "ticket-reservation/internal/usecase/concert"
)

type testHelper struct {
	ctrl                  *gomock.Controller
	appConfig             config.AppConfig
	mockTransactorFactory *db_mocks.MockSqlxTransactorFactory
	mockConcertRepository *repository_mocks.MockConcertRepository
	concertUsecase        concertusecase.ConcertUsecase
}

func initTest(t *testing.T) *testHelper {
	ctrl := gomock.NewController(t)

	// Create test app config
	appConfig := config.AppConfig{
		AdminAPIKey:    "test-api-key",
		AdminAPISecret: "test-api-secret",
		Timezone:       "Asia/Bangkok",
		SeatLockTTL:    5 * time.Minute,
	}

	mockTransactorFactory := db_mocks.NewMockSqlxTransactorFactory(ctrl)
	mockConcertRepository := repository_mocks.NewMockConcertRepository(ctrl)

	usecase := concertusecase.NewConcertUsecase(
		appConfig,
		mockTransactorFactory,
		mockConcertRepository,
	)

	return &testHelper{
		ctrl:                  ctrl,
		appConfig:             appConfig,
		mockTransactorFactory: mockTransactorFactory,
		mockConcertRepository: mockConcertRepository,
		concertUsecase:        usecase,
	}
}

func (h *testHelper) Done() {
	h.ctrl.Finish()
}

func TestNewConcertUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appConfig := config.AppConfig{
		AdminAPIKey:    "test-api-key",
		AdminAPISecret: "test-api-secret",
		Timezone:       "Asia/Bangkok",
		SeatLockTTL:    5 * time.Minute,
	}
	mockTransactorFactory := db_mocks.NewMockSqlxTransactorFactory(ctrl)
	mockConcertRepo := repository_mocks.NewMockConcertRepository(ctrl)

	// Execute
	usecase := concertusecase.NewConcertUsecase(appConfig, mockTransactorFactory, mockConcertRepo)

	// Assert
	assert.NotNil(t, usecase)
}
