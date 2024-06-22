package usecase

import (
	"cleaner-service/internal/domain"
	"context"
	"fmt"
)

//go:generate mockgen -source=cleaner.go -destination=mocks/cleaner_mock.go
type UrlRepo interface {
	GetURL(ctx context.Context, monthAmount int) ([]*domain.Url, error)
	DeleteURL(ctx context.Context, url *domain.Url) error
}

type Cache interface {
	GetURL(ctx context.Context, shortURL string) (*domain.Url, error)
}

type KeyRepo interface {
	UpdateKey(ctx context.Context, url *domain.Url) error
}

type URLService struct {
	repoUrl     UrlRepo
	repoKey     KeyRepo
	cache       Cache
	monthAmount int
}

func NewURLService(repoUrl UrlRepo, keyRepo KeyRepo, cache Cache, monthAmount int) *URLService {
	return &URLService{repoUrl: repoUrl, repoKey: keyRepo, cache: cache, monthAmount: monthAmount}
}

func (us *URLService) CleanURL(ctx context.Context) error {
	urls, err := us.repoUrl.GetURL(ctx, us.monthAmount)
	if err != nil {
		return fmt.Errorf("get urls: %w", err)
	} else if urls == nil {
		return nil
	}
	for _, url := range urls {
		urlCache, err := us.cache.GetURL(ctx, url.ShortURL)
		if err != nil {
			return fmt.Errorf("get url from cache: %w", err)
		} else if urlCache == nil {
			if err := us.repoUrl.DeleteURL(ctx, url); err != nil {
				return fmt.Errorf("delete url: %w", err)
			}
			if err := us.repoKey.UpdateKey(ctx, url); err != nil {
				return fmt.Errorf("update key: %w", err)
			}
			continue
		}

		if urlCache.VisitedAT.IsZero() || (!url.VisitedAT.IsZero() && urlCache.VisitedAT.After(url.VisitedAT)) {
			continue
		}
		if err := us.repoUrl.DeleteURL(ctx, url); err != nil {
			return fmt.Errorf("delete url: %w", err)
		}
		if err := us.repoKey.UpdateKey(ctx, url); err != nil {
			return fmt.Errorf("update key: %w", err)
		}
	}
	return nil
}
