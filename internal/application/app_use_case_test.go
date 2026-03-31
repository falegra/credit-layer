package application_test

import (
	"context"
	"credit-layer/internal/application"
	"credit-layer/internal/domain"
	"credit-layer/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateApp_Success(t *testing.T) {
	mockRepo := mocks.NewAppRepository(t)
	ctx := context.Background()

	mockRepo.On("ExistsByName", ctx, "Mi App").Return(false, nil)
	mockRepo.On("Create", ctx, "Mi App", mock.AnythingOfType("string")).Return(&domain.App{
		Name:   "Mi App",
		APIKey: "cl_abc123",
	}, nil)

	useCase := application.NewAppUseCase(mockRepo)
	app, err := useCase.CreateApp(ctx, "Mi App")

	assert.NoError(t, err)
	assert.Equal(t, "Mi App", app.Name)
	assert.NotEmpty(t, app.APIKey)
}

func TestCreateApp_NameTaken(t *testing.T) {
	mockRepo := mocks.NewAppRepository(t)
	ctx := context.Background()

	mockRepo.On("ExistsByName", ctx, "Mi App").Return(true, nil)

	useCase := application.NewAppUseCase(mockRepo)
	app, err := useCase.CreateApp(ctx, "Mi App")

	assert.Nil(t, app)
	assert.ErrorIs(t, err, application.ErrAppNameTaken)
}

func TestCreateApp_EmptyName(t *testing.T) {
	mockRepo := mocks.NewAppRepository(t)
	ctx := context.Background()

	useCase := application.NewAppUseCase(mockRepo)
	app, err := useCase.CreateApp(ctx, "")

	assert.Nil(t, app)
	assert.Error(t, err)
}
