package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
)

// ProcessLogoutRequest deletes the token from Redis
func ProcessLogoutRequest(userID string) *model.ApiError {
	err := dao.DeleteTokenFromRedis(userID)
	if err != nil {
		return &model.ApiError{
			Code:    code.DeleteUserTokenError,
			Message: "delete token from redis error",
			Error:   err,
		}
	}
	return nil
}
