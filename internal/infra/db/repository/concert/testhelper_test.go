package concertrepo_test

import (
	"testing"
	"ticket-reservation/internal/domain/repository"
	concertrepo "ticket-reservation/internal/infra/db/repository/concert"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type helper struct {
	repository repository.ConcertRepository
	mock       sqlmock.Sqlmock
	done       func()
}

func initTest(t *testing.T) *helper {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &helper{
		repository: concertrepo.NewConcertRepository(sqlxDB),
		mock:       mock,
		done: func() {
			db.Close()
		},
	}
}
