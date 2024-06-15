package repository

import (
	"context"
	"fmt"
	"url-service/url-service/internal/domain"

	"github.com/google/uuid"
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

func (s *StorageURL) CreateShort(ctx context.Context, url *domain.Url) (*uuid.UUID, error) {
	url.ID = uuid.New()
	if _, err := s.db.Exec(ctx,
		`INSERT INTO url_data (id, long_url, short_url, expires_at) 
			VALUES ($1, $2, $3, $4)`,
		&url.ID, &url.OriginalURL, &url.ShortURL, &url.ExpiresAT,
	); err != nil {
		return nil, fmt.Errorf("create url: %w", err)
	}
	return &url.ID, nil
}

func (s *StorageURL) GetUrl(ctx context.Context, short string) (*domain.Url, error) {
	var url *domain.Url
	if err := s.db.QueryRow(
		ctx, `
		SELECT id, long_url, short_url, expires_at FROM url_data WHERE url_data.short_url= $1
	`, short).Scan(&url.ID, &url.OriginalURL, &url.ShortURL, &url.ExpiresAT); err != nil {
		return nil, fmt.Errorf("get url: %w", err)
	}

	return url, nil
}
