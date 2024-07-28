package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
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
			ctx = context.WithValue(ctx, "error", &model.ApiError{
				Code:         code.LoginParamsError,
				Message:      code.LoginParamsError.Message(),
				ErrorMessage: err,
			})
			cancel()
			return
		}

		// Handle the login logic using the service layer
		loginResponse, apiError := ctrl.ProcessLoginRequest(ctx, loginRequest)
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
		ResponseSuccess(c, ctx.Value("value").(*model.UserLoginResponse))
		return
	}
}

// AdminLoginHandle handles admin login requests
func (ctrl *UserControllerImpl) AdminLoginHandle(c *gin.Context) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Create a channel to receive the login result
	resultChannel := make(chan *model.UserLoginResponse)

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
		resultChannel <- loginResponse
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
		ResponseErrorWithCode(c, code.ServerBusy)
	case result := <-resultChannel:
		ResponseSuccess(c, result)
		return
	}
}
