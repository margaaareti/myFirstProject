package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisCon() (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //no password set
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {

		return nil, err
	}

	rdb.WithTimeout(10 * time.Second)

	return rdb, nil
}
