package controller

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/service"
	"JWT_authorization/model"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// LogoutHandle handles logout requests
func LogoutHandle(c *gin.Context) {
	userID := GetUserID(c)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError)
	resultChannel := make(chan bool)

	go func() {
		apiError := service.ProcessLogoutRequest(userID)
		if apiError != nil {
			errorChannel <- apiError
			return
		}
		resultChannel <- true
		return
	}()

	select {
	case <-ctx.Done():
		ResponseErrorWithCode(c, code.RequestTimeout)
		return
	case err := <-errorChannel:
		ResponseErrorWithApiError(c, err)
		return
	case <-resultChannel:
		ResponseSuccess(c, nil)
		return
	}
}
