package email

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendMessage(ctx context.Context, sub string, content []byte, to []string, attachFiles []string) error {
	args := m.Called(ctx, sub, content, to, attachFiles)
	return args.Error(0)
}
