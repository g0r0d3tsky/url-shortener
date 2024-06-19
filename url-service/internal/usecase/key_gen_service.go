package usecase

import (
	"context"
	"fmt"
	"url-service/url-service/internal/domain"
)

type KeyRepo interface {
	GetNextKeyFromSequence(ctx context.Context) (*uint64, error)
	CreateNewKey(ctx context.Context, key *domain.Key) error
	GetFreeKey(ctx context.Context) (*domain.Key, error)
	UpdateKey(ctx context.Context, key *domain.Key) error
}

type KeyGenService struct {
	repoKey KeyRepo
}

func NewKeyGenService(repo KeyRepo) *KeyGenService {
	return &KeyGenService{repoKey: repo}
}

func (s *KeyGenService) GetNextKeyFromSequence(ctx context.Context) (*uint64, error) {
	key, err := s.repoKey.GetNextKeyFromSequence(ctx)
	if err != nil {
		return nil, fmt.Errorf("get next key from sequence: %w", err)
	}
	return key, nil
}

func (s *KeyGenService) CreateNewKey(ctx context.Context, key *domain.Key) error {
	err := s.repoKey.CreateNewKey(ctx, key)
	if err != nil {
		return fmt.Errorf("create key: %w", err)
	}
	return nil
}

func (s *KeyGenService) GetFreeKey(ctx context.Context) (*domain.Key, error) {
	key, err := s.repoKey.GetFreeKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get free key: %w", err)
	}
	return key, nil
}

func (s *KeyGenService) UpdateKey(ctx context.Context, key *domain.Key) error {
	err := s.repoKey.UpdateKey(ctx, key)
	if err != nil {
		return fmt.Errorf("update key: %w", err)
	}
	return nil
}
