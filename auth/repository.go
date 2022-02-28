package auth

import (
	"Test_derictory/models"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User2) (uint64, error)
	GetUser(ctx context.Context, username, password string) (*models.User2, error)
}
