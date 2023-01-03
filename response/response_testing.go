package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

var _ Sender = (*MockSender)(nil)

type MockSender struct {
	mock.Mock
}

func (m *MockSender) Send(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}
