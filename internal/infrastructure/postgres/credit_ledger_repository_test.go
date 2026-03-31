package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"credit-layer/internal/infrastructure/postgres"
)

func ptr(s string) *string {
	return &s
}

func setupApp(t *testing.T) string {
	t.Helper()

	appRepo := postgres.NewAppRepository(testPool)
	app, err := appRepo.Create(context.Background(), "ledger-test-app", "ledger-test-api-key")
	if err != nil {
		t.Fatalf("setupApp: failed to create app: %v", err)
	}

	return app.ID.String()
}

func TestAddCredits_Success(t *testing.T) {
	cleanDB(t)
	appID := setupApp(t)

	repo := postgres.NewCreditLedgerRepository(testPool)

	entry, err := repo.AddCredits(
		context.Background(),
		appID,
		"user-1",
		100,
		ptr("initial credits"),
		ptr("idem-key-add-1"),
	)

	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, int64(100), entry.Amount)
	assert.Equal(t, "user-1", entry.UserID)
}

func TestAddCredits_Idempotency(t *testing.T) {
	cleanDB(t)
	appID := setupApp(t)

	repo := postgres.NewCreditLedgerRepository(testPool)

	_, err := repo.AddCredits(
		context.Background(),
		appID,
		"user-1",
		100,
		ptr("first add"),
		ptr("idem-key-idempotent"),
	)
	assert.NoError(t, err)

	_, err = repo.AddCredits(
		context.Background(),
		appID,
		"user-1",
		100,
		ptr("second add same key"),
		ptr("idem-key-idempotent"),
	)
	assert.NoError(t, err)

	balance, err := repo.GetBalance(context.Background(), appID, "user-1")
	assert.NoError(t, err)
	assert.Equal(t, int64(100), balance)
}

func TestDeductCredits_Success(t *testing.T) {
	cleanDB(t)
	appID := setupApp(t)

	repo := postgres.NewCreditLedgerRepository(testPool)

	_, err := repo.AddCredits(
		context.Background(),
		appID,
		"user-1",
		100,
		ptr("funding"),
		ptr("idem-key-fund"),
	)
	assert.NoError(t, err)

	entry, err := repo.DeductCredits(
		context.Background(),
		appID,
		"user-1",
		-30,
		ptr("deduction"),
		ptr("idem-key-deduct"),
	)

	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, int64(-30), entry.Amount)
}

func TestGetBalance_Empty(t *testing.T) {
	cleanDB(t)
	appID := setupApp(t)

	repo := postgres.NewCreditLedgerRepository(testPool)

	balance, err := repo.GetBalance(context.Background(), appID, "user-with-no-credits")

	assert.NoError(t, err)
	assert.Equal(t, int64(0), balance)
}

func TestGetBalance_Success(t *testing.T) {
	cleanDB(t)
	appID := setupApp(t)

	repo := postgres.NewCreditLedgerRepository(testPool)

	_, err := repo.AddCredits(
		context.Background(),
		appID,
		"user-1",
		100,
		ptr("add credits"),
		ptr("idem-key-balance-add"),
	)
	assert.NoError(t, err)

	_, err = repo.DeductCredits(
		context.Background(),
		appID,
		"user-1",
		-30,
		ptr("deduct credits"),
		ptr("idem-key-balance-deduct"),
	)
	assert.NoError(t, err)

	balance, err := repo.GetBalance(context.Background(), appID, "user-1")

	assert.NoError(t, err)
	assert.Equal(t, int64(70), balance)
}
