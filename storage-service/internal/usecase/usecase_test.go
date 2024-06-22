package usecase

import (
	"context"
	"fmt"
	"storage-service/internal/domain"
	"testing"

	mock_repo "storage-service/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUrlRepo := mock_repo.NewMockRepository(ctrl)
	service := NewService(mockUrlRepo)

	testCases := []struct {
		name     string
		url      *domain.Url
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Create URL successfully",
			url:  &domain.Url{},
			mockFunc: func() {
				mockUrlRepo.EXPECT().CheckExist(gomock.Any(), &domain.Url{}).Return(false, nil)
				mockUrlRepo.EXPECT().CreateShort(gomock.Any(), &domain.Url{}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Create URL exist",
			url:  &domain.Url{},
			mockFunc: func() {
				mockUrlRepo.EXPECT().CheckExist(gomock.Any(), &domain.Url{}).Return(true, nil)
				mockUrlRepo.EXPECT().UpdateURL(gomock.Any(), &domain.Url{}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Create URL exist error",
			url:  &domain.Url{},
			mockFunc: func() {
				mockUrlRepo.EXPECT().CheckExist(gomock.Any(), &domain.Url{}).Return(true, nil)
				mockUrlRepo.EXPECT().UpdateURL(gomock.Any(), &domain.Url{}).Return(fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
		{
			name: "Create URL exist error",
			url:  &domain.Url{},
			mockFunc: func() {
				mockUrlRepo.EXPECT().CheckExist(gomock.Any(), &domain.Url{}).Return(false, nil)
				mockUrlRepo.EXPECT().CreateShort(gomock.Any(), &domain.Url{}).Return(fmt.Errorf("dummy error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := service.CreateURL(context.Background(), tc.url)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
