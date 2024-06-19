package usecase

import (
	"context"
	"fmt"
	"time"
	"url-service/url-service/internal/domain"
	"url-service/url-service/internal/kafka"
	"url-service/url-service/internal/shorter"

	"github.com/google/uuid"
)

type UrlRepo interface {
	GetURL(ctx context.Context, short string) (*domain.Url, error)
}

type Cache interface {
	GetURL(ctx context.Context, short string) (*domain.Url, error)
	SetURL(ctx context.Context, url *domain.Url) error
}

type Broker interface {
	Push(topic string, message *domain.Url) error
}

type URLService struct {
	repoUrl    UrlRepo
	repoKey    KeyRepo
	cache      Cache
	kafka      Broker
	kafkaTopic string
}

func NewURLService(repoUrl UrlRepo, keyRepo KeyRepo, cache Cache, producer *kafka.Producer, kafkaTopic string) *URLService {
	return &URLService{repoUrl: repoUrl, repoKey: keyRepo, cache: cache, kafka: producer, kafkaTopic: kafkaTopic}
}

func (us *URLService) GetURL(ctx context.Context, short string) (*domain.Url, error) {
	url, err := us.cache.GetURL(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("get url: %w", err)
	} else if url == nil {
		url, err := us.repoUrl.GetURL(ctx, short)
		if err != nil {
			return nil, fmt.Errorf("get url: %w", err)
		}
		url.VisitedAT = time.Now()
		err = us.cache.SetURL(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("set url: %w", err)
		}

		err = us.kafka.Push(us.kafkaTopic, url)
		if err != nil {
			return nil, fmt.Errorf("pushing kafka: %w", err)
		}
		return url, nil
	}

	return url, nil
}

func (us *URLService) CreateShortURL(ctx context.Context, originalUrl string) (*domain.Url, error) {
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
		id := uuid.New()

		url = &domain.Url{
			ID:          id,
			OriginalURL: originalUrl,
			ShortURL:    short,
			VisitedAT:   time.Now(),
		}

		err = us.kafka.Push(us.kafkaTopic, url)
		if err != nil {
			return nil, fmt.Errorf("pushing kafka: %w", err)
		}

		key = &domain.Key{
			Key:   *newKey,
			Code:  short,
			UrlID: id,
		}

		err = us.repoKey.CreateNewKey(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("create key: %w", err)

		}
	} else {
		id := uuid.New()
		url = &domain.Url{
			ID:          id,
			OriginalURL: originalUrl,
			ShortURL:    key.Code,
			VisitedAT:   time.Now(),
		}

		err = us.kafka.Push(us.kafkaTopic, url)
		if err != nil {
			return nil, fmt.Errorf("pushing kafka: %w", err)
		}

		key.UrlID = id
		err = us.repoKey.UpdateKey(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("update key: %w", err)
		}
	}
	return url, nil
}
