package mock

import (
	"Test_derictory/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (s *UserRepositoryMock) CreateUser(ctx context.Context, user models.User2) (int, error) {
	args := s.Called(ctx, user)

	return args.Int(0), args.Error(1)
}

func (s *UserRepositoryMock) GetUser(ctx context.Context, username, password string) (*models.User2, error) {
	args := s.Called(ctx, username, password)

	return args.Get(0).(*models.User2), args.Error(1)

}
