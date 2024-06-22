package usecase

import (
	"context"
	"fmt"
	"storage-service/internal/domain"
)

//go:generate mockgen -source=services.go -destination=mocks/service_mock.go
type Repository interface {
	CreateShort(ctx context.Context, url *domain.Url) error
	UpdateURL(ctx context.Context, url *domain.Url) error
	CheckExist(ctx context.Context, url *domain.Url) (bool, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateURL(ctx context.Context, url *domain.Url) error {
	exist, err := s.repository.CheckExist(ctx, url)
	if err != nil {
		return fmt.Errorf("checking exist: %w", err)
	}
	if exist {
		err = s.repository.UpdateURL(ctx, url)
		if err != nil {
			return fmt.Errorf("updating message: %w", err)
		}
		return nil
	}
	err = s.repository.CreateShort(ctx, url)
	if err != nil {
		return fmt.Errorf("creating message: %w", err)
	}
	return nil
}
