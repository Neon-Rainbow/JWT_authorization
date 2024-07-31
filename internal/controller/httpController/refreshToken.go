package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func (ctrl *UserControllerImpl) RefreshTokenHandle(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		ResponseErrorWithMessage(c, code.RefreshTokenError, "refresh_token is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		accessToken, err := ctrl.ProcessRefreshToken(ctx, refreshToken)
		if err != nil {
			ctx = context.WithValue(ctx, "error", err)
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "is_success", true)
		ctx = context.WithValue(ctx, "access_token", accessToken)
		ctx = context.WithValue(ctx, "refresh_token", refreshToken)
		cancel()
		return
	}()

	select {
	case <-ctx.Done():
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
			ResponseErrorWithMessage(c, code.RequestTimeout, "Request Timeout")
			return
		}
		if ctx.Value("error") != nil {
			ResponseErrorWithApiError(c, ctx.Value("error").(*model.ApiError))
			return
		}
		if ctx.Value("is_success").(bool) {
			accessToken := ctx.Value("access_token").(string)
			refreshToken := ctx.Value("refresh_token").(string)
			ResponseSuccess(c, gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			})
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	}
}
