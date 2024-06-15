package usecase

import (
	"context"
	"fmt"
	"time"
	"url-service/url-service/internal/domain"
	"url-service/url-service/internal/shorter"

	"github.com/google/uuid"
)

type UrlRepo interface {
	GetUrl(ctx context.Context, short string) (*domain.Url, error)
	CreateShort(ctx context.Context, url *domain.Url) (*uuid.UUID, error)
}

type UrlService struct {
	repoUrl UrlRepo
	repoKey KeyRepo
}

func newUrlService(repoUrl UrlRepo, repoKey KeyRepo) *UrlService {
	return &UrlService{repoUrl: repoUrl, repoKey: repoKey}
}

func (us *UrlService) CreateShort(ctx context.Context, originalUrl string) (*domain.Url, error) {
	key, err := us.repoKey.GetFreeKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get free key: %w", err)
	}

	var url *domain.Url

	if key == nil {
		newKey, err := us.repoKey.GetNextKeyFromSequence(ctx)
		if err != nil {
			return nil, fmt.Errorf("get next key from sequence: %w", err)
		}

		short := shorter.Shorten(*newKey)

		url = &domain.Url{
			OriginalURL: originalUrl,
			ShortURL:    short,
			ExpiresAT:   time.Now().AddDate(1, 0, 0),
		}

		id, err := us.repoUrl.CreateShort(ctx, url)

		if err != nil {
			return nil, fmt.Errorf("create short: %w", err)
		}

		key = &domain.Key{
			Key:   *newKey,
			Code:  short,
			UrlID: *id,
		}

		err = us.repoKey.CreateNewKey(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("create key: %w", err)

		}
	} else {
		url = &domain.Url{
			OriginalURL: originalUrl,
			ShortURL:    key.Code,
			ExpiresAT:   time.Now().AddDate(1, 0, 0),
		}

		id, err := us.repoUrl.CreateShort(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("create short: %w", err)
		}

		key.UrlID = *id
		err = us.repoKey.CreateNewKey(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("create key: %w", err)
		}
	}
	return url, nil
}

func (us *UrlService) GetUrl(ctx context.Context, short string) (*domain.Url, error) {
	url, err := us.repoUrl.GetUrl(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("get url: %w", err)
	}
	return url, nil
}
