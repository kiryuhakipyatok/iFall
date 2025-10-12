package services

import (
	"context"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	mock_repositories "iFall/internal/domain/repositories/mocks"
	"iFall/internal/utils"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create(t *testing.T) {

	type ttData struct {
		name          string
		email         string
		telegram      *string
		expectedError error
	}

	type mockBehavior = func(s *mock_repositories.MockUserRepository, ctx context.Context, tt ttData)

	tests := []struct {
		testName string
		ttData
		mockBehavior mockBehavior
	}{
		{
			testName: "success create without telegram",
			ttData: ttData{
				name:          "sanya",
				email:         "sanyaemail@gmail.com",
				telegram:      nil,
				expectedError: nil,
			},

			mockBehavior: func(s *mock_repositories.MockUserRepository, ctx context.Context, tt ttData) {
				s.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *models.User) error {
					assert.Equal(t, tt.name, user.Name)
					assert.Equal(t, tt.email, user.Email)
					assert.Nil(t, user.Telegram)
					assert.NotEqual(t, uuid.Nil, user.Id)
					return tt.expectedError
				})
			},
		},
		{
			testName: "success create with telegram",
			ttData: ttData{
				name:          "kir",
				email:         "kiremail@gmail.com",
				telegram:      utils.StrToPtr("kirtg"),
				expectedError: nil,
			},
			mockBehavior: func(s *mock_repositories.MockUserRepository, ctx context.Context, tt ttData) {
				s.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *models.User) error {
					assert.Equal(t, tt.name, user.Name)
					assert.Equal(t, tt.email, user.Email)
					assert.Equal(t, tt.telegram, user.Telegram)
					assert.NotEqual(t, uuid.Nil, user.Id)
					return tt.expectedError
				})
			},
		},
		{
			testName: "already exists with telegram",
			ttData: ttData{
				name:          "kir",
				email:         "kiremail@gmail.com",
				telegram:      utils.StrToPtr("kirtg"),
				expectedError: errs.ErrAlreadyExistsBase,
			},
			mockBehavior: func(s *mock_repositories.MockUserRepository, ctx context.Context, tt ttData) {
				s.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, user *models.User) error {
					assert.Equal(t, tt.name, user.Name)
					assert.Equal(t, tt.email, user.Email)
					assert.Equal(t, tt.telegram, user.Telegram)
					assert.NotEqual(t, uuid.Nil, user.Id)
					return tt.expectedError
				})
			},
		},
		{
			testName: "already exists without telegram",
			ttData: ttData{
				name:          "sanya",
				email:         "sanyaemail@gmail.com",
				telegram:      nil,
				expectedError: errs.ErrAlreadyExistsBase,
			},
			mockBehavior: func(s *mock_repositories.MockUserRepository, ctx context.Context, tt ttData) {
				s.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *models.User) error {
					assert.Equal(t, tt.name, user.Name)
					assert.Equal(t, tt.email, user.Email)
					assert.Nil(t, user.Telegram)
					assert.NotEqual(t, uuid.Nil, user.Id)
					return tt.expectedError
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			logger := logger.NewLogger(config.AppConfig{Name: "test", Env: "test", Version: "test", LogPath: ""})
			mockUserRepo := mock_repositories.NewMockUserRepository(c)
			userService := NewUserService(mockUserRepo, logger)
			ctx := context.Background()

			tt.mockBehavior(mockUserRepo, ctx, tt.ttData)
			err := userService.Create(context.Background(), tt.name, tt.email, tt.telegram)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}
