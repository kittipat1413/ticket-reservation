package zonerepo_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db/model_gen/ticket-reservation/public/model"
	zonerepo "ticket-reservation/internal/infra/db/repository/zone"

	"github.com/kittipat1413/go-common/util/pointer"
)

func TestZone_ToEntity(t *testing.T) {
	testID := uuid.New()
	testConcertID := uuid.New()
	testName := "VIP Section"
	testDescription := "Premium seating area"
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          zonerepo.Zone
		expectedEntity *entity.Zone
		expectedNil    bool
	}{
		{
			name: "successful conversion with description",
			input: zonerepo.Zone{
				Zones: model.Zones{
					ID:          testID,
					ConcertID:   testConcertID,
					Name:        testName,
					Description: &testDescription,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedEntity: &entity.Zone{
				ID:          testID,
				ConcertID:   testConcertID,
				Name:        testName,
				Description: &testDescription,
				CreatedAt:   testCreatedAt,
				UpdatedAt:   testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "successful conversion with nil description",
			input: zonerepo.Zone{
				Zones: model.Zones{
					ID:          testID,
					ConcertID:   testConcertID,
					Name:        "General Admission",
					Description: nil,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedEntity: &entity.Zone{
				ID:          testID,
				ConcertID:   testConcertID,
				Name:        "General Admission",
				Description: nil,
				CreatedAt:   testCreatedAt,
				UpdatedAt:   testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "conversion with empty name",
			input: zonerepo.Zone{
				Zones: model.Zones{
					ID:          testID,
					ConcertID:   testConcertID,
					Name:        "",
					Description: &testDescription,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedEntity: &entity.Zone{
				ID:          testID,
				ConcertID:   testConcertID,
				Name:        "",
				Description: &testDescription,
				CreatedAt:   testCreatedAt,
				UpdatedAt:   testUpdatedAt,
			},
			expectedNil: false,
		},
		{
			name: "conversion with zero time values",
			input: zonerepo.Zone{
				Zones: model.Zones{
					ID:          testID,
					ConcertID:   testConcertID,
					Name:        testName,
					Description: &testDescription,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			expectedEntity: &entity.Zone{
				ID:          testID,
				ConcertID:   testConcertID,
				Name:        testName,
				Description: &testDescription,
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			expectedNil: false,
		},
		{
			name: "conversion with empty description string",
			input: zonerepo.Zone{
				Zones: model.Zones{
					ID:          testID,
					ConcertID:   testConcertID,
					Name:        testName,
					Description: pointer.ToPointer(""),
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedEntity: &entity.Zone{
				ID:          testID,
				ConcertID:   testConcertID,
				Name:        testName,
				Description: pointer.ToPointer(""),
				CreatedAt:   testCreatedAt,
				UpdatedAt:   testUpdatedAt,
			},
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			result := tt.input.ToEntity()

			// Assert
			if tt.expectedNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedEntity.ID, result.ID)
				assert.Equal(t, tt.expectedEntity.ConcertID, result.ConcertID)
				assert.Equal(t, tt.expectedEntity.Name, result.Name)

				if tt.expectedEntity.Description != nil {
					require.NotNil(t, result.Description)
					assert.Equal(t, *tt.expectedEntity.Description, *result.Description)
				} else {
					assert.Nil(t, result.Description)
				}

				assert.Equal(t, tt.expectedEntity.CreatedAt.UTC(), result.CreatedAt.UTC())
				assert.Equal(t, tt.expectedEntity.UpdatedAt.UTC(), result.UpdatedAt.UTC())
			}
		})
	}
}

func TestZones_ToEntities(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()
	testConcertID1 := uuid.New()
	testConcertID2 := uuid.New()
	testName1 := "VIP Section"
	testName2 := "Regular Section"
	testDescription1 := "Premium seating area"
	testDescription2 := "Standard seating area"
	testCreatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	testUpdatedAt := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		input            zonerepo.Zones
		expectedEntities []entity.Zone
		expectedLength   int
		expectedNil      bool
	}{
		{
			name: "successful conversion with multiple zones",
			input: zonerepo.Zones{
				{
					Zones: model.Zones{
						ID:          testID1,
						ConcertID:   testConcertID1,
						Name:        testName1,
						Description: &testDescription1,
						CreatedAt:   testCreatedAt,
						UpdatedAt:   testUpdatedAt,
					},
				},
				{
					Zones: model.Zones{
						ID:          testID2,
						ConcertID:   testConcertID2,
						Name:        testName2,
						Description: &testDescription2,
						CreatedAt:   testCreatedAt,
						UpdatedAt:   testUpdatedAt,
					},
				},
			},
			expectedEntities: []entity.Zone{
				{
					ID:          testID1,
					ConcertID:   testConcertID1,
					Name:        testName1,
					Description: &testDescription1,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
				{
					ID:          testID2,
					ConcertID:   testConcertID2,
					Name:        testName2,
					Description: &testDescription2,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedLength: 2,
			expectedNil:    false,
		},
		{
			name: "successful conversion with mixed descriptions",
			input: zonerepo.Zones{
				{
					Zones: model.Zones{
						ID:          testID1,
						ConcertID:   testConcertID1,
						Name:        testName1,
						Description: &testDescription1,
						CreatedAt:   testCreatedAt,
						UpdatedAt:   testUpdatedAt,
					},
				},
				{
					Zones: model.Zones{
						ID:          testID2,
						ConcertID:   testConcertID2,
						Name:        testName2,
						Description: nil,
						CreatedAt:   testCreatedAt,
						UpdatedAt:   testUpdatedAt,
					},
				},
			},
			expectedEntities: []entity.Zone{
				{
					ID:          testID1,
					ConcertID:   testConcertID1,
					Name:        testName1,
					Description: &testDescription1,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
				{
					ID:          testID2,
					ConcertID:   testConcertID2,
					Name:        testName2,
					Description: nil,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedLength: 2,
			expectedNil:    false,
		},
		{
			name:             "empty zones slice",
			input:            zonerepo.Zones{},
			expectedEntities: []entity.Zone{},
			expectedLength:   0,
			expectedNil:      false,
		},
		{
			name: "single zone conversion",
			input: zonerepo.Zones{
				{
					Zones: model.Zones{
						ID:          testID1,
						ConcertID:   testConcertID1,
						Name:        testName1,
						Description: &testDescription1,
						CreatedAt:   testCreatedAt,
						UpdatedAt:   testUpdatedAt,
					},
				},
			},
			expectedEntities: []entity.Zone{
				{
					ID:          testID1,
					ConcertID:   testConcertID1,
					Name:        testName1,
					Description: &testDescription1,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedLength: 1,
			expectedNil:    false,
		},
		{
			name: "conversion with empty names",
			input: zonerepo.Zones{
				{
					Zones: model.Zones{
						ID:          testID1,
						ConcertID:   testConcertID1,
						Name:        "",
						Description: &testDescription1,
						CreatedAt:   testCreatedAt,
						UpdatedAt:   testUpdatedAt,
					},
				},
			},
			expectedEntities: []entity.Zone{
				{
					ID:          testID1,
					ConcertID:   testConcertID1,
					Name:        "",
					Description: &testDescription1,
					CreatedAt:   testCreatedAt,
					UpdatedAt:   testUpdatedAt,
				},
			},
			expectedLength: 1,
			expectedNil:    false,
		},
		{
			name: "conversion with zero time values",
			input: zonerepo.Zones{
				{
					Zones: model.Zones{
						ID:          testID1,
						ConcertID:   testConcertID1,
						Name:        testName1,
						Description: &testDescription1,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				},
			},
			expectedEntities: []entity.Zone{
				{
					ID:          testID1,
					ConcertID:   testConcertID1,
					Name:        testName1,
					Description: &testDescription1,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				},
			},
			expectedLength: 1,
			expectedNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			result := tt.input.ToEntities()

			// Assert
			if tt.expectedNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Len(t, *result, tt.expectedLength)

				actualEntities := *result
				for i, expectedEntity := range tt.expectedEntities {
					if i < len(actualEntities) {
						assert.Equal(t, expectedEntity.ID, actualEntities[i].ID)
						assert.Equal(t, expectedEntity.ConcertID, actualEntities[i].ConcertID)
						assert.Equal(t, expectedEntity.Name, actualEntities[i].Name)

						if expectedEntity.Description != nil {
							require.NotNil(t, actualEntities[i].Description)
							assert.Equal(t, *expectedEntity.Description, *actualEntities[i].Description)
						} else {
							assert.Nil(t, actualEntities[i].Description)
						}

						assert.Equal(t, expectedEntity.CreatedAt.UTC(), actualEntities[i].CreatedAt.UTC())
						assert.Equal(t, expectedEntity.UpdatedAt.UTC(), actualEntities[i].UpdatedAt.UTC())
					}
				}
			}
		})
	}
}

func TestZones_ToEntities_EmptyAndNilChecks(t *testing.T) {
	tests := []struct {
		name   string
		input  zonerepo.Zones
		assert func(t *testing.T, result *entity.Zones)
	}{
		{
			name:  "nil input",
			input: nil,
			assert: func(t *testing.T, result *entity.Zones) {
				require.NotNil(t, result)
				assert.Equal(t, 0, len(*result))
			},
		},
		{
			name:  "empty slice",
			input: zonerepo.Zones{},
			assert: func(t *testing.T, result *entity.Zones) {
				require.NotNil(t, result)
				assert.Equal(t, 0, len(*result))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToEntities()
			tt.assert(t, result)
		})
	}
}

func TestZones_ToEntities_ReturnType(t *testing.T) {
	// Test that ToEntities always returns a pointer to entity.Zones
	input := zonerepo.Zones{}
	result := input.ToEntities()

	// Check that it's a pointer
	assert.IsType(t, &entity.Zones{}, result)
	assert.NotNil(t, result)

	// Check that it's not a nil pointer
	assert.NotNil(t, pointer.ToPointer(entity.Zones{}))
}
