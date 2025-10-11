package services

import (
	"context"
	mock_client "iFall/internal/client/mocks"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	mock_repositories "iFall/internal/domain/repositories/mocks"
	mock_email "iFall/internal/email/mocks"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIPhoneService_Get(t *testing.T) {

	type mockBehavior = func(m *mock_repositories.MockIPhoneRepository, ctx context.Context, id string)

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
			mockBehavior: func(m *mock_repositories.MockIPhoneRepository, ctx context.Context, id string) {
				m.EXPECT().Get(ctx, id).Return(&models.IPhone{
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
			mockBehavior: func(m *mock_repositories.MockIPhoneRepository, ctx context.Context, id string) {
				m.EXPECT().Get(ctx, id).Return(nil, errs.ErrNotFoundBase)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			iphoneRepo := mock_repositories.NewMockIPhoneRepository(c)
			client := mock_client.NewMockApiClient(c)
			logger := logger.NewLogger(config.AppConfig{Name: "test", Env: "test", Version: "test", LogPath: ""})
			emailSender := mock_email.NewMockEmailSender(c)
			ctx := context.Background()
			tt.mockBehavior(iphoneRepo, ctx, tt.id)
			iphoneService := NewIPhoneService(iphoneRepo, client, logger, emailSender, config.IPhonesConfig{})
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

func TestIphoneService_Update(t *testing.T) {
	type ttData struct {
		id             string
		expectedResult *models.IPhone
		expectedError  error
	}
	type mockBehavior = func(mr *mock_repositories.MockIPhoneRepository, mc *mock_client.MockApiClient, ctx context.Context, ttData ttData)
	tests := []struct {
		testName     string
		ttData       ttData
		mockBehavior mockBehavior
	}{
		{
			testName: "success update",
			ttData: ttData{
				id: "iphone1-id",
				expectedResult: &models.IPhone{
					Id:     "iphone1-id",
					Name:   "iphone1",
					Price:  900.0,
					Change: 100.0,
					Color:  "ffffff",
				},
				expectedError: nil,
			},

			mockBehavior: func(mr *mock_repositories.MockIPhoneRepository, mc *mock_client.MockApiClient, ctx context.Context, ttData ttData) {
				gomock.InOrder(
					mc.EXPECT().GetIPhoneData(ttData.id).Return(&models.IPhone{
						Id:    "iphone1-id",
						Name:  "iphone1",
						Price: 900.0,
						Color: "ffffff",
					}, nil),
					mr.EXPECT().Update(ctx, ttData.id, 900.0).Return(&models.IPhone{
						Id:     "iphone1-id",
						Name:   "iphone1",
						Price:  900.0,
						Change: 100.0,
						Color:  "ffffff",
					}, ttData.expectedError),
				)
			},
		},
		{
			testName: "not found",
			ttData: ttData{
				id:             "iphone1-id",
				expectedResult: nil,
				expectedError:  errs.ErrNotFoundBase,
			},

			mockBehavior: func(mr *mock_repositories.MockIPhoneRepository, mc *mock_client.MockApiClient, ctx context.Context, ttData ttData) {
				gomock.InOrder(
					mc.EXPECT().GetIPhoneData(ttData.id).Return(&models.IPhone{
						Id:    "iphone1-id",
						Name:  "iphone1",
						Price: 900.0,
						Color: "ffffff",
					}, nil),
					mr.EXPECT().Update(ctx, ttData.id, 900.0).Return(nil, errs.ErrNotFoundBase),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			mockClient := mock_client.NewMockApiClient(c)
			mockRepository := mock_repositories.NewMockIPhoneRepository(c)
			logger := logger.NewLogger(config.AppConfig{Name: "test", Env: "test", Version: "test", LogPath: ""})
			emailSender := mock_email.NewMockEmailSender(c)
			ctx := context.Background()
			tt.mockBehavior(mockRepository, mockClient, ctx, tt.ttData)
			service := NewIPhoneService(mockRepository, mockClient, logger, emailSender, config.IPhonesConfig{})
			iphone, err := service.Update(ctx, tt.ttData.id)
			if tt.ttData.expectedError != nil {
				assert.ErrorIs(t, err, tt.ttData.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.ttData.expectedResult, iphone)
			}

		})
	}
}

func TestIphoneService_UpdateAll(t *testing.T) {
	type ttData struct {
		expectedResult []models.IPhone
		expectedError  error
	}
	type mockBehavior = func(mr *mock_repositories.MockIPhoneRepository, mc *mock_client.MockApiClient, ctx context.Context, ttData ttData)
	tests := []struct {
		testName     string
		ttData       ttData
		mockBehavior mockBehavior
	}{
		{
			testName: "success updating",
			ttData: ttData{
				expectedResult: []models.IPhone{
					{
						Id:     "iphone-black-id",
						Name:   "iphone-black-name",
						Price:  900.0,
						Change: 0.0,
						Color:  "black",
					},
					{
						Id:     "iphone-white-id",
						Name:   "iphone-white-name",
						Price:  920.0,
						Change: 20.0,
						Color:  "white",
					},
					{
						Id:     "iphone-blue-id",
						Name:   "iphone-blue-name",
						Price:  1000,
						Change: 100,
						Color:  "blue",
					},
				},
				expectedError: nil,
			},
			mockBehavior: func(mr *mock_repositories.MockIPhoneRepository, mc *mock_client.MockApiClient, ctx context.Context, ttData ttData) {

				mc.EXPECT().GetIPhoneData("iphone-black-id").Return(&models.IPhone{
					Id:    "iphone-black-id",
					Name:  "iphone-black-name",
					Price: 900.0,
					Color: "black",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-black-id", 900.0).Return(&models.IPhone{
					Id:     "iphone-black-id",
					Name:   "iphone-black-name",
					Price:  900.0,
					Change: 0.0,
					Color:  "black",
				}, nil)

				mc.EXPECT().GetIPhoneData("iphone-white-id").Return(&models.IPhone{
					Id:    "iphone-white-id",
					Name:  "iphone-white-name",
					Price: 920.0,
					Color: "white",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-white-id", 920.0).Return(&models.IPhone{
					Id:     "iphone-white-id",
					Name:   "iphone-white-name",
					Price:  920.0,
					Change: 20.0,
					Color:  "white",
				}, nil)

				mc.EXPECT().GetIPhoneData("iphone-blue-id").Return(&models.IPhone{
					Id:    "iphone-blue-id",
					Name:  "iphone-blue-name",
					Price: 1000.0,
					Color: "blue",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-blue-id", 1000.0).Return(&models.IPhone{
					Id:     "iphone-blue-id",
					Name:   "iphone-blue-name",
					Price:  1000.0,
					Change: 100.0,
					Color:  "blue",
				}, nil)

			},
		},
		{
			testName: "all not found",
			ttData: ttData{
				expectedResult: nil,
				expectedError:  errs.ErrNotFoundBase,
			},
			mockBehavior: func(mr *mock_repositories.MockIPhoneRepository, mc *mock_client.MockApiClient, ctx context.Context, ttData ttData) {

				mc.EXPECT().GetIPhoneData("iphone-black-id").Return(&models.IPhone{
					Id:    "iphone-black-id",
					Name:  "iphone-black-name",
					Price: 900.0,
					Color: "black",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-black-id", 900.0).Return(&models.IPhone{
					Id:     "iphone-black-id",
					Name:   "iphone-black-name",
					Price:  900.0,
					Change: 0.0,
					Color:  "black",
				}, errs.ErrNotFoundBase)

				mc.EXPECT().GetIPhoneData("iphone-white-id").Return(&models.IPhone{
					Id:    "iphone-white-id",
					Name:  "iphone-white-name",
					Price: 920.0,
					Color: "white",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-white-id", 920.0).Return(&models.IPhone{
					Id:     "iphone-white-id",
					Name:   "iphone-white-name",
					Price:  920.0,
					Change: 20.0,
					Color:  "white",
				}, errs.ErrNotFoundBase)

				mc.EXPECT().GetIPhoneData("iphone-blue-id").Return(&models.IPhone{
					Id:    "iphone-blue-id",
					Name:  "iphone-blue-name",
					Price: 1000.0,
					Color: "blue",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-blue-id", 1000.0).Return(&models.IPhone{
					Id:     "iphone-blue-id",
					Name:   "iphone-blue-name",
					Price:  1000.0,
					Change: 100.0,
					Color:  "blue",
				}, errs.ErrNotFoundBase)

			},
		},
		{
			testName: "one not found",
			ttData: ttData{
				expectedResult: nil,
				expectedError:  errs.ErrNotFoundBase,
			},
			mockBehavior: func(mr *mock_repositories.MockIPhoneRepository, mc *mock_client.MockApiClient, ctx context.Context, ttData ttData) {

				mc.EXPECT().GetIPhoneData("iphone-black-id").Return(&models.IPhone{
					Id:    "iphone-black-id",
					Name:  "iphone-black-name",
					Price: 900.0,
					Color: "black",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-black-id", 900.0).Return(&models.IPhone{
					Id:     "iphone-black-id",
					Name:   "iphone-black-name",
					Price:  900.0,
					Change: 0.0,
					Color:  "black",
				}, errs.ErrNotFoundBase)

				mc.EXPECT().GetIPhoneData("iphone-white-id").Return(&models.IPhone{
					Id:    "iphone-white-id",
					Name:  "iphone-white-name",
					Price: 920.0,
					Color: "white",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-white-id", 920.0).Return(&models.IPhone{
					Id:     "iphone-white-id",
					Name:   "iphone-white-name",
					Price:  920.0,
					Change: 20.0,
					Color:  "white",
				}, nil)

				mc.EXPECT().GetIPhoneData("iphone-blue-id").Return(&models.IPhone{
					Id:    "iphone-blue-id",
					Name:  "iphone-blue-name",
					Price: 1000.0,
					Color: "blue",
				}, nil)
				mr.EXPECT().Update(gomock.Any(), "iphone-blue-id", 1000.0).Return(&models.IPhone{
					Id:     "iphone-blue-id",
					Name:   "iphone-blue-name",
					Price:  1000.0,
					Change: 100.0,
					Color:  "blue",
				}, nil)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repoMock := mock_repositories.NewMockIPhoneRepository(c)
			clientMock := mock_client.NewMockApiClient(c)
			logger := logger.NewLogger(config.AppConfig{Name: "test", Env: "test", Version: "test", LogPath: ""})
			emailMock := mock_email.NewMockEmailSender(c)
			cfg := config.IPhonesConfig{
				Black: "iphone-black-id",
				White: "iphone-white-id",
				Blue:  "iphone-blue-id",
			}
			ctx := context.Background()
			service := NewIPhoneService(repoMock, clientMock, logger, emailMock, cfg)
			tt.mockBehavior(repoMock, clientMock, ctx, tt.ttData)
			iphones, err := service.UpdateAll()
			if tt.ttData.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.ttData.expectedError)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.ttData.expectedResult, iphones)
			}
		})
	}
}
