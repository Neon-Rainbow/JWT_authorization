package httpController

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func (ctrl *UserControllerImpl) CheckUserPermissionsHandle(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError)
	resultChannel := make(chan bool)

	go func() {
		inputPermission := c.Query("permission_number")
		userID := GetUserID(c)

		permissionNumber, err := strconv.Atoi(inputPermission)
		if err != nil {
			errorChannel <- &model.ApiError{
				Code:         code.PermissionParamsError,
				Message:      code.PermissionParamsError.Message(),
				ErrorMessage: err,
			}
			return
		}

		isAllowed, apiError := ctrl.CheckPermission(userID, permissionNumber)
		if apiError != nil {
			errorChannel <- apiError
			return
		}

		resultChannel <- isAllowed
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	case apiError := <-errorChannel:
		ResponseErrorWithApiError(c, apiError)
		return
	case isAllowed := <-resultChannel:
		ResponseSuccess(c, isAllowed)
		return
	}
}

func (ctrl *UserControllerImpl) GetUserPermissionHandle(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError)
	resultChannel := make(chan int)

	go func() {
		userID := GetUserID(c)

		userPermissions, apiError := ctrl.GetUserPermissions(userID)
		if apiError != nil {
			errorChannel <- apiError
			return
		}

		resultChannel <- userPermissions
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	case apiError := <-errorChannel:
		ResponseErrorWithApiError(c, apiError)
		return
	case userPermissions := <-resultChannel:

		permissions := gin.H{}
		for i := 1; i <= 7; i++ {
			// 相当于添加 "Permission1": userPermissions&(1<<0) != 0 这个字段
			permissions[fmt.Sprintf("Permission%d", i)] = userPermissions&(1<<(i-1)) != 0
		}

		ResponseSuccess(c, permissions)
		return
	}
}

func (ctrl *UserControllerImpl) AddUserPermissionHandle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError)
	resultChannel := make(chan bool)

	go func() {
		inputPermission := c.Query("permission_number")
		userID := GetUserID(c)

		permissionNumber, err := strconv.Atoi(inputPermission)
		if err != nil {
			errorChannel <- &model.ApiError{
				Code:         code.PermissionParamsError,
				Message:      code.PermissionParamsError.Message(),
				ErrorMessage: err,
			}
			return
		}

		apiError := ctrl.AddPermission(userID, permissionNumber)
		if apiError != nil {
			errorChannel <- apiError
			return
		}

		resultChannel <- true
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	case apiError := <-errorChannel:
		ResponseErrorWithApiError(c, apiError)
		return
	case <-resultChannel:
		ResponseSuccess(c, nil)
		return
	}
}

func (ctrl *UserControllerImpl) DeleteUserPermissionHandle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError)
	resultChannel := make(chan bool)

	go func() {
		inputPermission := c.Query("permission_number")
		userID := GetUserID(c)

		permissionNumber, err := strconv.Atoi(inputPermission)
		if err != nil {
			errorChannel <- &model.ApiError{
				Code:         code.PermissionParamsError,
				Message:      code.PermissionParamsError.Message(),
				ErrorMessage: err,
			}
			return
		}

		apiError := ctrl.DeletePermission(userID, permissionNumber)
		if apiError != nil {
			errorChannel <- apiError
			return
		}

		resultChannel <- true
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	case apiError := <-errorChannel:
		ResponseErrorWithApiError(c, apiError)
		return
	case <-resultChannel:
		ResponseSuccess(c, nil)
		return
	}
}
