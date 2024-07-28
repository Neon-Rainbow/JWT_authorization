package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// LogoutHandle handles logout requests
func (ctrl *UserControllerImpl) LogoutHandle(c *gin.Context) {
	userID := GetUserID(c)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		apiError := ctrl.ProcessLogoutRequest(ctx, userID)
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "result", "Logout successful")
		cancel()
		return
	}()

	select {
	case <-ctx.Done():
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
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
}
