package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type loginResult struct {
	Response *model.UserLoginResponse
	ApiError *model.ApiError
}

// LoginHandler handles login requests
func (ctrl *UserControllerImpl) LoginHandler(c *gin.Context) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive the login result
	resultChannel := make(chan loginResult)

	go func() {
		var loginRequest model.UserLoginRequest
		// Bind the request data to the loginRequest struct
		if err := c.ShouldBind(&loginRequest); err != nil {
			// Send the error to the result channel
			resultChannel <- loginResult{
				ApiError: &model.ApiError{
					Code:    code.LoginParamsError,
					Message: code.LoginParamsError.Message(),
				},
			}
			return
		}

		// Handle the login logic using the service layer
		loginResponse, apiError := ctrl.ProcessLoginRequest(loginRequest)
		// Send the response or error to the result channel
		resultChannel <- loginResult{
			Response: loginResponse,
			ApiError: apiError,
		}
	}()

	// Use select to handle context timeout and completion
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			// Send a timeout error response if the context times out
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
	case result := <-resultChannel:
		if result.ApiError != nil {
			ResponseErrorWithApiError(c, result.ApiError)
			return
		}
		if result.Response != nil {
			ResponseSuccess(c, *result.Response)
			return
		}
	}

	// Default error response if no other conditions are met
	ResponseErrorWithCode(c, code.ServerBusy)
}

// AdminLoginHandle handles admin login requests
func (ctrl *UserControllerImpl) AdminLoginHandle(c *gin.Context) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive the login result
	resultChannel := make(chan loginResult)

	go func() {
		var loginRequest model.UserLoginRequest
		// Bind the request data to the loginRequest struct
		if err := c.ShouldBind(&loginRequest); err != nil {
			// Send the error to the result channel
			resultChannel <- loginResult{
				ApiError: &model.ApiError{
					Code:    code.LoginParamsError,
					Message: code.LoginParamsError.Message(),
				},
			}
			return
		}

		// Handle the login logic using the service layer
		loginResponse, apiError := ctrl.ProcessAdminLoginRequest(loginRequest)
		// Send the response or error to the result channel
		resultChannel <- loginResult{
			Response: loginResponse,
			ApiError: apiError,
		}
	}()

	// Use select to handle context timeout and completion
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			// Send a timeout error response if the context times out
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
	case result := <-resultChannel:
		if result.ApiError != nil {
			ResponseErrorWithApiError(c, result.ApiError)
			return
		}
		if result.Response != nil {
			ResponseSuccess(c, *result.Response)
			return
		}
	}

	// Default error response if no other conditions are met
	ResponseErrorWithCode(c, code.ServerBusy)
}
