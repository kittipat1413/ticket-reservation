package concertrepo_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/domain/repository"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	"github.com/kittipat1413/go-common/util/pointer"
)

func TestConcertRepositoryImpl_FindAll(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()
	testDate1 := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	testDate2 := time.Date(2025, 12, 26, 21, 0, 0, 0, time.UTC)
	createdAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		filter           repository.FindAllConcertsFilter
		setupMock        func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter)
		expectedConcerts *entity.Concerts
		expectedTotal    int64
		expectedError    bool
		errorType        error
	}{
		{
			name:   "successful retrieval with no filter",
			filter: repository.FindAllConcertsFilter{},
			setupMock: func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter) {
				// Count query
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(2)
				mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts`).
					WillReturnRows(countRows)

				// Main query
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.date",
					"concerts.venue", "concerts.created_at", "concerts.updated_at",
				}).
					AddRow(testID1, "Concert 1", testDate1, "Venue 1", createdAt, updatedAt).
					AddRow(testID2, "Concert 2", testDate2, "Venue 2", createdAt, updatedAt)

				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts`).
					WillReturnRows(rows)
			},
			expectedConcerts: &entity.Concerts{
				{ID: testID1, Name: "Concert 1", Venue: "Venue 1", Date: testDate1, CreatedAt: createdAt, UpdatedAt: updatedAt},
				{ID: testID2, Name: "Concert 2", Venue: "Venue 2", Date: testDate2, CreatedAt: createdAt, UpdatedAt: updatedAt},
			},
			expectedTotal: 2,
			expectedError: false,
		},
		{
			name: "successful retrieval with venue filter",
			filter: repository.FindAllConcertsFilter{
				Venue: pointer.ToPointer("Test Venue"),
			},
			setupMock: func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter) {
				// Count query with WHERE clause
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(1)
				mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts WHERE \(concerts\.venue LIKE \$1::text\)`).
					WithArgs("%Test Venue%").
					WillReturnRows(countRows)

				// Main query with WHERE clause
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.date",
					"concerts.venue", "concerts.created_at", "concerts.updated_at",
				}).AddRow(testID1, "Concert 1", testDate1, "Test Venue", createdAt, updatedAt)

				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE \(concerts\.venue LIKE \$1::text\)`).
					WithArgs("%Test Venue%").
					WillReturnRows(rows)
			},
			expectedConcerts: &entity.Concerts{
				{ID: testID1, Name: "Concert 1", Venue: "Test Venue", Date: testDate1, CreatedAt: createdAt, UpdatedAt: updatedAt},
			},
			expectedTotal: 1,
			expectedError: false,
		},
		{
			name: "successful retrieval with date range filter",
			filter: repository.FindAllConcertsFilter{
				StartDate: pointer.ToPointer(time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)),
				EndDate:   pointer.ToPointer(time.Date(2025, 12, 26, 23, 59, 59, 0, time.UTC)),
			},
			setupMock: func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter) {
				// Count query with date range
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(2)
				mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts WHERE \( \(concerts\.date >= \$1::timestamp with time zone\) AND \(concerts\.date <= \$2::timestamp with time zone\) \)`).
					WithArgs(*filter.StartDate, *filter.EndDate).
					WillReturnRows(countRows)

				// Main query with date range
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.date",
					"concerts.venue", "concerts.created_at", "concerts.updated_at",
				}).
					AddRow(testID1, "Concert 1", testDate1, "Venue 1", createdAt, updatedAt).
					AddRow(testID2, "Concert 2", testDate2, "Venue 2", createdAt, updatedAt)

				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts WHERE \( \(concerts\.date >= \$1::timestamp with time zone\) AND \(concerts\.date <= \$2::timestamp with time zone\) \)`).
					WithArgs(*filter.StartDate, *filter.EndDate).
					WillReturnRows(rows)
			},
			expectedConcerts: &entity.Concerts{
				{ID: testID1, Name: "Concert 1", Venue: "Venue 1", Date: testDate1, CreatedAt: createdAt, UpdatedAt: updatedAt},
				{ID: testID2, Name: "Concert 2", Venue: "Venue 2", Date: testDate2, CreatedAt: createdAt, UpdatedAt: updatedAt},
			},
			expectedTotal: 2,
			expectedError: false,
		},
		{
			name: "successful retrieval with pagination and sorting",
			filter: repository.FindAllConcertsFilter{
				Limit:     pointer.ToPointer(int64(10)),
				Offset:    pointer.ToPointer(int64(0)),
				SortBy:    pointer.ToPointer("name"),
				SortOrder: pointer.ToPointer(entity.SortOrderAsc),
			},
			setupMock: func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter) {
				// Count query
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(2)
				mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts`).
					WillReturnRows(countRows)

				// Main query with ORDER BY, LIMIT, OFFSET
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.date",
					"concerts.venue", "concerts.created_at", "concerts.updated_at",
				}).
					AddRow(testID1, "Concert 1", testDate1, "Venue 1", createdAt, updatedAt).
					AddRow(testID2, "Concert 2", testDate2, "Venue 2", createdAt, updatedAt)

				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts ORDER BY concerts\.name ASC LIMIT \$1 OFFSET \$2`).
					WithArgs(*filter.Limit, *filter.Offset).
					WillReturnRows(rows)
			},
			expectedConcerts: &entity.Concerts{
				{ID: testID1, Name: "Concert 1", Venue: "Venue 1", Date: testDate1, CreatedAt: createdAt, UpdatedAt: updatedAt},
				{ID: testID2, Name: "Concert 2", Venue: "Venue 2", Date: testDate2, CreatedAt: createdAt, UpdatedAt: updatedAt},
			},
			expectedTotal: 2,
			expectedError: false,
		},
		{
			name:   "count query database error",
			filter: repository.FindAllConcertsFilter{},
			setupMock: func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter) {
				// Count query fails
				mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedConcerts: nil,
			expectedTotal:    0,
			expectedError:    true,
			errorType:        &errsFramework.DatabaseError{},
		},
		{
			name:   "main query database error",
			filter: repository.FindAllConcertsFilter{},
			setupMock: func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter) {
				// Count query succeeds
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(2)
				mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts`).
					WillReturnRows(countRows)

				// Main query fails
				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts`).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedConcerts: nil,
			expectedTotal:    0,
			expectedError:    true,
			errorType:        &errsFramework.DatabaseError{},
		},
		{
			name:   "no concerts found",
			filter: repository.FindAllConcertsFilter{},
			setupMock: func(mock sqlmock.Sqlmock, filter repository.FindAllConcertsFilter) {
				// Count query returns 0
				countRows := sqlmock.NewRows([]string{"total"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts`).
					WillReturnRows(countRows)

				// Main query returns empty result
				rows := sqlmock.NewRows([]string{
					"concerts.id", "concerts.name", "concerts.date",
					"concerts.venue", "concerts.created_at", "concerts.updated_at",
				})

				mock.ExpectQuery(`SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts`).
					WillReturnRows(rows)
			},
			expectedConcerts: &entity.Concerts{},
			expectedTotal:    0,
			expectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.done()

			tt.setupMock(h.mock, tt.filter)
			concerts, total, err := h.repository.FindAll(context.Background(), tt.filter)

			// Assert
			if tt.expectedError {
				require.Error(t, err)

				// Verify it's wrapped with the expected error prefix
				assert.Contains(t, err.Error(), "[repository concert/find_all FindAll]")

				// Verify it's the expected error type
				if tt.errorType != nil {
					assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
				}

				assert.Nil(t, concerts)
				assert.Equal(t, int64(0), total)
			} else {
				require.NoError(t, err)
				require.NotNil(t, concerts)
				assert.Equal(t, tt.expectedTotal, total)

				// Compare concerts
				assert.Equal(t, len(*tt.expectedConcerts), len(*concerts))
				for i, expectedConcert := range *tt.expectedConcerts {
					actualConcert := (*concerts)[i]
					assert.Equal(t, expectedConcert.ID, actualConcert.ID)
					assert.Equal(t, expectedConcert.Name, actualConcert.Name)
					assert.Equal(t, expectedConcert.Venue, actualConcert.Venue)
					assert.Equal(t, expectedConcert.Date.UTC(), actualConcert.Date.UTC())
					assert.Equal(t, expectedConcert.CreatedAt.UTC(), actualConcert.CreatedAt.UTC())
					assert.Equal(t, expectedConcert.UpdatedAt.UTC(), actualConcert.UpdatedAt.UTC())
				}
			}

			// Verify all expectations were met
			err = h.mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestConcertRepositoryImpl_FindAll_SortingOptions(t *testing.T) {
	testID := uuid.New()
	testDate := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	createdAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name            string
		sortBy          string
		sortOrder       entity.SortOrder
		expectedOrderBy string
	}{
		{
			name:            "sort by name ascending",
			sortBy:          "name",
			sortOrder:       entity.SortOrderAsc,
			expectedOrderBy: "ORDER BY concerts\\.name ASC",
		},
		{
			name:            "sort by name descending",
			sortBy:          "name",
			sortOrder:       entity.SortOrderDesc,
			expectedOrderBy: "ORDER BY concerts\\.name DESC",
		},
		{
			name:            "sort by venue ascending",
			sortBy:          "venue",
			sortOrder:       entity.SortOrderAsc,
			expectedOrderBy: "ORDER BY concerts\\.venue ASC",
		},
		{
			name:            "sort by venue descending",
			sortBy:          "venue",
			sortOrder:       entity.SortOrderDesc,
			expectedOrderBy: "ORDER BY concerts\\.venue DESC",
		},
		{
			name:            "sort by date ascending",
			sortBy:          "date",
			sortOrder:       entity.SortOrderAsc,
			expectedOrderBy: "ORDER BY concerts\\.date ASC",
		},
		{
			name:            "sort by date descending",
			sortBy:          "date",
			sortOrder:       entity.SortOrderDesc,
			expectedOrderBy: "ORDER BY concerts\\.date DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := initTest(t)
			defer h.done()

			filter := repository.FindAllConcertsFilter{
				SortBy:    &tt.sortBy,
				SortOrder: &tt.sortOrder,
			}

			// Count query
			countRows := sqlmock.NewRows([]string{"total"}).AddRow(1)
			h.mock.ExpectQuery(`SELECT COUNT\(concerts\.id\) AS "total" FROM public\.concerts`).
				WillReturnRows(countRows)

			// Main query with specific ORDER BY
			rows := sqlmock.NewRows([]string{
				"concerts.id", "concerts.name", "concerts.date",
				"concerts.venue", "concerts.created_at", "concerts.updated_at",
			}).AddRow(testID, "Test Concert", testDate, "Test Venue", createdAt, updatedAt)

			expectedQuery := `SELECT concerts\.id AS "concerts\.id", concerts\.name AS "concerts\.name", concerts\.date AS "concerts\.date", concerts\.venue AS "concerts\.venue", concerts\.created_at AS "concerts\.created_at", concerts\.updated_at AS "concerts\.updated_at" FROM public\.concerts ` + tt.expectedOrderBy
			h.mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

			_, _, err := h.repository.FindAll(context.Background(), filter)

			require.NoError(t, err)

			// Verify all expectations were met
			err = h.mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
