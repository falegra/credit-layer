package domain

import "context"

type CreditLedgerRepository interface {
	AddCredits(ctx context.Context, appID, userID string, amount int64, description, idempotencyKey *string) (*CreditLedger, error)
	DeductCredits(ctx context.Context, appID, userID string, amount int64, description, idempotencyKey *string) (*CreditLedger, error)
	GetBalance(ctx context.Context, appID, userID string) (int64, error)
}
