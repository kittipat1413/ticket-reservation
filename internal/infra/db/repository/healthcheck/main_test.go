package healthcheckrepo_test

import (
	"testing"
	"ticket-reservation/internal/domain/repository"
	healthcheckrepo "ticket-reservation/internal/infra/db/repository/healthcheck"
	"ticket-reservation/pkg/testhelper"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTest(t *testing.T) *testhelper.RepoTestHelper[repository.HealthCheckRepository] {
	return testhelper.NewRepoTestHelper(t, func(db *sqlx.DB) repository.HealthCheckRepository {
		return healthcheckrepo.NewHealthCheckRepository(db)
	})
}

func TestNewHealthCheckRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockDB := sqlx.NewDb(db, "sqlmock")

	// Execute
	repo := healthcheckrepo.NewHealthCheckRepository(mockDB)

	// Assert
	assert.NotNil(t, repo)
}

func TestHealthCheckRepositoryImpl_WithTx(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	// Create a mock transaction
	txDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer txDB.Close()

	transactionDB := sqlx.NewDb(txDB, "sqlmock")

	// Execute
	txRepo := h.Repository.WithTx(transactionDB)

	// Assert
	assert.NotNil(t, txRepo)

	// Verify that the returned repository is a new instance with the transaction
	assert.NotEqual(t, h.Repository, txRepo, "WithTx should return a new repository instance")
}
