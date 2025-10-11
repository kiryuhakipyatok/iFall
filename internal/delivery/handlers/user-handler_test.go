package handlers

import (
	"bytes"
	"iFall/internal/config"
	mock_services "iFall/internal/domain/services/mocks"
	"iFall/internal/utils"
	"iFall/pkg/server"
	"iFall/pkg/validator"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Create(t *testing.T) {
	type mockBehavior = func(m *mock_services.MockUserService)
	tests := []struct {
		testName     string
		mockBehavior mockBehavior
		request      string
		expectedCode int
	}{
		{
			testName:     "success with telegram",
			request:      `{"name": "sanya", "email": "sanya@gmail.com", "telegram": "tg"}`,
			expectedCode: 200,
			mockBehavior: func(m *mock_services.MockUserService) {
				m.EXPECT().Create(gomock.Any(), "sanya", "sanya@gmail.com", utils.StrToPtr("tg")).Return(nil)
			},
		},
		{
			testName:     "success without telegram",
			request:      `{"name": "sanya", "email": "sanya@gmail.com"}`,
			expectedCode: 200,
			mockBehavior: func(m *mock_services.MockUserService) {
				m.EXPECT().Create(gomock.Any(), "sanya", "sanya@gmail.com", nil).Return(nil)
			},
		},
		{
			testName:     "failed without name",
			request:      `{"email": "sanya@gmail.com"}`,
			expectedCode: 400,
			mockBehavior: func(m *mock_services.MockUserService) {},
		},
		{
			testName:     "failed with invalid email",
			request:      `{"name": "sanya", "email": "sanyagmail.com"}`,
			expectedCode: 400,
			mockBehavior: func(m *mock_services.MockUserService) {},
		},
		{
			testName:     "failed with bad json",
			request:      `{"name": "sanya", "email": "sanya@gmail.com",}`,
			expectedCode: 422,
			mockBehavior: func(m *mock_services.MockUserService) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			validator := validator.NewValidator()
			mockService := mock_services.NewMockUserService(c)
			handler := NewUsersHandler(mockService, validator)
			a := server.NewServer(config.ServerConfig{}, config.AppConfig{})
			a.App.Post("/users", handler.CreateUser)
			tt.mockBehavior(mockService)
			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(tt.request))
			req.Header.Set("Content-Type", "application/json")
			resp, err := a.App.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
