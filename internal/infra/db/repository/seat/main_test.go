package seatrepo_test

import (
	"testing"
	"ticket-reservation/internal/domain/repository"
	seatrepo "ticket-reservation/internal/infra/db/repository/seat"
	"ticket-reservation/pkg/testhelper"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTest(t *testing.T) *testhelper.RepoTestHelper[repository.SeatRepository] {
	return testhelper.InitRepoTest(t, func(db *sqlx.DB) repository.SeatRepository {
		return seatrepo.NewSeatRepository(db)
	})
}

func TestNewSeatRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockDB := sqlx.NewDb(db, "sqlmock")

	// Execute
	repo := seatrepo.NewSeatRepository(mockDB)

	// Assert
	assert.NotNil(t, repo)
}

func TestSeatRepositoryImpl_WithTx(t *testing.T) {
	h := initTest(t)
	defer h.Done()

	// Create a mock transaction database
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
