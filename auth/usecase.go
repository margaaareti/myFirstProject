package auth

import (
	"Test_derictory/models"
	"context"
	"io"
)

const CtxUserId = "userId"
const CtxUserName = "name"
const CtxUserSurname = "surname"
const CtxUserPatronymic = "patronymic"

type UseCase interface {
	SignUp(ctx context.Context, user models.User2) (uint64, error)
	SignIn(ctx context.Context, username, password string) (*models.User2, error)
	CreateTokens(ctx context.Context, user *models.User2) (*models.TokenDetails, *models.User2, error)
	ParseAcsToken(ctx context.Context, accessToken string) (*models.AccessDetails, error)
	ParseRefToken(ctx context.Context, refreshToken string) (string, *models.User2, error)
	DeleteTokens(ctx context.Context, tokensUUID ...string) (uint64, error)
	CreateAuth(ctx context.Context, userId uint64, td *models.TokenDetails) error
	LogOut(ctx context.Context, givenUuid ...string) (uint64, error)
	UploadImage(ctx context.Context, file io.Reader, size int64, contentType string)
}
