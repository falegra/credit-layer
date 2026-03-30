package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreditLedger struct {
	ID             uuid.UUID
	AppID          uuid.UUID
	UserID         string
	Amount         int64
	Description    *string
	IdempotencyKey *string
	CreatedAt      time.Time
}
