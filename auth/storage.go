package auth

import (
	"Test_derictory/models"

	"context"
)

type TokenStorage interface {
	CreateAuth(ctx context.Context, userid uint64, td *models.TokenDetails) error
	DeleteToken(ctx context.Context, givenUuid []string) (int64, error)
}
