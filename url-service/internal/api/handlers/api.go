package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"url-service/url-service/internal/api/handlers/model"
	"url-service/url-service/internal/domain"
)

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
	// Parse request body
	var input model.RequestInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create short URL
	shortURL, err := h.service.CreateShortURL(r.Context(), input.OriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Return short URL
	if err := json.NewEncoder(w).Encode(shortURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Redirect to actual URL
// (GET /{shortenedUrl})
func (h *APIHandler) RedirectURL(w http.ResponseWriter, r *http.Request, shortenedUrl string) {
	// Get actual URL
	url, err := h.service.GetURL(r.Context(), shortenedUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to actual URL
	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

//type APIHandler struct {
//	*MyHandler1
//	*MyHandler2
//}
