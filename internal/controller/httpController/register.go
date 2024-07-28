package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type registerResult struct {
	Response *model.UserRegisterResponse
	ApiError *model.ApiError
}

// RegisterHandle handles user registration requests
func (ctrl *UserControllerImpl) RegisterHandle(c *gin.Context) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive the registration result
	resultChannel := make(chan registerResult)

	go func() {
		var registerRequest model.UserRegisterRequest
		err := c.ShouldBind(&registerRequest)
		if err != nil {
			// Send the error to the result channel
			resultChannel <- registerResult{
				ApiError: &model.ApiError{
					Code:         code.RegisterParamsError,
					Message:      code.RegisterParamsError.Message(),
					ErrorMessage: err,
				},
			}
			return
		}

		registerResponse, apiError := ctrl.ProcessRegisterRequest(ctx, &registerRequest)
		// Send the response or error to the result channel
		resultChannel <- registerResult{
			Response: registerResponse,
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
