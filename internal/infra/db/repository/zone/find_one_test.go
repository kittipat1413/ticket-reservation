package zonerepo_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func TestZoneRepositoryImpl_FindOne(t *testing.T) {
	testID := uuid.New()
	testConcertID := uuid.New()
	testDescription := "VIP Section"
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		zoneID        uuid.UUID
		setupMock     func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectedZone  *entity.Zone
		expectedError bool
		errorType     error
	}{
		{
			name:   "successful zone retrieval",
			zoneID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"zones.id", "zones.concert_id", "zones.name", "zones.description",
					"zones.created_at", "zones.updated_at",
				}).AddRow(
					id, testConcertID, "VIP", &testDescription, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`SELECT zones\.id AS "zones\.id", zones\.concert_id AS "zones\.concert_id", zones\.name AS "zones\.name", zones\.description AS "zones\.description", zones\.created_at AS "zones\.created_at", zones\.updated_at AS "zones\.updated_at" FROM public\.zones WHERE zones\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedZone: &entity.Zone{
				ID:          testID,
				ConcertID:   testConcertID,
				Name:        "VIP",
				Description: &testDescription,
				CreatedAt:   testCreatedAt,
				UpdatedAt:   testUpdatedAt,
			},
			expectedError: false,
		},
		{
			name:   "zone not found",
			zoneID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT zones\.id AS "zones\.id", zones\.concert_id AS "zones\.concert_id", zones\.name AS "zones\.name", zones\.description AS "zones\.description", zones\.created_at AS "zones\.created_at", zones\.updated_at AS "zones\.updated_at" FROM public\.zones WHERE zones\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedZone:  nil,
			expectedError: true,
			errorType:     &errsFramework.NotFoundError{},
		},
		{
			name:   "database error",
			zoneID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT zones\.id AS "zones\.id", zones\.concert_id AS "zones\.concert_id", zones\.name AS "zones\.name", zones\.description AS "zones\.description", zones\.created_at AS "zones\.created_at", zones\.updated_at AS "zones\.updated_at" FROM public\.zones WHERE zones\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnError(sql.ErrConnDone)
			},
			expectedZone:  nil,
			expectedError: true,
			errorType:     &errsFramework.DatabaseError{},
		},
		{
			name:   "nil description",
			zoneID: testID,
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"zones.id", "zones.concert_id", "zones.name", "zones.description",
					"zones.created_at", "zones.updated_at",
				}).AddRow(
					id, testConcertID, "General", nil, testCreatedAt, testUpdatedAt,
				)

				mock.ExpectQuery(`SELECT zones\.id AS "zones\.id", zones\.concert_id AS "zones\.concert_id", zones\.name AS "zones\.name", zones\.description AS "zones\.description", zones\.created_at AS "zones\.created_at", zones\.updated_at AS "zones\.updated_at" FROM public\.zones WHERE zones\.id = \$1 FOR UPDATE`).
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedZone: &entity.Zone{
				ID:          testID,
				ConcertID:   testConcertID,
				Name:        "General",
				Description: nil,
				CreatedAt:   testCreatedAt,
				UpdatedAt:   testUpdatedAt,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.Done()

			ctx := context.Background()
			tt.setupMock(h.Mock, tt.zoneID)

			// Execute
			result, err := h.Repository.FindOne(ctx, tt.zoneID)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Nil(t, result)
				assert.ErrorAs(t, err, &tt.errorType)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedZone.ID, result.ID)
				assert.Equal(t, tt.expectedZone.ConcertID, result.ConcertID)
				assert.Equal(t, tt.expectedZone.Name, result.Name)
				assert.Equal(t, tt.expectedZone.Description, result.Description)
				assert.Equal(t, tt.expectedZone.CreatedAt, result.CreatedAt)
				assert.Equal(t, tt.expectedZone.UpdatedAt, result.UpdatedAt)
			}

			// Verify all expectations were met
			h.AssertExpectationsMet(t)
		})
	}
}

func TestZoneRepositoryImpl_FindOne_QueryValidation(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	testID := uuid.New()

	// Setup expectations - verify exact query
	rows := sqlmock.NewRows([]string{
		"zones.id", "zones.concert_id", "zones.name", "zones.description",
		"zones.created_at", "zones.updated_at",
	}).AddRow(
		testID, uuid.New(), "VIP", "Front Row", time.Now(), time.Now(),
	)

	h.Mock.ExpectQuery(`SELECT zones\.id AS "zones\.id", zones\.concert_id AS "zones\.concert_id", zones\.name AS "zones\.name", zones\.description AS "zones\.description", zones\.created_at AS "zones\.created_at", zones\.updated_at AS "zones\.updated_at" FROM public\.zones WHERE zones\.id = \$1 FOR UPDATE`).
		WithArgs(testID).
		WillReturnRows(rows)

	ctx := context.Background()

	// Execute
	result, err := h.Repository.FindOne(ctx, testID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, testID, result.ID)

	// Verify all expectations were met
	h.AssertExpectationsMet(t)
}
