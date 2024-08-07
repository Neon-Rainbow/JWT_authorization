package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

// FreezeUserHandle is a function to frozen user
func (ctrl *UserControllerImpl) FreezeUserHandle(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ResponseErrorWithMessage(c, code.RequestUnauthorized, "User ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		apiError := ctrl.ProcessFreezeUser(ctx, userID)
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		return
	}()

	go func() {
		defer wg.Done()
		apiError := ctrl.ChangeUserPermissions(ctx, userID, 0) // 0 means no permission
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		return
	}()

	go func() {
		wg.Wait()
		ctx = context.WithValue(ctx, "result", true)
		cancel()
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

func (ctrl *UserControllerImpl) ThawUserHandle(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ResponseErrorWithMessage(c, code.RequestUnauthorized, "User ID is required")
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		apiError := ctrl.ProcessThawUser(ctx, userID)
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
