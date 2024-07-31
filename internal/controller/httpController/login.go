package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// LoginHandler handles login requests
func (ctrl *UserControllerImpl) LoginHandler(c *gin.Context) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		var loginRequest model.UserLoginRequest
		// Bind the request data to the loginRequest struct
		if err := c.ShouldBind(&loginRequest); err != nil {
			// Send the error to the result channel
			apiError := &model.ApiError{
				Code:         code.LoginParamsError,
				Message:      code.LoginParamsError.Message(),
				ErrorMessage: err,
			}
			ctx = context.WithValue(ctx, "error", apiError)
			zap.L().Error("LoginHandler:请求参数错误", zap.Error(apiError))

			cancel()
			return
		}

		// Handle the login logic using the service layer
		loginResponse, apiError := ctrl.ProcessLoginRequest(ctx, loginRequest)
		// Send the response or error to the result channel
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			zap.L().Error("LoginHandler:登录失败", zap.Error(apiError))
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "result", loginResponse)
		cancel()
		return
	}()

	select {
	case <-ctx.Done():
		fmt.Println(ctx)
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			// Send a timeout error response if the context times out
			ResponseErrorWithCode(c, code.RequestTimeout)
			zap.L().Error("LoginHandler:请求超时", zap.Error(ctx.Err()))
			return
		}
		if ctx.Value("error") != nil {
			ResponseErrorWithApiError(c, ctx.Value("error").(*model.ApiError))
			return
		}
		if ctx.Value("result") != nil {
			ResponseSuccess(c, ctx.Value("result").(*model.UserLoginResponse))
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	}
}

// AdminLoginHandle handles admin login requests
func (ctrl *UserControllerImpl) AdminLoginHandle(c *gin.Context) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		var loginRequest model.UserLoginRequest
		// Bind the request data to the loginRequest struct
		if err := c.ShouldBind(&loginRequest); err != nil {
			// Send the error to the result channel
			ctx = context.WithValue(ctx, "error", &model.ApiError{
				Code:         code.LoginParamsError,
				Message:      code.LoginParamsError.Message(),
				ErrorMessage: err,
			})
			cancel()
			return
		}

		// Handle the login logic using the service layer
		loginResponse, apiError := ctrl.ProcessAdminLoginRequest(ctx, loginRequest)
		// Send the response or error to the result channel
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "result", loginResponse)
		cancel()
		return
	}()

	// Use select to handle context timeout and completion
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			// Send a timeout error response if the context times out
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		if ctx.Value("error") != nil {
			ResponseErrorWithApiError(c, ctx.Value("error").(*model.ApiError))
			return
		}
		if ctx.Value("result") != nil {
			ResponseSuccess(c, ctx.Value("result").(*model.UserLoginResponse))
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
	}
}
