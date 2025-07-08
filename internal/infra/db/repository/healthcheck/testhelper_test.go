package healthcheckrepo_test

import (
	"testing"
	"ticket-reservation/internal/domain/repository"
	healthcheckrepo "ticket-reservation/internal/infra/db/repository/healthcheck"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type helper struct {
	repository repository.HealthCheckRepository
	mock       sqlmock.Sqlmock
	done       func()
}

func initTest(t *testing.T) *helper {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &helper{
		repository: healthcheckrepo.NewHealthCheckRepository(sqlxDB),
		mock:       mock,
		done: func() {
			db.Close()
		},
	}
}
