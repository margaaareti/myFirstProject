package redst

import (
	"Test_derictory/models"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type AuthRedis struct {
	storage *redis.Client
}

func NewAuthRedis(storage *redis.Client) *AuthRedis {
	return &AuthRedis{storage: storage}
}

func (a *AuthRedis) CreateAuth(ctx context.Context, userid uint64, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // convert Unix to UTC
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := a.storage.Set(ctx, td.AccessUuid.String(), strconv.FormatUint(userid, 10), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := a.storage.Set(ctx, td.RefreshUuid.String(), strconv.FormatUint(userid, 10), rt.Sub(now)).Err()
	if errAccess != nil {
		return errRefresh
	}

	return nil

}

func (a *AuthRedis) DeleteToken(ctx context.Context, givenUuid []string) (int64, error) {

	if len(givenUuid) == 2 {
		aToken := givenUuid[0]
		rToken := givenUuid[1]

		deleted, err := a.storage.Del(ctx, aToken, rToken).Result()
		if err != nil {
			return 0, err
		}
		return deleted, nil

	} else {
		rToken := givenUuid[0]

		deleted, err := a.storage.Del(ctx, rToken).Result()
		if err != nil {
			return 0, err
		}
		return deleted, nil

	}
}
