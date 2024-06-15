package domain

import (
	"time"

	"github.com/google/uuid"
)

type Url struct {
	ID          uuid.UUID
	OriginalURL string
	ShortURL    string
	ExpiresAT   time.Time
}

type Key struct {
	ID    uuid.UUID
	UrlID uuid.UUID
	Key   uint64
	Code  string
}
