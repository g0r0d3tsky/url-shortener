package repository

import (
	"cleaner-service/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (s *StorageURL) GetURL(ctx context.Context, monthAmount int) ([]*domain.Url, error) {
	urls := []*domain.Url{}
	rows, err := s.db.Query(
		ctx,
		`
    SELECT id, long_url, short_url, visited_at FROM url_data WHERE url_data.visited_at < NOW() - INTERVAL '1 month' * $1
    `,
		monthAmount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get urls: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		url := &domain.Url{}
		if err := rows.Scan(&url.ID, &url.OriginalURL, &url.ShortURL, &url.VisitedAT); err != nil {
			return nil, fmt.Errorf("scan url: %w", err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *StorageURL) DeleteURL(ctx context.Context, url *domain.Url) error {
	_, err := s.db.Exec(ctx,
		`DELETE FROM url_data WHERE id = $1`, url.ID)
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}

	return nil
}
