package redis

import (
	"context"

	"github.com/Adityadangi14/ecomm_ai/config"
	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(config *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.RedisAddr,
		Password: config.Redis.RedisPassword,
		DB:       config.Redis.DB,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil

}
