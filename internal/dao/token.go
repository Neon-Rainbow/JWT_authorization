package dao

import (
	"context"
	"time"
)

var ctx = context.Background()

// SetTokenToRedis is a function to set token to redis
func (dao *UserDAOImpl) SetTokenToRedis(userID string, refreshToken string) error {
	rdb := dao.Client
	ctx := dao.Context
	err := rdb.Set(ctx, userID, refreshToken, time.Hour*24*15).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetTokenFromRedis is a function to get token from redis
func (dao *UserDAOImpl) GetTokenFromRedis(userID string) (refreshToken string, err error) {
	rdb := dao.Client
	ctx := dao.Context
	refreshToken, err = rdb.Get(ctx, userID).Result()
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

// DeleteTokenFromRedis is a function to delete token from redis
func (dao *UserDAOImpl) DeleteTokenFromRedis(userID string) error {
	rdb := dao.Client
	ctx := dao.Context
	err := rdb.Del(ctx, userID).Err()
	if err != nil {
		return err
	}
	return nil
}
