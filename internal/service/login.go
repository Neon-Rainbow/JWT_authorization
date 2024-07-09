package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
	"JWT_authorization/util/jwt"
	"strconv"
)

func ProcessLoginRequest(req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError) {
	dbUser, err := dao.GetUserInformationByUsername(req.Username)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &model.ApiError{
				Code:    code.LoginUserNotFound,
				Message: code.LoginUserNotFound.Message(),
				Error:   err,
			}
		}
		return nil, &model.ApiError{
			Code:    code.LoginGetUserInformationError,
			Message: code.LoginGetUserInformationError.Message(),
			Error:   err,
		}
	}

	if dbUser.IsFrozen {
		return nil, &model.ApiError{
			Code:    code.LoginUserIsFrozen,
			Message: code.LoginUserIsFrozen.Message(),
			Error:   nil,
		}
	}

	if dbUser.Password != EncryptPassword(req.Password) {
		return nil, &model.ApiError{
			Code:    code.LoginPasswordError,
			Message: code.LoginPasswordError.Message(),
			Error:   nil,
		}
	}

	accessToken, refreshToken, err := jwt.GenerateToken(dbUser.Username, dbUser.ID, false)
	if err != nil {
		return nil, &model.ApiError{
			Code:    code.LoginGenerateTokenError,
			Message: code.LoginGenerateTokenError.Message(),
			Error:   err,
		}
	}

	loginResponse := &model.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       dbUser.ID,
		Username:     dbUser.Username,
		IsAdmin:      dbUser.IsAdmin,
		IsFrozen:     dbUser.IsFrozen,
	}

	err = dao.SetTokenToRedis(strconv.Itoa(int(dbUser.ID)), refreshToken)
	if err != nil {
		return nil, &model.ApiError{
			Code:    code.LoginGenerateTokenError,
			Message: code.LoginGenerateTokenError.Message(),
			Error:   err,
		}
	}

	return loginResponse, nil
}

func ProcessAdminLoginRequest(req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError) {
	dbUser, err := dao.GetUserInformationByUsername(req.Username)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &model.ApiError{
				Code:    code.LoginUserNotFound,
				Message: code.LoginUserNotFound.Message(),
				Error:   err,
			}
		}
		return nil, &model.ApiError{
			Code:    code.LoginGetUserInformationError,
			Message: code.LoginGetUserInformationError.Message(),
			Error:   err,
		}
	}

	if dbUser.IsFrozen {
		return nil, &model.ApiError{
			Code:    code.LoginUserIsFrozen,
			Message: code.LoginUserIsFrozen.Message(),
			Error:   nil,
		}
	}

	if dbUser.Password != EncryptPassword(req.Password) {
		return nil, &model.ApiError{
			Code:    code.LoginPasswordError,
			Message: code.LoginPasswordError.Message(),
			Error:   nil,
		}
	}

	if !dbUser.IsAdmin {
		return nil, &model.ApiError{
			Code:    code.LoginPasswordError,
			Message: "user is not admin",
			Error:   nil,
		}
	}

	accessToken, refreshToken, err := jwt.GenerateToken(dbUser.Username, dbUser.ID, dbUser.IsAdmin)
	if err != nil {
		return nil, &model.ApiError{
			Code:    code.LoginGenerateTokenError,
			Message: code.LoginGenerateTokenError.Message(),
			Error:   err,
		}
	}

	loginResponse := &model.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       dbUser.ID,
		Username:     dbUser.Username,
		IsAdmin:      dbUser.IsAdmin,
		IsFrozen:     dbUser.IsFrozen,
	}

	err = dao.SetTokenToRedis(strconv.Itoa(int(dbUser.ID)), refreshToken)
	if err != nil {
		return nil, &model.ApiError{
			Code:    code.LoginGenerateTokenError,
			Message: code.LoginGenerateTokenError.Message(),
			Error:   err,
		}
	}

	return loginResponse, nil
}
