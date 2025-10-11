package services

import (
	"errors"
	mock_bot "iFall/internal/bot/mocks"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	mock_repositories "iFall/internal/domain/repositories/mocks"
	mock_email "iFall/internal/email/mocks"
	"iFall/internal/utils"

	"iFall/pkg/logger"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIphoneReportService_Test(t *testing.T) {
	type mockBehavior = func(um *mock_repositories.MockUserRepository, em *mock_email.MockEmailSender, bm *mock_bot.MockTelegramBot)
	sendingError := errors.New("sending error")
	type ttData struct {
		iphones []models.IPhone

		expectedError error
	}
	tests := []struct {
		testName     string
		ttData       ttData
		mockBehavior mockBehavior
	}{
		{
			testName: "success reporting",
			ttData: ttData{
				iphones: []models.IPhone{
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
			mockBehavior: func(um *mock_repositories.MockUserRepository, em *mock_email.MockEmailSender, bm *mock_bot.MockTelegramBot) {
				um.EXPECT().FetchContacts(gomock.Any()).Return([]models.Contacts{
					{
						Telegram: utils.StrToPtr("tg1"),
						Email:    "kiremail@gmail.com",
						ChatId:   utils.Int64ToPtr(000000),
					},
					{
						Telegram: nil,
						Email:    "gusemail$gmail.com",
						ChatId:   nil,
					},
				}, nil)
				em.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(nil)
				bm.EXPECT().SendIPhonesInfo(gomock.Any(), []models.IPhone{
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
				}).Return(nil)
			},
		},
		{
			testName: "zero contacts",
			ttData: ttData{
				iphones: []models.IPhone{
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
			mockBehavior: func(um *mock_repositories.MockUserRepository, em *mock_email.MockEmailSender, bm *mock_bot.MockTelegramBot) {
				um.EXPECT().FetchContacts(gomock.Any()).Return([]models.Contacts{}, nil)
			},
		},
		{
			testName: "only emails",
			ttData: ttData{
				iphones: []models.IPhone{
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
			mockBehavior: func(um *mock_repositories.MockUserRepository, em *mock_email.MockEmailSender, bm *mock_bot.MockTelegramBot) {
				um.EXPECT().FetchContacts(gomock.Any()).Return([]models.Contacts{
					{
						Telegram: nil,
						Email:    "gusemail$gmail.com",
						ChatId:   nil,
					},
				}, nil)
				em.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(nil)
			},
		},
		{
			testName: "failed to send",
			ttData: ttData{
				iphones: []models.IPhone{
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
				expectedError: sendingError,
			},
			mockBehavior: func(um *mock_repositories.MockUserRepository, em *mock_email.MockEmailSender, bm *mock_bot.MockTelegramBot) {
				um.EXPECT().FetchContacts(gomock.Any()).Return([]models.Contacts{
					{
						Telegram: utils.StrToPtr("tg1"),
						Email:    "kiremail@gmail.com",
						ChatId:   utils.Int64ToPtr(000000),
					},
					{
						Telegram: nil,
						Email:    "gusemail$gmail.com",
						ChatId:   nil,
					},
				}, nil)
				em.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(sendingError)
				bm.EXPECT().SendIPhonesInfo(gomock.Any(), []models.IPhone{
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
				}).Return(sendingError)
			},
		},
		{
			testName: "failed to send emails",
			ttData: ttData{
				iphones: []models.IPhone{
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
				expectedError: sendingError,
			},
			mockBehavior: func(um *mock_repositories.MockUserRepository, em *mock_email.MockEmailSender, bm *mock_bot.MockTelegramBot) {
				um.EXPECT().FetchContacts(gomock.Any()).Return([]models.Contacts{
					{
						Telegram: utils.StrToPtr("tg1"),
						Email:    "kiremail@gmail.com",
						ChatId:   utils.Int64ToPtr(000000),
					},
					{
						Telegram: nil,
						Email:    "gusemail$gmail.com",
						ChatId:   nil,
					},
				}, nil)
				em.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(sendingError)
				bm.EXPECT().SendIPhonesInfo(gomock.Any(), []models.IPhone{
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
				}).Return(nil)
			},
		},
		{
			testName: "failed to send telegram",
			ttData: ttData{
				iphones: []models.IPhone{
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
				expectedError: sendingError,
			},
			mockBehavior: func(um *mock_repositories.MockUserRepository, em *mock_email.MockEmailSender, bm *mock_bot.MockTelegramBot) {
				um.EXPECT().FetchContacts(gomock.Any()).Return([]models.Contacts{
					{
						Telegram: utils.StrToPtr("tg1"),
						Email:    "kiremail@gmail.com",
						ChatId:   utils.Int64ToPtr(000000),
					},
					{
						Telegram: nil,
						Email:    "gusemail$gmail.com",
						ChatId:   nil,
					},
				}, nil)
				em.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(nil)
				bm.EXPECT().SendIPhonesInfo(gomock.Any(), []models.IPhone{
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
				}).Return(sendingError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			userMockRepo := mock_repositories.NewMockUserRepository(c)
			emailMock := mock_email.NewMockEmailSender(c)
			botMock := mock_bot.NewMockTelegramBot(c)
			logger := logger.NewLogger(config.AppConfig{Name: "test", Env: "test", Version: "test", LogPath: ""})
			service := NewIPhoneReportService(userMockRepo, logger, botMock, emailMock, config.IPhonesConfig{Timeout: time.Second})
			tt.mockBehavior(userMockRepo, emailMock, botMock)
			err := service.SendIPhonesInfo(tt.ttData.iphones)
			if tt.ttData.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.ttData.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
