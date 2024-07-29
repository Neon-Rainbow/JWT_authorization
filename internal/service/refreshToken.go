package service

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"JWT_authorization/util/jwt"
	"context"
	"strconv"
)

// ProcessRefreshToken is a function to process refresh token
func (s *UserServiceImpl) ProcessRefreshToken(ctx context.Context, refreshToken string) (newAccessToken string, apiError *model.ApiError) {
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return "", &model.ApiError{
			Code:         code.RefreshTokenError,
			Message:      "refresh token invalid",
			ErrorMessage: err,
		}
	}

	redisRefreshToken, err := s.GetTokenFromRedis(ctx, strconv.Itoa(int(claims.UserID)))
	if err != nil {
		return "", &model.ApiError{
			Code:         code.RefreshTokenError,
			Message:      code.RefreshTokenError.Message(),
			ErrorMessage: err,
		}
	}
	if redisRefreshToken != refreshToken {
		return "", &model.ApiError{
			Code:         code.RefreshTokenError,
			Message:      "refresh token not match, please login again",
			ErrorMessage: nil,
		}
	}

	isFrozen, err := s.CheckUserFrozen(ctx, strconv.Itoa(int(claims.UserID)))
	if err != nil {
		return "", &model.ApiError{
			Code:         code.RefreshTokenError,
			Message:      "check user frozen error",
			ErrorMessage: err,
		}
	}

	if isFrozen {
		return "", &model.ApiError{
			Code:         code.RefreshTokenError,
			Message:      "user is frozen",
			ErrorMessage: nil,
		}
	}

	newAccessToken, _, err = jwt.GenerateToken(claims.Username, claims.UserID, claims.IsAdmin)
	if err != nil {
		return "", &model.ApiError{
			Code:         code.RefreshTokenError,
			Message:      "generate token error",
			ErrorMessage: err,
		}
	}
	return newAccessToken, nil
}
