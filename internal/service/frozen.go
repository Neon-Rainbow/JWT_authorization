package service

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
)

func (s *UserServiceImpl) ProcessFreezeUser(ctx context.Context, userID string) (apiError *model.ApiError) {
	err := s.FreezeUser(ctx, userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.FrozenUserError,
			Message:      "frozen user to MySQL error",
			ErrorMessage: err,
		}
	}
	err = s.DeleteTokenFromRedis(ctx, userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.FrozenUserError,
			Message:      "delete token from redis error",
			ErrorMessage: err,
		}
	}

	return nil
}

// ProcessThawUser is a function to thaw user
func (s *UserServiceImpl) ProcessThawUser(ctx context.Context, userID string) *model.ApiError {
	err := s.ThawUser(ctx, userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.ThawUserError,
			Message:      "thaw user to MySQL error",
			ErrorMessage: err,
		}
	}
	return nil
}
