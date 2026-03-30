package postgres

import (
	"context"
	"credit-layer/internal/domain"
	"credit-layer/internal/infrastructure/postgres/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppRepository struct {
	queries *db.Queries
}

func NewAppRepository(pool *pgxpool.Pool) *AppRepository {
	return &AppRepository{
		queries: db.New(pool),
	}
}

func (aR *AppRepository) Create(ctx context.Context, name, apiKey string) (*domain.App, error) {
	row, err := aR.queries.CreateApp(ctx, db.CreateAppParams{
		Name:   name,
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &domain.App{
		ID:        row.ID,
		Name:      row.Name,
		APIKey:    row.ApiKey,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (aR *AppRepository) GetByAPIKey(ctx context.Context, apiKey string) (*domain.App, error) {
	row, err := aR.queries.GetAppByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return &domain.App{
		ID:        row.ID,
		Name:      row.Name,
		APIKey:    row.ApiKey,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (aR *AppRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return aR.queries.ExistsAppByName(ctx, name)
}
