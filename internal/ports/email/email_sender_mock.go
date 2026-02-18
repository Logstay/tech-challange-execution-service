package email

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type EmailSenderMock struct {
	mock.Mock
}

func (m *EmailSenderMock) SendEmail(ctx context.Context, e Email) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}
