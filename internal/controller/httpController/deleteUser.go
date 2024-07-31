package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func (ctrl *UserControllerImpl) DeleteUserHandle(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ResponseErrorWithMessage(c, code.RequestUnauthorized, "User ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		apiError := ctrl.ProcessDeleteUser(ctx, userID)
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "result", true)
		cancel()
		return
	}()

	select {
	case <-ctx.Done():
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
			ResponseErrorWithCode(c, code.ServerBusy)
			return
		}
		if ctx.Err().Error() == context.Canceled.Error() {
			if ctx.Value("error") != nil {
				ResponseErrorWithApiError(c, ctx.Value("error").(*model.ApiError))
				return
			}
			if ctx.Value("result") != nil {
				ResponseSuccess(c, nil)
				return
			}
			ResponseErrorWithCode(c, code.ServerBusy)
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	}
}
