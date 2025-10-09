package client

import (
	"iFall/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetIPhoneData(id string) (*models.IPhone, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IPhone), args.Error(1)
}
