package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"storage-service/internal/config"
	"storage-service/internal/handlers"
	"storage-service/internal/kafka"
	"storage-service/internal/repository"
	"storage-service/internal/usecase"
	"sync"
	"syscall"

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

	repo := repository.NewStorageURL(dbPool)
	service := usecase.NewService(&repo)
	handler := handlers.NewServiceMessage(service)

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	wg.Add(1)
	consumer, err := kafka.RunConsumer(ctx, wg, cfg, handler.Handler())
	if err != nil {
		slog.Error("kafka consumer", err)
	}
	slog.Info("kafka consumer started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			slog.Info("terminate: context done")
			keepRunning = false
		case <-ch:
			slog.Info("terminate: signal")
			keepRunning = false
		}
	}

	cancel()

	wg.Wait()

	if err = consumer.Close(); err != nil {
		slog.Error("closing consumer", err)
		return
	}
}
