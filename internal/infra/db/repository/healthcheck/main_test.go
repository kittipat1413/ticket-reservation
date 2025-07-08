package healthcheckrepo_test

import (
	"testing"
	healthcheckrepo "ticket-reservation/internal/infra/db/repository/healthcheck"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	defer h.done()

	// Create a mock transaction
	txDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer txDB.Close()

	transactionDB := sqlx.NewDb(txDB, "sqlmock")

	// Execute
	txRepo := h.repository.WithTx(transactionDB)

	// Assert
	assert.NotNil(t, txRepo)

	// Verify that the returned repository is a new instance with the transaction
	assert.NotEqual(t, h.repository, txRepo, "WithTx should return a new repository instance")
}
