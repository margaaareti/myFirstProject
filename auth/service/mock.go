package service

import (
	"Test_derictory/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) SignUp(ctx context.Context, user models.User) (int, error) {
	args := m.Called(user)

	return 0, args.Error(0)
}
