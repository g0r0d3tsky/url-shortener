package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"url-service/url-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	GetURL(ctx context.Context, short string) (*domain.Url, error)
}

type Redis interface {
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type Cache struct {
	redis Redis
	key   string
}

func NewCache(client Redis, key string) Cache {
	cache := Cache{
		redis: client,
		key:   key,
	}
	return cache
}

func (c Cache) GetURL(ctx context.Context, shortURL string) (*domain.Url, error) {
	url, err := c.redis.HGet(ctx, c.key, shortURL).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get url from cache: %w", err)
	}

	var urlRedis *Url

	err = json.Unmarshal([]byte(url), &urlRedis)
	if err != nil {
		return nil, fmt.Errorf("unmarshal url: %w", err)
	}

	urlDomain := &domain.Url{
		ID:          urlRedis.ID,
		OriginalURL: urlRedis.OriginalURL,
		ShortURL:    urlRedis.ShortURL,
		VisitedAT:   urlRedis.VisitedAT,
	}
	return urlDomain, nil
}

func (c Cache) SetURL(ctx context.Context, url *domain.Url) error {
	urlRedis := &Url{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		ShortURL:    url.ShortURL,
		VisitedAT:   url.VisitedAT,
	}

	urlJson, err := json.Marshal(urlRedis)

	err = c.redis.HSet(ctx, c.key, urlRedis.ShortURL, urlJson).Err()
	if err != nil {
		return fmt.Errorf("store url in cache: %w", err)
	}

	return nil
}
