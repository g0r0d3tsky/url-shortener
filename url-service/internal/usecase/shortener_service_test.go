package usecase

import (
	"context"
	"fmt"
	"testing"
	"url-service/internal/domain"

	mock_repo "url-service/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUrlRepo := mock_repo.NewMockUrlRepo(ctrl)
	mockKeyRepo := mock_repo.NewMockKeyRepo(ctrl)
	mockCache := mock_repo.NewMockCache(ctrl)
	mockBrocker := mock_repo.NewMockBroker(ctrl)

	urlService := NewURLService(mockUrlRepo, mockKeyRepo, mockCache, mockBrocker, "test")

	testCases := []struct {
		name     string
		url      *domain.Url
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Get url from cache successfully",
			url:  &domain.Url{},
			mockFunc: func() {
				mockCache.EXPECT().GetURL(gomock.Any(), "a").Return(&domain.Url{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Get url from cache error",
			url:  &domain.Url{},
			mockFunc: func() {
				mockCache.EXPECT().GetURL(gomock.Any(), "a").Return(nil, fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
		{
			name: "Get url (cache nil) successfully",
			url:  &domain.Url{},
			mockFunc: func() {
				mockCache.EXPECT().GetURL(gomock.Any(), "a").Return(nil, nil)
				mockUrlRepo.EXPECT().GetURL(gomock.Any(), "a").Return(&domain.Url{}, nil)
				mockCache.EXPECT().SetURL(gomock.Any(), gomock.Any()).Return(nil)
				mockBrocker.EXPECT().Push(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Get url (cache nil) repo error",
			url:  &domain.Url{},
			mockFunc: func() {
				mockCache.EXPECT().GetURL(gomock.Any(), "a").Return(nil, nil)
				mockUrlRepo.EXPECT().GetURL(gomock.Any(), "a").Return(nil, fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
		{
			name: "Get url (cache nil) cache error",
			url:  &domain.Url{},
			mockFunc: func() {
				mockCache.EXPECT().GetURL(gomock.Any(), "a").Return(nil, nil)
				mockUrlRepo.EXPECT().GetURL(gomock.Any(), "a").Return(&domain.Url{}, nil)
				mockCache.EXPECT().SetURL(gomock.Any(), gomock.Any()).Return(fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
		{
			name: "Get url (cache nil) broker error",
			url:  &domain.Url{},
			mockFunc: func() {
				mockCache.EXPECT().GetURL(gomock.Any(), "a").Return(nil, nil)
				mockUrlRepo.EXPECT().GetURL(gomock.Any(), "a").Return(&domain.Url{}, nil)
				mockCache.EXPECT().SetURL(gomock.Any(), gomock.Any()).Return(nil)
				mockBrocker.EXPECT().Push(gomock.Any(), gomock.Any()).Return(fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			_, err := urlService.GetURL(context.Background(), "a")

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUrlRepo := mock_repo.NewMockUrlRepo(ctrl)
	mockKeyRepo := mock_repo.NewMockKeyRepo(ctrl)
	mockCache := mock_repo.NewMockCache(ctrl)
	mockBrocker := mock_repo.NewMockBroker(ctrl)

	urlService := NewURLService(mockUrlRepo, mockKeyRepo, mockCache, mockBrocker, "test")

	testCases := []struct {
		name     string
		url      *domain.Url
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Create short url successfully with free key",
			url:  &domain.Url{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().GetFreeKey(gomock.Any()).Return(&domain.Key{}, nil)
				mockBrocker.EXPECT().Push(gomock.Any(), gomock.Any()).Return(nil)
				mockKeyRepo.EXPECT().UpdateKey(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Get free key error",
			url:  &domain.Url{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().GetFreeKey(gomock.Any()).Return(nil, fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
		{
			name: "Create short url successfully without free key",
			url:  &domain.Url{},
			mockFunc: func() {
				var value uint64 = 42
				ptrUint64 := &value
				mockKeyRepo.EXPECT().GetFreeKey(gomock.Any()).Return(nil, nil)
				mockKeyRepo.EXPECT().GetNextKeyFromSequence(gomock.Any()).Return(ptrUint64, nil)
				mockBrocker.EXPECT().Push(gomock.Any(), gomock.Any()).Return(nil)
				mockKeyRepo.EXPECT().CreateNewKey(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			_, err := urlService.CreateShortURL(context.Background(), "a")

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
