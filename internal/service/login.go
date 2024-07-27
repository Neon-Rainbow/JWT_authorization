package service

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"JWT_authorization/util/jwt"
	"context"
	"strconv"
)

func (s *UserServiceImpl) ProcessLoginRequest(ctx context.Context, req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError) {
	select {
	case <-ctx.Done():
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
			return nil, &model.ApiError{
				Code:         code.RequestTimeout,
				Message:      code.RequestTimeout.Message(),
				ErrorMessage: nil,
			}
		}
		return nil, &model.ApiError{
			Code:         code.RequestCanceled,
			Message:      code.RequestCanceled.Message(),
			ErrorMessage: nil,
		}
	default:
	}

	dbUser, err := s.GetUserInformationByUsername(ctx, req.Username)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &model.ApiError{
				Code:         code.LoginUserNotFound,
				Message:      code.LoginUserNotFound.Message(),
				ErrorMessage: err,
			}
		}
		return nil, &model.ApiError{
			Code:         code.LoginGetUserInformationError,
			Message:      code.LoginGetUserInformationError.Message(),
			ErrorMessage: err,
		}
	}

	if dbUser.IsFrozen {
		return nil, &model.ApiError{
			Code:         code.LoginUserIsFrozen,
			Message:      code.LoginUserIsFrozen.Message(),
			ErrorMessage: nil,
		}
	}

	if dbUser.Password != EncryptPassword(req.Password) {
		return nil, &model.ApiError{
			Code:         code.LoginPasswordError,
			Message:      code.LoginPasswordError.Message(),
			ErrorMessage: nil,
		}
	}

	accessToken, refreshToken, err := jwt.GenerateToken(dbUser.Username, dbUser.ID, false)
	if err != nil {
		return nil, &model.ApiError{
			Code:         code.LoginGenerateTokenError,
			Message:      code.LoginGenerateTokenError.Message(),
			ErrorMessage: err,
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

	err = s.SetTokenToRedis(strconv.Itoa(int(dbUser.ID)), refreshToken)
	if err != nil {
		return nil, &model.ApiError{
			Code:         code.LoginGenerateTokenError,
			Message:      code.LoginGenerateTokenError.Message(),
			ErrorMessage: err,
		}
	}

	return loginResponse, nil
}

func (s *UserServiceImpl) ProcessAdminLoginRequest(ctx context.Context, req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError) {
	select {
	case <-ctx.Done():
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
			return nil, &model.ApiError{
				Code:         code.RequestTimeout,
				Message:      code.RequestTimeout.Message(),
				ErrorMessage: nil,
			}
		}
		return nil, &model.ApiError{
			Code:         code.RequestCanceled,
			Message:      code.RequestCanceled.Message(),
			ErrorMessage: nil,
		}
	default:
	}

	dbUser, err := s.GetUserInformationByUsername(ctx, req.Username)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &model.ApiError{
				Code:         code.LoginUserNotFound,
				Message:      code.LoginUserNotFound.Message(),
				ErrorMessage: err,
			}
		}
		return nil, &model.ApiError{
			Code:         code.LoginGetUserInformationError,
			Message:      code.LoginGetUserInformationError.Message(),
			ErrorMessage: err,
		}
	}

	if dbUser.IsFrozen {
		return nil, &model.ApiError{
			Code:         code.LoginUserIsFrozen,
			Message:      code.LoginUserIsFrozen.Message(),
			ErrorMessage: nil,
		}
	}

	if dbUser.Password != EncryptPassword(req.Password) {
		return nil, &model.ApiError{
			Code:         code.LoginPasswordError,
			Message:      code.LoginPasswordError.Message(),
			ErrorMessage: nil,
		}
	}

	if !dbUser.IsAdmin {
		return nil, &model.ApiError{
			Code:         code.LoginPasswordError,
			Message:      "user is not admin",
			ErrorMessage: nil,
		}
	}

	accessToken, refreshToken, err := jwt.GenerateToken(dbUser.Username, dbUser.ID, dbUser.IsAdmin)
	if err != nil {
		return nil, &model.ApiError{
			Code:         code.LoginGenerateTokenError,
			Message:      code.LoginGenerateTokenError.Message(),
			ErrorMessage: err,
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

	err = s.SetTokenToRedis(strconv.Itoa(int(dbUser.ID)), refreshToken)
	if err != nil {
		return nil, &model.ApiError{
			Code:         code.LoginGenerateTokenError,
			Message:      code.LoginGenerateTokenError.Message(),
			ErrorMessage: err,
		}
	}

	return loginResponse, nil
}
