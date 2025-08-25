package handler_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	handler "ticket-reservation/internal/api/http/handler/healthcheck"
	healthcheck_mocks "ticket-reservation/internal/usecase/healthcheck/mocks"
)

type testHelper struct {
	ctrl                   *gomock.Controller
	mockHealthcheckUsecase *healthcheck_mocks.MockHealthCheckUsecase
	healthcheckHandler     handler.HealthCheckHandler
}

func initTest(t *testing.T) *testHelper {
	ctrl := gomock.NewController(t)

	mockHealthcheckUsecase := healthcheck_mocks.NewMockHealthCheckUsecase(ctrl)

	healthcheckHandler := handler.NewHealthCheckHandler(mockHealthcheckUsecase)

	return &testHelper{
		ctrl:                   ctrl,
		mockHealthcheckUsecase: mockHealthcheckUsecase,
		healthcheckHandler:     healthcheckHandler,
	}
}

func (h *testHelper) Done() {
	h.ctrl.Finish()
}

func TestNewHealthCheckHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHealthcheckUsecase := healthcheck_mocks.NewMockHealthCheckUsecase(ctrl)

	// Execute
	handler := handler.NewHealthCheckHandler(mockHealthcheckUsecase)

	// Assert
	assert.NotNil(t, handler)
}
