package postgres

import (
	"context"
	"credit-layer/internal/domain"
	"credit-layer/internal/infrastructure/postgres/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreditLedgerRepository struct {
	queries *db.Queries
}

func NewCreditLedgerRepository(pool *pgxpool.Pool) *CreditLedgerRepository {
	return &CreditLedgerRepository{
		queries: db.New(pool),
	}
}

func toDomainCreditLedger(row db.CreditLedger) *domain.CreditLedger {
	return &domain.CreditLedger{
		ID:             row.ID,
		AppID:          row.AppID,
		UserID:         row.UserID,
		Amount:         row.Amount,
		Description:    fromPgtypeText(row.Description),
		IdempotencyKey: fromPgtypeText(row.IdempotencyKey),
		CreatedAt:      row.CreatedAt,
	}
}

func (cLR *CreditLedgerRepository) AddCredits(ctx context.Context, appID, userID string, amount int64, description, idempotencyKey *string) (*domain.CreditLedger, error) {
	parsedAppID, err := uuid.Parse(appID)
	if err != nil {
		return nil, err
	}

	row, err := cLR.queries.AddCredits(ctx, db.AddCreditsParams{
		AppID:          parsedAppID,
		UserID:         userID,
		Amount:         amount,
		Description:    toPgtypeText(description),
		IdempotencyKey: toPgtypeText(idempotencyKey),
	})
	if err != nil {
		return nil, err
	}

	return toDomainCreditLedger(row), nil
}

func (cLR *CreditLedgerRepository) DeductCredits(ctx context.Context, appID, userID string, amount int64, description, idempotencyKey *string) (*domain.CreditLedger, error) {
	parsedAppID, err := uuid.Parse(appID)
	if err != nil {
		return nil, err
	}

	row, err := cLR.queries.DeductCredits(ctx, db.DeductCreditsParams{
		AppID:          parsedAppID,
		UserID:         userID,
		Amount:         amount,
		Description:    toPgtypeText(description),
		IdempotencyKey: toPgtypeText(idempotencyKey),
	})
	if err != nil {
		return nil, err
	}

	return toDomainCreditLedger(row), nil
}

func (r *CreditLedgerRepository) GetBalance(ctx context.Context, appID, userID string) (int64, error) {
	parsedAppID, err := uuid.Parse(appID)
	if err != nil {
		return 0, err
	}

	return r.queries.GetBalance(ctx, db.GetBalanceParams{
		AppID:  parsedAppID,
		UserID: userID,
	})
}

func toPgtypeText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func fromPgtypeText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
