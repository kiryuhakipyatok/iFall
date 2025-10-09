package mocks

import (
	"context"
	"iFall/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockIPhoneRepository struct {
	mock.Mock
}

func (m *MockIPhoneRepository) Get(ctx context.Context, id string) (*models.IPhone, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IPhone), args.Error(1)
}

func (m *MockIPhoneRepository) Update(ctx context.Context, id string, price float64) (*models.IPhone, error) {
	args := m.Called(ctx, id, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IPhone), args.Error(1)
}
