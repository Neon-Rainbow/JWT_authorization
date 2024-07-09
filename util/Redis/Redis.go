package Redis

import (
	"JWT_authorization/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis() error {
	cfg := config.GetConfig()

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	// Test the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetRedis() *redis.Client {
	if rdb == nil {
		return nil
	}
	return rdb
}
