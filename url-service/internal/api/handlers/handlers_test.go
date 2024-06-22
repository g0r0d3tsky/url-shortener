package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-service/internal/domain"

	mock_service "url-service/internal/api/handlers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateActorHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockURLService, url *domain.Url)
	testCases := []struct {
		name               string
		inputUrl           *domain.Url
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name: "OK",
			inputUrl: &domain.Url{
				OriginalURL: "http://example.com",
			},
			mockBehavior: func(r *mock_service.MockURLService, url *domain.Url) {
				r.EXPECT().CreateShortURL(gomock.Any(), url.OriginalURL).Return(url, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Bad Request",
			inputUrl: &domain.Url{
				OriginalURL: "http://example.com",
			},
			mockBehavior: func(r *mock_service.MockURLService, url *domain.Url) {
				r.EXPECT().CreateShortURL(gomock.Any(), url.OriginalURL).Return(nil, dummyError)
			},
			expectedStatusCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockURLService(c)
			tc.mockBehavior(service, tc.inputUrl)

			handler := NewAPIHandler(service)

			jsonData, err := json.Marshal(tc.inputUrl)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/data", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.CreateURL(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

		})
	}
}

func TestRedirectUrlHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockURLService, url *domain.Url)
	testCases := []struct {
		name               string
		inputUrl           *domain.Url
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name: "OK",
			inputUrl: &domain.Url{
				OriginalURL: "http://example.com",
			},
			mockBehavior: func(r *mock_service.MockURLService, url *domain.Url) {
				r.EXPECT().GetURL(gomock.Any(), "shorturl").Return(url, nil)
			},
			expectedStatusCode: 301,
		},
		{
			name: "Bad Request",
			inputUrl: &domain.Url{
				OriginalURL: "http://example.com",
			},
			mockBehavior: func(r *mock_service.MockURLService, url *domain.Url) {
				r.EXPECT().GetURL(gomock.Any(), "shorturl").Return(nil, dummyError)
			},
			expectedStatusCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockURLService(c)
			tc.mockBehavior(service, tc.inputUrl)

			handler := NewAPIHandler(service)

			req, err := http.NewRequest("GET", "/api/v1/data", nil)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.RedirectURL(recorder, req, "shorturl")

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

		})
	}
}
