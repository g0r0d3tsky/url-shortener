package main

import (
	"cleaner-service/config"
	"cleaner-service/internal/cache"
	"cleaner-service/internal/repository"
	"cleaner-service/internal/usecase"
	"context"
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("loading envs", err)
	}

	cfg, err := config.Read()
	if err != nil {
		slog.Error("reading config:", err)
		return
	}

	dbPool, err := repository.Connect(cfg)
	if err != nil {
		slog.Error("database connect", err)
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	redis, err := cache.Connect(cfg)
	defer func() {
		if redis != nil {
			err := redis.Close()
			if err != nil {
				slog.Error("closing redis", err)
				return
			}
		}
	}()
	if err != nil {
		slog.Error("redis connect", err)
	}

	storageURL := repository.NewStorageURL(dbPool)
	storageKey := repository.NewStorageKey(dbPool)

	cacheService := cache.NewCache(redis, cfg.Redis.Key)

	cleanerService := usecase.NewURLService(&storageURL, &storageKey, cacheService, cfg.MonthAmount)

	err = cleanerService.CleanURL(context.Background())

	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("creating scheduler", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_, err := s.NewJob(
			gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(1, 15),
				gocron.NewAtTimes(gocron.NewAtTime(1, 0, 0))),
			gocron.NewTask(
				func() {
					if ctx.Err() != nil {
						return
					}
					err := cleanerService.CleanURL(ctx)
					if err != nil {
						return
					}
				},
			),
		)
		if err != nil {
			slog.Error("creating job", err)
		}

		s.Start()
	}()

	slog.Info("cleaner service started")
	<-ctx.Done()

	cancel()

	err = s.Shutdown()
	if err != nil {
		slog.Error("shutdown scheduler", err)
	}
}
