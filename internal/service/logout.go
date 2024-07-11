package service

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
)

// ProcessLogoutRequest deletes the token from Redis
func (s *UserServiceImpl) ProcessLogoutRequest(userID string) *model.ApiError {
	err := s.DeleteTokenFromRedis(userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.DeleteUserTokenError,
			Message:      "delete token from redis error",
			ErrorMessage: err,
		}
	}
	return nil
}
