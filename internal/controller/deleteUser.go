package controller

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/service"
	"JWT_authorization/model"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func DeleteUserHandle(c *gin.Context) {
	userID := GetUserID(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError, 1)
	resultChannel := make(chan bool, 1)

	go func() {
		apiError := service.ProcessDeleteUser(userID)
		if apiError != nil {
			errorChannel <- apiError
			return
		}
		resultChannel <- true
	}()

	select {
	case apiError := <-errorChannel:
		ResponseErrorWithApiError(c, apiError)
		return
	case <-resultChannel:
		ResponseSuccess(c, nil)
		return
	case <-ctx.Done():
		ResponseErrorWithCode(c, code.RequestTimeout)
		return
	default:

	}
	ResponseErrorWithCode(c, code.ServerBusy)
	return

}
