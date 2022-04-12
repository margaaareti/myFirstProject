package mainpage

import "context"

type HomeRepo interface {
	LogOut(ctx context.Context, givenUuid ...string) (int64, error)
}
