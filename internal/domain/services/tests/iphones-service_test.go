package services

import (
	"context"
	mock_client "iFall/internal/client/mocks"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	mock_repositories "iFall/internal/domain/repositories/mocks"
	"iFall/internal/domain/services"
	mock_email "iFall/internal/email/mocks"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIPhoneService_Get(t *testing.T) {

	type mockBehavior = func(s *mock_repositories.MockIPhoneRepository, ctx context.Context, id string)

	tests := []struct {
		testName       string
		id             string
		mockBehavior   mockBehavior
		expectedError  error
		expectedResult *models.IPhone
	}{
		{
			testName:      "success receiving",
			id:            "00000000-0000-0000-0000-000000000000",
			expectedError: nil,
			expectedResult: &models.IPhone{
				Id:     "00000000-0000-0000-0000-000000000000",
				Name:   "iphone1",
				Price:  1000,
				Change: 0,
				Color:  "000000",
			},
			mockBehavior: func(s *mock_repositories.MockIPhoneRepository, ctx context.Context, id string) {
				s.EXPECT().Get(ctx, id).Return(&models.IPhone{
					Id:     "00000000-0000-0000-0000-000000000000",
					Name:   "iphone1",
					Price:  1000,
					Change: 0,
					Color:  "000000",
				}, nil)
			},
		},
		{
			testName:       "not found",
			id:             "00000000-0000-0000-0000-000000000000",
			expectedError:  errs.ErrNotFoundBase,
			expectedResult: nil,
			mockBehavior: func(s *mock_repositories.MockIPhoneRepository, ctx context.Context, id string) {
				s.EXPECT().Get(ctx, id).Return(nil, errs.ErrNotFoundBase)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			iphoneRepo := mock_repositories.NewMockIPhoneRepository(c)
			client := mock_client.NewMockApiClient(c)
			logger := logger.NewLogger(config.AppConfig{Name: "test", Env: "test", Version: "test", LogPath: "test.log"})
			emailSender := mock_email.NewMockEmailSender(c)
			ctx := context.Background()
			tt.mockBehavior(iphoneRepo, ctx, tt.id)
			iphoneService := services.NewIPhoneService(iphoneRepo, client, logger, emailSender, config.IPhonesConfig{})
			iphone, err := iphoneService.Get(ctx, tt.id)
			assert.Equal(t, tt.expectedResult, iphone)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}
