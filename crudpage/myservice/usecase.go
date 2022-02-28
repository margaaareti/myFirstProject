package myservice

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type CrudUseCase struct {
	redisConn *redis.Client
}

func NewCrudUseCase(redisConn *redis.Client) *CrudUseCase {
	return &CrudUseCase{redisConn: redisConn}
}

func (c *CrudUseCase) LogOut(ctx context.Context, givenUuid string) (int64, error) {

	deleted, err := c.redisConn.Del(ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
