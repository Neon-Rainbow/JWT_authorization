package controller

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/service"
	"JWT_authorization/model"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

type refreshTokenResult struct {
	Response *model.RefreshTokenResponse
	ApiError *model.ApiError
}

func RefreshTokenHandle(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		ResponseErrorWithMessage(c, code.RefreshTokenError, "refresh_token is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChannel := make(chan refreshTokenResult)

	go func() {
		accessToken, err := service.ProcessRefreshToken(refreshToken)
		if err != nil {
			resultChannel <- refreshTokenResult{
				ApiError: err,
				Response: nil,
			}
			return
		}
		resultChannel <- refreshTokenResult{
			ApiError: nil,
			Response: &model.RefreshTokenResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}
	}()

	select {
	case result := <-resultChannel:
		if result.ApiError != nil {
			ResponseErrorWithApiError(c, result.ApiError)
			return
		}
		if result.Response != nil {
			ResponseSuccess(c, result.Response)
			return
		}
	case <-ctx.Done():
		ResponseErrorWithCode(c, code.RequestTimeout)
		return
	default:

	}
	ResponseErrorWithCode(c, code.ServerBusy)
	return

}
