package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"credit-layer/internal/infrastructure/postgres"
)

func TestCreate_Success(t *testing.T) {
	cleanDB(t)

	repo := postgres.NewAppRepository(testPool)

	app, err := repo.Create(context.Background(), "test-app", "test-api-key-123")

	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.NotEmpty(t, app.Name)
	assert.NotEmpty(t, app.APIKey)
}

func TestGetByAPIKey_Success(t *testing.T) {
	cleanDB(t)

	repo := postgres.NewAppRepository(testPool)

	created, err := repo.Create(context.Background(), "my-app", "my-secret-key")
	assert.NoError(t, err)
	assert.NotNil(t, created)

	found, err := repo.GetByAPIKey(context.Background(), "my-secret-key")

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, created.Name, found.Name)
}

func TestGetByAPIKey_NotFound(t *testing.T) {
	cleanDB(t)

	repo := postgres.NewAppRepository(testPool)

	found, err := repo.GetByAPIKey(context.Background(), "nonexistent-api-key")

	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestExistsByName_Exists(t *testing.T) {
	cleanDB(t)

	repo := postgres.NewAppRepository(testPool)

	_, err := repo.Create(context.Background(), "existing-app", "some-api-key")
	assert.NoError(t, err)

	exists, err := repo.ExistsByName(context.Background(), "existing-app")

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestExistsByName_NotExists(t *testing.T) {
	cleanDB(t)

	repo := postgres.NewAppRepository(testPool)

	exists, err := repo.ExistsByName(context.Background(), "ghost-app")

	assert.NoError(t, err)
	assert.False(t, exists)
}
