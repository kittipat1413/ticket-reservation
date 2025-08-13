package healthcheckrepo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	healthcheckrepo "ticket-reservation/internal/infra/redis/repository/healthcheck"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func TestHealthCheckRepositoryImpl_CheckRedisReadiness(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock redismock.ClientMock)
		expectedOK    bool
		expectedError bool
		errorType     error
	}{
		{
			name: "successful health check - all operations pass",
			setupMock: func(mock redismock.ClientMock) {
				// Mock PING command
				mock.ExpectPing().SetVal("PONG")

				// Mock SET command
				mock.Regexp().ExpectSet(`health_check_.*`, "health_check_value", time.Minute).SetVal("OK")

				// Mock GET command
				mock.Regexp().ExpectGet(`health_check_.*`).SetVal("health_check_value")

				// Mock DEL command
				mock.Regexp().ExpectDel(`health_check_.*`).SetVal(1)
			},
			expectedOK:    true,
			expectedError: false,
		},
		{
			name: "ping command fails",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetErr(errors.New("connection refused"))
			},
			expectedOK:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "ping returns unexpected response",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetVal("INVALID")
			},
			expectedOK:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "set operation fails",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetVal("PONG")
				mock.MatchExpectationsInOrder(false)
				mock.Regexp().ExpectSet(`health_check_.*`, "health_check_value", time.Minute).SetErr(errors.New("set failed"))
			},
			expectedOK:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "get operation fails",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetVal("PONG")
				mock.MatchExpectationsInOrder(false)
				mock.Regexp().ExpectSet(`health_check_.*`, "health_check_value", time.Minute).SetVal("OK")
				mock.Regexp().ExpectGet(`health_check_.*`).SetErr(errors.New("get failed"))

				// Expect cleanup DEL command when GET fails
				mock.Regexp().ExpectDel(`health_check_.*`).SetVal(1)
			},
			expectedOK:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "get returns unexpected value",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetVal("PONG")
				mock.MatchExpectationsInOrder(false)
				mock.Regexp().ExpectSet(`health_check_.*`, "health_check_value", time.Minute).SetVal("OK")
				mock.Regexp().ExpectGet(`health_check_.*`).SetVal("wrong_value")

				// Expect cleanup DEL command when value mismatch
				mock.Regexp().ExpectDel(`health_check_.*`).SetVal(1)
			},
			expectedOK:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "del operation fails",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetVal("PONG")
				mock.MatchExpectationsInOrder(false)
				mock.Regexp().ExpectSet(`health_check_.*`, "health_check_value", time.Minute).SetVal("OK")
				mock.Regexp().ExpectGet(`health_check_.*`).SetVal("health_check_value")
				mock.Regexp().ExpectDel(`health_check_.*`).SetErr(errors.New("del failed"))
			},
			expectedOK:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "del returns unexpected count",
			setupMock: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetVal("PONG")
				mock.MatchExpectationsInOrder(false)
				mock.Regexp().ExpectSet(`health_check_.*`, "health_check_value", time.Minute).SetVal("OK")
				mock.Regexp().ExpectGet(`health_check_.*`).SetVal("health_check_value")
				mock.Regexp().ExpectDel(`health_check_.*`).SetVal(0) // Should return 1 for successful deletion
			},
			expectedOK:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			repository := healthcheckrepo.NewHealthCheckRepository(db)

			tt.setupMock(mock)

			// Execute
			ok, err := repository.CheckRedisReadiness(context.Background())

			// Assert
			assert.Equal(t, tt.expectedOK, ok)
			if tt.expectedError {
				require.Error(t, err)
				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository healthcheck/check_redis_readiness CheckRedisReadiness]")
				// Verify it's the expected error type
				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
