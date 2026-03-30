package application

import (
	"context"
	"credit-layer/internal/domain"
	"errors"
)

var (
	ErrInsufficientCredits = errors.New("insufficient credits")
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
)

type CreditLedgerUseCase struct {
	creditLedgerRepo domain.CreditLedgerRepository
}

func NewCreditLedgerUseCase(creditLegderRepo domain.CreditLedgerRepository) *CreditLedgerUseCase {
	return &CreditLedgerUseCase{
		creditLedgerRepo: creditLegderRepo,
	}
}

func (cLUC *CreditLedgerUseCase) GetBalance(ctx context.Context, appID, userID string) (int64, error) {
	if appID == "" || userID == "" {
		return 0, errors.New("appID and userID are required")
	}

	return cLUC.creditLedgerRepo.GetBalance(ctx, appID, userID)
}

func (cLUC *CreditLedgerUseCase) AddCredits(ctx context.Context, appID, userID string, amount int64, description, idempotencyKey *string) (*domain.CreditLedger, error) {
	if appID == "" || userID == "" || description == nil || idempotencyKey == nil {
		return nil, errors.New("appID, userID, description and idempotencyKey are required")
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	return cLUC.creditLedgerRepo.AddCredits(ctx, appID, userID, amount, description, idempotencyKey)
}

func (cLUC *CreditLedgerUseCase) DeductCredits(ctx context.Context, appID, userID string, amount int64, description, idempotencyKey *string) (*domain.CreditLedger, error) {
	if appID == "" || userID == "" || description == nil || idempotencyKey == nil {
		return nil, errors.New("appID, userID, description and idempotencyKey are required")
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	balance, err := cLUC.creditLedgerRepo.GetBalance(ctx, appID, userID)
	if err != nil {
		return nil, err
	}

	if balance < amount {
		return nil, ErrInsufficientCredits
	}

	return cLUC.creditLedgerRepo.DeductCredits(ctx, appID, userID, -amount, description, idempotencyKey)
}
