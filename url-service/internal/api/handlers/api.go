package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"url-service/internal/api/handlers/model"
	"url-service/internal/domain"
	"url-service/internal/metrics"
)

//go:generate mockgen -source=api.go -destination=mocks/handlers_mock.go
type URLService interface {
	CreateShortURL(ctx context.Context, originalUrl string) (*domain.Url, error)
	GetURL(ctx context.Context, short string) (*domain.Url, error)
}

type APIHandler struct {
	service URLService
}

func NewAPIHandler(service URLService) *APIHandler {
	return &APIHandler{service: service}
}

// (POST /data/shorten)
func (h *APIHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Parse request body
	var input model.RequestInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		duration := time.Since(start)
		statusCode := http.StatusBadRequest
		metrics.ObserveRequest(duration, statusCode, r.URL.Path)
		return
	}

	// Create short URL
	shortURL, err := h.service.CreateShortURL(r.Context(), input.OriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		duration := time.Since(start)
		statusCode := http.StatusInternalServerError
		metrics.ObserveRequest(duration, statusCode, r.URL.Path)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(shortURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		duration := time.Since(start)
		statusCode := http.StatusInternalServerError
		metrics.ObserveRequest(duration, statusCode, r.URL.Path)

		return
	}

	duration := time.Since(start)
	statusCode := http.StatusOK
	metrics.ObserveRequest(duration, statusCode, r.URL.Path)
}

// Redirect to actual URL
// (GET /{shortenedUrl})
func (h *APIHandler) RedirectURL(w http.ResponseWriter, r *http.Request, shortenedUrl string) {
	// Get actual URL
	start := time.Now()

	url, err := h.service.GetURL(r.Context(), shortenedUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		duration := time.Since(start)
		statusCode := http.StatusInternalServerError
		metrics.ObserveRequest(duration, statusCode, r.URL.Path)
		return
	}

	// Redirect to actual URL
	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
	duration := time.Since(start)
	statusCode := http.StatusMovedPermanently
	metrics.ObserveRequest(duration, statusCode, r.URL.Path)
}
