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
	userID := GetUserID(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError)
	doneChannel := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		apiError := ctrl.ProcessFreezeUser(userID)
		if apiError != nil {
			//errorChannel is a channel without buffer, so it will block until the error is read
			errorChannel <- apiError
			return
		}
		return
	}()

	go func() {
		defer wg.Done()
		apiError := ctrl.ChangeUserPermissions(userID, 0) // 0 means no permission
		if apiError != nil {
			//errorChannel is a channel without buffer, so it will block until the error is read
			errorChannel <- apiError
			return
		}
		return
	}()

	go func() {
		// Wait for the two goroutines to finish
		wg.Wait()
		doneChannel <- true
	}()

	select {
	case apiError := <-errorChannel:
		ResponseErrorWithApiError(c, apiError)
		return
	case <-doneChannel:
		ResponseSuccess(c, nil)
		return
	case <-ctx.Done():
		ResponseErrorWithCode(c, code.RequestTimeout)
		return
	}
}

func (ctrl *UserControllerImpl) ThawUserHandle(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		ResponseErrorWithCode(c, code.ThawUserIDRequired)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError, 1)
	resultChannel := make(chan bool, 1)

	go func() {
		apiError := ctrl.ProcessThawUser(userID)
		if apiError != nil {
			errorChannel <- apiError
			return
		}
		resultChannel <- true
		return
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
	}
}
