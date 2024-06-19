package domain

import (
	"time"

	"github.com/google/uuid"
)

type Url struct {
	ID          uuid.UUID
	OriginalURL string
	ShortURL    string
	VisitedAT   time.Time
}
