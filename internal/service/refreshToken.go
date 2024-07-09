package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
	"JWT_authorization/util/jwt"
	"strconv"
)

// ProcessRefreshToken is a function to process refresh token
func ProcessRefreshToken(refreshToken string) (newAccessToken string, apiError *model.ApiError) {
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return "", &model.ApiError{
			Code:    code.RefreshTokenError,
			Message: "refresh token invalid",
			Error:   err,
		}
	}

	redisRefreshToken, err := dao.GetTokenFromRedis(strconv.Itoa(int(claims.UserID)))
	if err != nil {
		return "", &model.ApiError{
			Code:    code.RefreshTokenError,
			Message: code.RefreshTokenError.Message(),
			Error:   err,
		}
	}
	if redisRefreshToken != refreshToken {
		return "", &model.ApiError{
			Code:    code.RefreshTokenError,
			Message: "refresh token not match, please login again",
			Error:   nil,
		}
	}

	isFrozen, err := dao.CheckUserFrozen(strconv.Itoa(int(claims.UserID)))
	if err != nil {
		return "", &model.ApiError{
			Code:    code.RefreshTokenError,
			Message: "check user frozen error",
			Error:   err,
		}
	}

	if isFrozen {
		return "", &model.ApiError{
			Code:    code.RefreshTokenError,
			Message: "user is frozen",
			Error:   nil,
		}
	}

	newAccessToken, _, err = jwt.GenerateToken(claims.Username, claims.UserID, claims.IsAdmin)
	if err != nil {
		return "", &model.ApiError{
			Code:    code.RefreshTokenError,
			Message: "generate token error",
			Error:   err,
		}
	}
	return newAccessToken, nil
}
