package repository

import (
	"context"
	"errors"
	"fmt"
	"url-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *StorageKey) GetNextKeyFromSequence(ctx context.Context) (*uint64, error) {
	var key uint64
	err := s.db.QueryRow(ctx, "SELECT nextval('key_serial_seq')").Scan(&key)
	if err != nil {
		return nil, fmt.Errorf("get next key from sequence: %v", err)
	}

	return &key, nil
}

func (s *StorageKey) CreateNewKey(ctx context.Context, key *domain.Key) error {
	key.ID = uuid.New()
	if _, err := s.db.Exec(ctx,
		`INSERT INTO url_keys (id, key_serial, encode, url_id) 
			VALUES ($1, $2, $3, $4)`,
		&key.ID, &key.Key, &key.Code, &key.UrlID,
	); err != nil {
		return fmt.Errorf("create url: %w", err)
	}
	return nil
}

func (s *StorageKey) GetFreeKey(ctx context.Context) (*domain.Key, error) {
	key := &domain.Key{}
	err := s.db.QueryRow(
		ctx, `
            SELECT id, key_serial, encode, url_id FROM url_keys WHERE url_keys.url_id IS NULL
        `).Scan(&key.ID, &key.Key, &key.Code, &key.UrlID)
	if err != nil {
		errors.Unwrap(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get url: %w", err)
	}
	return key, nil
}

func (s *StorageKey) UpdateKey(ctx context.Context, key *domain.Key) error {
	if _, err := s.db.Exec(ctx,
		`UPDATE url_keys SET url_id = $1 WHERE id = $2`, key.UrlID, key.ID); err != nil {
		return fmt.Errorf("update url: %w", err)
	}
	return nil
}
