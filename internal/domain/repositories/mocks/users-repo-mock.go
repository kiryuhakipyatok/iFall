package mocks

import (
	"context"
	"iFall/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FetchContacts(ctx context.Context) ([]models.Contacts, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Contacts), args.Error(1)
}

func (m *MockUserRepository) SetChatId(ctx context.Context, telegram string, chatId int64) error {
	args := m.Called(ctx, telegram, chatId)
	return args.Error(0)
}

func (m *MockUserRepository) DropChatId(ctx context.Context, telegram string, chatId int64) error {
	args := m.Called(ctx, telegram, chatId)
	return args.Error(0)
}
