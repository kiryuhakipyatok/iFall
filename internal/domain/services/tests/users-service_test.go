package services

import (
	"context"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	repoMocks "iFall/internal/domain/repositories/mocks"
	"iFall/internal/domain/services"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		testName string
		name     string
		email    string
		telegram *string
		isError  bool
		result   error
	}{
		{
			testName: "success create without telegram",
			name:     "sanya",
			email:    "sanyaemail@gmail.com",
			telegram: nil,
			isError:  false,
			result:   nil,
		},
		{
			testName: "success create with telegram",
			name:     "kir",
			email:    "kiremail@gmail.com",
			telegram: toPtr("kirtg"),
			isError:  false,
			result:   nil,
		},
		{
			testName: "already exists",
			name:     "kir",
			email:    "kiremail@gmail.com",
			telegram: toPtr("kirtg"),
			isError:  true,
			result:   errs.ErrAlreadyExistsBase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockUserRepo := new(repoMocks.MockUserRepository)
			logger := logger.NewLogger(config.AppConfig{Name: "test", Version: "1.0.0", Env: "test", LogPath: "test.log"})
			if tt.isError {
				mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(tt.result)
			} else {
				mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			}
			userService := services.NewUserService(mockUserRepo, logger)
			err := userService.Create(context.Background(), tt.name, tt.email, tt.telegram)
			if tt.isError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.result)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, err)
			}
			mockUserRepo.AssertCalled(t, "Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
				return u.Name == tt.name && u.Email == tt.email && ((u.Telegram == nil && tt.telegram == nil) || u.Telegram != nil && tt.telegram != nil && *u.Telegram == *tt.telegram)
			}))
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func toPtr(s string) *string {
	return &s
}
