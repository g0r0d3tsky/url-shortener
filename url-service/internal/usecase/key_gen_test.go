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

func TestGetNextKeyFromSequence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyRepo := mock_repo.NewMockKeyRepo(ctrl)

	urlService := NewKeyGenService(mockKeyRepo)

	testCases := []struct {
		name     string
		key      *domain.Key
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Get next key successfully",
			key:  &domain.Key{},
			mockFunc: func() {
				var value uint64 = 42
				ptrUint64 := &value
				mockKeyRepo.EXPECT().GetNextKeyFromSequence(gomock.Any()).Return(ptrUint64, nil)
			},
			wantErr: false,
		},
		{
			name: "Get next key error",
			key:  &domain.Key{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().GetNextKeyFromSequence(gomock.Any()).Return(nil, fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			_, err := urlService.GetNextKeyFromSequence(context.Background())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateNewKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyRepo := mock_repo.NewMockKeyRepo(ctrl)

	urlService := NewKeyGenService(mockKeyRepo)

	testCases := []struct {
		name     string
		key      *domain.Key
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Create new key successfully",
			key:  &domain.Key{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().CreateNewKey(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Create new key error",
			key:  &domain.Key{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().CreateNewKey(gomock.Any(), gomock.Any()).Return(fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := urlService.CreateNewKey(context.Background(), tc.key)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetFreeKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyRepo := mock_repo.NewMockKeyRepo(ctrl)

	urlService := NewKeyGenService(mockKeyRepo)

	testCases := []struct {
		name     string
		key      *domain.Key
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Get free key successfully",
			key:  &domain.Key{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().GetFreeKey(gomock.Any()).Return(&domain.Key{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Get free key error",
			key:  &domain.Key{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().GetFreeKey(gomock.Any()).Return(nil, fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			_, err := urlService.GetFreeKey(context.Background())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyRepo := mock_repo.NewMockKeyRepo(ctrl)

	urlService := NewKeyGenService(mockKeyRepo)

	testCases := []struct {
		name     string
		key      *domain.Key
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Update key successfully",
			key:  &domain.Key{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().UpdateKey(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Update key error",
			key:  &domain.Key{},
			mockFunc: func() {
				mockKeyRepo.EXPECT().UpdateKey(gomock.Any(), gomock.Any()).Return(fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := urlService.UpdateKey(context.Background(), tc.key)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
