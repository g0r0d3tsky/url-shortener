package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"url-service/config"
	"url-service/internal/api/handlers"
	handlers_gen "url-service/internal/api/handlers/gen"
	"url-service/internal/cache"
	"url-service/internal/kafka"
	request_metrics "url-service/internal/metrics"
	"url-service/internal/repository"
	"url-service/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("failed to load .env file", slog.String("msg", err.Error()))
	}
	cfg, err := config.Read()

	if err != nil {
		log.Println("failed to read config:", err.Error())
		return
	}
	dbPool, err := repository.Connect(cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	producer, err := kafka.New(cfg)
	if err != nil {
		slog.Error("kafka connect", err)
		return
	}

	defer func() {
		if producer != nil {
			err := producer.Close()
			if err != nil {
				slog.Error("closing producer", err)
				return
			}
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

	serviceURL := usecase.NewURLService(&storageURL, &storageKey, cacheService, producer, cfg.KafkaTopic)

	handlerURL := handlers.NewAPIHandler(serviceURL)

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("../../../swagger-ui/dist"))
	r.Handle("/swagger-ui/*", http.StripPrefix("/swagger-ui/", fs))

	reg := prometheus.NewRegistry()

	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	))
	reg.MustRegister(request_metrics.RequestMetrics)

	r.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	handlers_gen.HandlerFromMuxWithBaseURL(handlerURL, r, "/api/v1")

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: r,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on port %v...\n", cfg.Port)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}
