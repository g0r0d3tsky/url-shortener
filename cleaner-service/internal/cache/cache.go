package cache

import (
	"cleaner-service/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Url struct {
	ID          uuid.UUID `json:"id,omitempty"`
	OriginalURL string    `json:"original_url,omitempty"`
	ShortURL    string    `json:"short_url,omitempty"`
	VisitedAT   time.Time `json:"visited_at"`
}

func Connect(c *config.Config) (*redis.Client, error) {
	connectionString := c.RedisDSN()
	opts, err := redis.ParseURL(connectionString)
	if err != nil {
		slog.Error("parsing redis url")
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection was refused %w", err)
	}
	if status != "PONG" {
		return nil, fmt.Errorf("redis connection was not successful")
	}

	return client, nil
}
