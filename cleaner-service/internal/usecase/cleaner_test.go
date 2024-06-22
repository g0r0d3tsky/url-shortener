package usecase

import (
	"context"
	"testing"

	mock_repo "cleaner-service/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCleanUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUrlRepo := mock_repo.NewMockUrlRepo(ctrl)
	mockKeyRepo := mock_repo.NewMockKeyRepo(ctrl)
	mockCache := mock_repo.NewMockCache(ctrl)

	urlService := NewURLService(mockUrlRepo, mockKeyRepo, mockCache, 1)

	testCases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Clean URL error",
			mockFunc: func() {
				mockUrlRepo.EXPECT().GetURL(gomock.Any(), 1).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := urlService.CleanURL(context.Background())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
