package repository

import (
	"context"
	"fmt"
	"url-service/internal/domain"

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

func (s *StorageURL) GetURL(ctx context.Context, short string) (*domain.Url, error) {
	url := &domain.Url{}
	if err := s.db.QueryRow(
		ctx, `
		SELECT id, long_url, short_url, visited_at FROM url_data WHERE url_data.short_url= $1
	`, short).Scan(&url.ID, &url.OriginalURL, &url.ShortURL, &url.VisitedAT); err != nil {
		return nil, fmt.Errorf("get url: %w", err)
	}

	return url, nil
}
