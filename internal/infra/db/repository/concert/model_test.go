package concertrepo_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"
	concertrepo "ticket-reservation/internal/infra/db/repository/concert"
)

func TestConcert_ToEntity(t *testing.T) {
	testID := uuid.New()
	testName := "Test Concert"
	testVenue := "Test Venue"
	testDate := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		concert  concertrepo.Concert
		expected *entity.Concert
	}{
		{
			name: "successful conversion with all fields",
			concert: concertrepo.Concert{
				Concerts: model.Concerts{
					ID:        testID,
					Name:      testName,
					Venue:     testVenue,
					Date:      testDate,
					CreatedAt: testCreatedAt,
					UpdatedAt: testUpdatedAt,
				},
			},
			expected: &entity.Concert{
				ID:        testID,
				Name:      testName,
				Venue:     testVenue,
				Date:      testDate,
				CreatedAt: testCreatedAt,
				UpdatedAt: testUpdatedAt,
			},
		},
		{
			name: "conversion with empty strings",
			concert: concertrepo.Concert{
				Concerts: model.Concerts{
					ID:        testID,
					Name:      "",
					Venue:     "",
					Date:      testDate,
					CreatedAt: testCreatedAt,
					UpdatedAt: testUpdatedAt,
				},
			},
			expected: &entity.Concert{
				ID:        testID,
				Name:      "",
				Venue:     "",
				Date:      testDate,
				CreatedAt: testCreatedAt,
				UpdatedAt: testUpdatedAt,
			},
		},
		{
			name: "conversion with zero time values",
			concert: concertrepo.Concert{
				Concerts: model.Concerts{
					ID:        testID,
					Name:      testName,
					Venue:     testVenue,
					Date:      time.Time{},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			expected: &entity.Concert{
				ID:        testID,
				Name:      testName,
				Venue:     testVenue,
				Date:      time.Time{},
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.concert.ToEntity()

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Venue, result.Venue)
			assert.Equal(t, tt.expected.Date, result.Date)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}

func TestConcerts_ToEntities(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()
	testDate1 := time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC)
	testDate2 := time.Date(2025, 12, 26, 21, 0, 0, 0, time.UTC)
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		concerts concertrepo.Concerts
		expected *entity.Concerts
	}{
		{
			name:     "empty slice",
			concerts: concertrepo.Concerts{},
			expected: &entity.Concerts{},
		},
		{
			name: "single concert",
			concerts: concertrepo.Concerts{
				{
					Concerts: model.Concerts{
						ID:        testID1,
						Name:      "Concert 1",
						Venue:     "Venue 1",
						Date:      testDate1,
						CreatedAt: testCreatedAt,
						UpdatedAt: testUpdatedAt,
					},
				},
			},
			expected: &entity.Concerts{
				{
					ID:        testID1,
					Name:      "Concert 1",
					Venue:     "Venue 1",
					Date:      testDate1,
					CreatedAt: testCreatedAt,
					UpdatedAt: testUpdatedAt,
				},
			},
		},
		{
			name: "multiple concerts",
			concerts: concertrepo.Concerts{
				{
					Concerts: model.Concerts{
						ID:        testID1,
						Name:      "Concert 1",
						Venue:     "Venue 1",
						Date:      testDate1,
						CreatedAt: testCreatedAt,
						UpdatedAt: testUpdatedAt,
					},
				},
				{
					Concerts: model.Concerts{
						ID:        testID2,
						Name:      "Concert 2",
						Venue:     "Venue 2",
						Date:      testDate2,
						CreatedAt: testCreatedAt,
						UpdatedAt: testUpdatedAt,
					},
				},
			},
			expected: &entity.Concerts{
				{
					ID:        testID1,
					Name:      "Concert 1",
					Venue:     "Venue 1",
					Date:      testDate1,
					CreatedAt: testCreatedAt,
					UpdatedAt: testUpdatedAt,
				},
				{
					ID:        testID2,
					Name:      "Concert 2",
					Venue:     "Venue 2",
					Date:      testDate2,
					CreatedAt: testCreatedAt,
					UpdatedAt: testUpdatedAt,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.concerts.ToEntities()

			require.NotNil(t, result)
			assert.Equal(t, len(*tt.expected), len(*result))

			for i, expectedConcert := range *tt.expected {
				actualConcert := (*result)[i]
				assert.Equal(t, expectedConcert.ID, actualConcert.ID)
				assert.Equal(t, expectedConcert.Name, actualConcert.Name)
				assert.Equal(t, expectedConcert.Venue, actualConcert.Venue)
				assert.Equal(t, expectedConcert.Date, actualConcert.Date)
				assert.Equal(t, expectedConcert.CreatedAt, actualConcert.CreatedAt)
				assert.Equal(t, expectedConcert.UpdatedAt, actualConcert.UpdatedAt)
			}
		})
	}
}
