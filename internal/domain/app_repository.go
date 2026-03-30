package domain

import "context"

type AppRepository interface {
	GetByAPIKey(ctx context.Context, apiKey string) (*App, error)
	Create(ctx context.Context, name, apiKey string) (*App, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
