package service

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
)

func (s *UserServiceImpl) ProcessFreezeUser(userID string) (apiError *model.ApiError) {
	err := s.FreezeUser(userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.FrozenUserError,
			Message:      "frozen user to MySQL error",
			ErrorMessage: err,
		}
	}
	err = s.DeleteTokenFromRedis(userID)
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
func (s *UserServiceImpl) ProcessThawUser(userID string) *model.ApiError {
	err := s.ThawUser(userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.ThawUserError,
			Message:      "thaw user to MySQL error",
			ErrorMessage: err,
		}
	}
	return nil
}
