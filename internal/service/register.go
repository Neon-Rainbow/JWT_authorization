package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
)

func ProcessRegisterRequest(req *model.UserRegisterRequest) (*model.UserRegisterResponse, *model.ApiError) {
	if req.Username == "" || req.Password == "" {
		return nil, &model.ApiError{
			Code:         code.RegisterParamsError,
			Message:      code.RegisterParamsError.Message(),
			ErrorMessage: nil,
		}
	}

	usernameExists, telephoneExists, err := dao.CheckUserExists(req.Username, req.Telephone)
	if err != nil {
		return nil, &model.ApiError{
			Code:         code.RegisterCheckUserExistsError,
			Message:      code.RegisterCheckUserExistsError.Message(),
			ErrorMessage: err,
		}
	}
	if usernameExists {
		return nil, &model.ApiError{
			Code:         code.RegisterUsernameExists,
			Message:      code.RegisterUsernameExists.Message(),
			ErrorMessage: nil,
		}
	}

	if telephoneExists {
		return nil, &model.ApiError{
			Code:         code.RegisterTelephoneExists,
			Message:      code.RegisterTelephoneExists.Message(),
			ErrorMessage: nil,
		}
	}

	user, err := dao.CreateUser(req.Username, EncryptPassword(req.Password), req.Telephone)
	if err != nil {
		return nil, &model.ApiError{
			Code:         code.RegisterCreateUserError,
			Message:      code.RegisterCreateUserError.Message(),
			ErrorMessage: err,
		}
	}

	return &model.UserRegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}
