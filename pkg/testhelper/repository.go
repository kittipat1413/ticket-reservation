package testhelper

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// RepoTestHelper is a generic test helper that can work with any repository type
type RepoTestHelper[T any] struct {
	Repository T
	Mock       sqlmock.Sqlmock
	done       func()
}

// NewRepoTestHelper creates a new test helper with a mock database for any repository type
// repoFunc should be a function that takes a *sqlx.DB and returns the repository instance
func NewRepoTestHelper[T any](t *testing.T, repoFunc func(*sqlx.DB) T) *RepoTestHelper[T] {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return &RepoTestHelper[T]{
		Repository: repoFunc(sqlxDB),
		Mock:       mock,
		done: func() {
			db.Close()
		},
	}
}

// Done closes the database connection and cleans up resources
func (h *RepoTestHelper[T]) Done() {
	h.done()
}

// AssertExpectationsMet is a convenience method to check if all mock expectations were met
func (h *RepoTestHelper[T]) AssertExpectationsMet(t *testing.T) {
	require.NoError(t, h.Mock.ExpectationsWereMet())
}
