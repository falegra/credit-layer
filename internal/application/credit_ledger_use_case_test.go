package application_test

import (
	"context"
	"credit-layer/internal/application"
	"credit-layer/internal/domain"
	"credit-layer/internal/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddCredits_Success(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	mockRepo.On("AddCredits", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64"), mock.Anything, mock.Anything).Return(&domain.CreditLedger{
		ID:     uuid.New(),
		AppID:  uuid.New(),
		UserID: "user",
		Amount: 100,
	}, nil)

	desc := "Compra de créditos"
	iKey := "test_001"

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	c, err := useCase.AddCredits(ctx, "app", "user", 100, &desc, &iKey)

	assert.NoError(t, err)
	assert.NotEmpty(t, c)
}

func TestAddCredits_InvalidAmount(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	desc := "Compra de créditos"
	iKey := "test_001"

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	c, err := useCase.AddCredits(ctx, "app", "user", -1, &desc, &iKey)

	assert.Nil(t, c)
	assert.ErrorIs(t, err, application.ErrInvalidAmount)
}

func TestDeductCredits_RequiredFields(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	desc := "Uso del servicio"
	iKey := "test_001"

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	c, err := useCase.DeductCredits(ctx, "", "user", 10, &desc, &iKey)

	assert.Nil(t, c)
	assert.Error(t, err)
}

func TestDeductCredits_InvalidAmount(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	desc := "Uso del servicio"
	iKey := "test_001"

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	c, err := useCase.DeductCredits(ctx, "app", "user", 0, &desc, &iKey)

	assert.Nil(t, c)
	assert.ErrorIs(t, err, application.ErrInvalidAmount)
}

func TestDeductCredits_InsufficientCredits(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	mockRepo.On("GetBalance", ctx, "app", "user").Return(int64(50), nil)

	desc := "Uso del servicio"
	iKey := "test_001"

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	c, err := useCase.DeductCredits(ctx, "app", "user", 100, &desc, &iKey)

	assert.Nil(t, c)
	assert.ErrorIs(t, err, application.ErrInsufficientCredits)
}

func TestDeductCredits_Success(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	mockRepo.On("GetBalance", ctx, "app", "user").Return(int64(200), nil)
	mockRepo.On("DeductCredits", ctx, "app", "user", int64(-100), mock.Anything, mock.Anything).Return(&domain.CreditLedger{
		ID:     uuid.New(),
		AppID:  uuid.New(),
		UserID: "user",
		Amount: -100,
	}, nil)

	desc := "Uso del servicio"
	iKey := "test_001"

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	c, err := useCase.DeductCredits(ctx, "app", "user", 100, &desc, &iKey)

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, int64(-100), c.Amount)
}

func TestGetBalance_RequiredFields(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	balance, err := useCase.GetBalance(ctx, "", "user")

	assert.Zero(t, balance)
	assert.Error(t, err)
}

func TestGetBalance_Success(t *testing.T) {
	mockRepo := mocks.NewCreditLedgerRepository(t)
	ctx := context.Background()

	mockRepo.On("GetBalance", ctx, "app", "user").Return(int64(150), nil)

	useCase := application.NewCreditLedgerUseCase(mockRepo)
	balance, err := useCase.GetBalance(ctx, "app", "user")

	assert.NoError(t, err)
	assert.Equal(t, int64(150), balance)
}
