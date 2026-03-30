package application

import (
	"context"
	"credit-layer/internal/domain"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var (
	ErrAppNotFound  = errors.New("app not found")
	ErrAppNameTaken = errors.New("app name already taken")
)

type AppUseCase struct {
	appRepo domain.AppRepository
}

func NewAppUseCase(appRepo domain.AppRepository) *AppUseCase {
	return &AppUseCase{
		appRepo: appRepo,
	}
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return "cl_" + hex.EncodeToString(bytes), nil
}

func (aUC *AppUseCase) ExistsByName(ctx context.Context, name string) (bool, error) {
	exists, err := aUC.appRepo.ExistsByName(ctx, name)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (aUC *AppUseCase) CreateApp(ctx context.Context, name string) (*domain.App, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	exists, err := aUC.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrAppNameTaken
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	return aUC.appRepo.Create(ctx, name, apiKey)
}

func (aUC *AppUseCase) GetAppByAPIKey(ctx context.Context, apiKey string) (*domain.App, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey is required")
	}

	app, err := aUC.appRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, ErrAppNotFound
	}

	return app, nil
}
