package services

import (
	"context"
	"iFall/internal/client"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	repoMocks "iFall/internal/domain/repositories/mocks"
	"iFall/internal/domain/services"
	"iFall/internal/email"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetIPhone(t *testing.T) {
	tests := []struct {
		testName string
		id       string
		isError  bool
		resError error
		result   *models.IPhone
	}{
		{
			testName: "success receiving",
			id:       "00000000-0000-0000-0000-000000000000",
			isError:  false,
			resError: nil,
			result: &models.IPhone{
				Id:     "00000000-0000-0000-0000-000000000000",
				Name:   "iphone1",
				Price:  1000,
				Change: 0,
				Color:  "000000",
			},
		},
		{
			testName: "not found",
			id:       "00000000-0000-0000-0000-000000000000",
			isError:  true,
			resError: errs.ErrNotFoundBase,
			result:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockIPhoneRepo := new(repoMocks.MockIPhoneRepository)
			mockClient := new(client.MockClient)
			mockEmailSender := new(email.MockEmailSender)
			logger := logger.NewLogger(config.AppConfig{Name: "test", Version: "1.0.0", Env: "test", LogPath: "test.log"})
			if tt.isError {
				mockIPhoneRepo.On("Get", mock.Anything, tt.id).Return(nil, tt.resError)
			} else {
				mockIPhoneRepo.On("Get", mock.Anything, tt.id).Return(tt.result, nil)
			}
			iphoneService := services.NewIPhoneService(mockIPhoneRepo, mockClient, logger, mockEmailSender, config.IPhonesConfig{})
			iphone, err := iphoneService.Get(context.Background(), tt.id)
			if tt.isError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.resError)
				assert.Nil(t, iphone)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.result, iphone)
			}
			mockIPhoneRepo.AssertExpectations(t)
		})
	}
}

// func TestUpdateAll(t *testing.T) {
// 	tests := []struct {
// 		testName string
// 		isError  bool
// 		result   error
// 	}{
// 		{
// 			testName: "successfully update all",
// 			isError:  false,
// 			result:   nil,
// 		},
// 	}
// }
