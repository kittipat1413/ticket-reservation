package healthcheckrepo_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func TestHealthCheckRepositoryImpl_CheckDatabaseReadiness(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedOk    bool
		expectedError bool
		errorType     error
	}{
		{
			name: "successful database check",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"?column?"}).AddRow(true)
				mock.ExpectQuery("SELECT 1=1").WillReturnRows(rows)
			},
			expectedOk:    true,
			expectedError: false,
		},
		{
			name: "database connection error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1=1").WillReturnError(sql.ErrConnDone)
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "database timeout error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1=1").WillReturnError(context.DeadlineExceeded)
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "generic database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1=1").WillReturnError(errors.New("database is down"))
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name: "no rows returned error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1=1").WillReturnError(sql.ErrNoRows)
			},
			expectedOk:    false,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			tt.setupMock(h.Mock)
			ok, err := h.Repository.CheckDatabaseReadiness(context.Background())

			// Assert
			assert.Equal(t, tt.expectedOk, ok)
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository healthcheck/check_database_readiness CheckDatabaseReadiness]")

				// Verify it's the expected error type
				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			h.AssertExpectationsMet(t)
		})
	}
}

func TestHealthCheckRepositoryImpl_CheckDatabaseReadiness_QueryValidation(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	// Setup expectations - verify exact query
	rows := sqlmock.NewRows([]string{"?column?"}).AddRow(true)
	h.Mock.ExpectQuery("SELECT 1=1").WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	ok, err := h.Repository.CheckDatabaseReadiness(ctx)

	// Assert
	require.NoError(t, err)
	assert.True(t, ok)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
