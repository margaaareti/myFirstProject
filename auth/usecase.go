package auth

import (
	"Test_derictory/models"
	"context"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, user models.User2) (uint64, error)
	SignIn(ctx context.Context, username, password string) (*models.User2, error)
	CreateToken(ctx context.Context, username string, userId uint64) (*models.TokenDetails, uint64, error)
	ParseToken(ctx context.Context, accessToken string) (*models.AccessDetails, error)
	CreateAuth(ctx context.Context, userId uint64, td *models.TokenDetails) error
	LogOut(ctx context.Context, givenUuid string) (int64, error)
}
