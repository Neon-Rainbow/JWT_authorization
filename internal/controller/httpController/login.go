package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// handleLoginRequest handles both user and admin login requests
func (ctrl *UserControllerImpl) handleLoginRequest(c *gin.Context, processLogin func(context.Context, model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError)) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resultChan := make(chan interface{})

	go func() {
		defer close(resultChan) // Ensure the channel is closed when goroutine exits
		var loginRequest model.UserLoginRequest
		if err := c.ShouldBind(&loginRequest); err != nil {
			apiError := &model.ApiError{
				Code:         code.LoginParamsError,
				Message:      code.LoginParamsError.Message(),
				ErrorMessage: err,
			}
			zap.L().Error("handleLoginRequest: 请求参数错误", zap.Error(apiError))
			resultChan <- apiError
			return
		}

		loginResponse, apiError := processLogin(ctx, loginRequest)
		if apiError != nil {
			zap.L().Error("handleLoginRequest: 登录失败", zap.Error(apiError))
			resultChan <- apiError
			return
		}

		resultChan <- loginResponse
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ResponseErrorWithCode(c, code.RequestTimeout)
			zap.L().Error("handleLoginRequest: 请求超时", zap.Error(ctx.Err()))
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	case result := <-resultChan:
		switch res := result.(type) {
		case *model.ApiError:
			ResponseErrorWithApiError(c, res)
			return
		case *model.UserLoginResponse:
			ResponseSuccess(c, res)
			return
		default:
			ResponseErrorWithCode(c, code.ServerBusy)
			return
		}
	}
}

// LoginHandler handles login requests
func (ctrl *UserControllerImpl) LoginHandler(c *gin.Context) {
	ctrl.handleLoginRequest(c, ctrl.ProcessLoginRequest)
}

// AdminLoginHandle handles admin login requests
func (ctrl *UserControllerImpl) AdminLoginHandle(c *gin.Context) {
	ctrl.handleLoginRequest(c, ctrl.ProcessAdminLoginRequest)
}
