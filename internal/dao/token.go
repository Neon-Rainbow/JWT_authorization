package dao

import (
	"JWT_authorization/util/Redis"
	"context"
	"time"
)

var ctx = context.Background()
var rdb = Redis.GetRedis()

// SetTokenToRedis is a function to set token to redis
func SetTokenToRedis(userID string, refreshToken string) error {
	rdb := Redis.GetRedis()
	err := rdb.Set(ctx, userID, refreshToken, time.Hour*24*15).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetTokenFromRedis is a function to get token from redis
func GetTokenFromRedis(userID string) (refreshToken string, err error) {
	refreshToken, err = rdb.Get(ctx, userID).Result()
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

// DeleteTokenFromRedis is a function to delete token from redis
func DeleteTokenFromRedis(userID string) error {
	err := rdb.Del(ctx, userID).Err()
	if err != nil {
		return err
	}
	return nil
}
