package repository

import (
	"cleaner-service/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageKey struct {
	db *pgxpool.Pool
}

func NewStorageKey(dbPool *pgxpool.Pool) StorageKey {
	StorageKey := StorageKey{
		db: dbPool,
	}
	return StorageKey
}

func (s *StorageKey) UpdateKey(ctx context.Context, url *domain.Url) error {
	if _, err := s.db.Exec(ctx,
		`UPDATE url_keys SET url_id = NULL WHERE id = $1`, url.ID); err != nil {
		return fmt.Errorf("update url: %w", err)
	}
	return nil
}
