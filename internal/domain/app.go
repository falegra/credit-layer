package domain

import (
	"time"

	"github.com/google/uuid"
)

type App struct {
	ID        uuid.UUID
	Name      string
	APIKey    string
	CreatedAt time.Time
}
