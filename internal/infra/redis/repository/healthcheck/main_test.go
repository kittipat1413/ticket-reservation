package healthcheckrepo_test

import (
	"testing"
	healthcheckrepo "ticket-reservation/internal/infra/redis/repository/healthcheck"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewHealthCheckRepository(t *testing.T) {
	client, _ := redismock.NewClientMock()

	// Execute
	repo := healthcheckrepo.NewHealthCheckRepository(client)

	// Assert
	assert.NotNil(t, repo)
}
