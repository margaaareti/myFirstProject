package auth

import (
	"Test_derictory/models"
	"context"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, user models.User2) (uint64, error)
	SignIn(ctx context.Context, username, password string) (*models.User2, error)
	CreateTokens(ctx context.Context, username string, userId uint64) (*models.TokenDetails, uint64, error)
	ParseAcsToken(ctx context.Context, accessToken string) (*models.AccessDetails, error)
	ParseRefToken(ctx context.Context, refreshToken string) (string, string, error)
	DeleteTokens(ctx context.Context, tokensUUID ...string) (uint64, error)
	CreateAuth(ctx context.Context, userId uint64, td *models.TokenDetails) error
	LogOut(ctx context.Context, givenUuid ...string) (uint64, error)
}
