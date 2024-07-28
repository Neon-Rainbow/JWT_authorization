package service

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
)

// ProcessLogoutRequest deletes the token from Redis
func (s *UserServiceImpl) ProcessLogoutRequest(ctx context.Context, userID string) *model.ApiError {
	err := s.DeleteTokenFromRedis(ctx, userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.DeleteUserTokenError,
			Message:      "delete token from redis error",
			ErrorMessage: err,
		}
	}
	return nil
}
