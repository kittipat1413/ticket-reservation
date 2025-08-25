package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func TestHealthCheckUsecase_CheckReadiness(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(h *testHelper)
		expectedOk    bool
		expectedError bool
		errorType     error
		errorContains string
	}{
		{
			name: "successful readiness check - both DB and Redis ready",
			setupMocks: func(h *testHelper) {
				// DB check succeeds immediately
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(true, nil)

				// Redis check succeeds immediately
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(true, nil)
			},
			expectedOk:    true,
			expectedError: false,
		},
		{
			name: "DB readiness fails - non-retryable error",
			setupMocks: func(h *testHelper) {
				// DB check fails with internal server error (non-retryable)
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(false, errsFramework.NewInternalServerError("internal error", nil))
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.InternalServerError{},
			errorContains: "service is not ready",
		},
		{
			name: "DB readiness fails - retryable database error eventually succeeds",
			setupMocks: func(h *testHelper) {
				// First attempt: database error (retryable)
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(false, errsFramework.NewDatabaseError("connection timeout", "timeout"))

				// Second attempt: succeeds
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(true, nil)

				// Redis check succeeds
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(true, nil)
			},
			expectedOk:    true,
			expectedError: false,
		},
		{
			name: "DB readiness fails - retryable database error exhausts retries",
			setupMocks: func(h *testHelper) {
				// All retry attempts fail with database error
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(false, errsFramework.NewDatabaseError("connection failed", "error")).
					Times(3) // MaxAttempts = 3
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.InternalServerError{},
			errorContains: "service is not ready",
		},
		{
			name: "Redis readiness fails - non-retryable error",
			setupMocks: func(h *testHelper) {
				// DB check succeeds
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(true, nil)

				// Redis check fails with internal server error (non-retryable)
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(false, errsFramework.NewInternalServerError("internal error", nil))
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.InternalServerError{},
			errorContains: "service is not ready",
		},
		{
			name: "Redis readiness fails - retryable database error eventually succeeds",
			setupMocks: func(h *testHelper) {
				// DB check succeeds
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(true, nil)

				// First attempt: database error (retryable)
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(false, errsFramework.NewDatabaseError("redis connection timeout", "timeout"))

				// Second attempt: succeeds
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(true, nil)
			},
			expectedOk:    true,
			expectedError: false,
		},
		{
			name: "Redis readiness fails - retryable database error exhausts retries",
			setupMocks: func(h *testHelper) {
				// DB check succeeds
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(true, nil)

				// All retry attempts fail with database error
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(false, errsFramework.NewDatabaseError("redis connection failed", "error")).
					Times(3) // MaxAttempts = 3
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.InternalServerError{},
			errorContains: "service is not ready",
		},
		{
			name: "DB returns false but no error - service not ready",
			setupMocks: func(h *testHelper) {
				// DB check returns false (not ready) but no error
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(false, nil)

				// Redis check succeeds
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(true, nil)
			},
			expectedOk:    false,
			expectedError: false,
		},
		{
			name: "Redis returns false but no error - service not ready",
			setupMocks: func(h *testHelper) {
				// DB check succeeds
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(true, nil)

				// Redis check returns false (not ready) but no error
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(false, nil)
			},
			expectedOk:    false,
			expectedError: false,
		},
		{
			name: "both DB and Redis return false - service not ready",
			setupMocks: func(h *testHelper) {
				// DB check returns false (not ready) but no error
				h.mockDBRepository.EXPECT().
					CheckDatabaseReadiness(gomock.Any()).
					Return(false, nil)

				// Redis check returns false (not ready) but no error
				h.mockRedisRepository.EXPECT().
					CheckRedisReadiness(gomock.Any()).
					Return(false, nil)
			},
			expectedOk:    false,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			// Setup mocks
			tt.setupMocks(h)

			// Execute
			ctx := context.Background()
			ok, err := h.healthCheckUsecase.CheckReadiness(ctx)

			// Assert
			assert.Equal(t, tt.expectedOk, ok, "Expected ok to be %v", tt.expectedOk)

			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[usecase healthcheck/check_readiness CheckReadiness]")

				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}

				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
