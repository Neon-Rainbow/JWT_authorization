package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
)

func ProcessFreezeUser(userID string) (apiError *model.ApiError) {
	err := dao.FreezeUser(userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.FrozenUserError,
			Message:      "frozen user to MySQL error",
			ErrorMessage: err,
		}
	}
	err = dao.DeleteTokenFromRedis(userID)
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
func ProcessThawUser(userID string) *model.ApiError {
	err := dao.ThawUser(userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.ThawUserError,
			Message:      "thaw user to MySQL error",
			ErrorMessage: err,
		}
	}
	return nil
}
