package repository

import (
	"context"
	"fmt"
	"storage-service/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageURL struct {
	db *pgxpool.Pool
}

func NewStorageURL(dbPool *pgxpool.Pool) StorageURL {
	StorageURL := StorageURL{
		db: dbPool,
	}
	return StorageURL
}

func (s *StorageURL) CreateShort(ctx context.Context, url *domain.Url) error {
	if _, err := s.db.Exec(ctx,
		`INSERT INTO url_data (id, long_url, short_url, visited_at) 
			VALUES ($1, $2, $3, $4)`,
		&url.ID, &url.OriginalURL, &url.ShortURL, &url.VisitedAT,
	); err != nil {
		return fmt.Errorf("create url: %w", err)
	}
	return nil
}

func (s *StorageURL) UpdateURL(ctx context.Context, url *domain.Url) error {
	if _, err := s.db.Exec(ctx,
		`UPDATE url_data SET visited_at = $1 WHERE id = $2`, time.Now(), url.ID); err != nil {
		return fmt.Errorf("update url: %w", err)
	}
	return nil
}

func (s *StorageURL) CheckExist(ctx context.Context, url *domain.Url) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM url_data WHERE id = $1)", url.ID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
