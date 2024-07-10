package controller

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/service"
	"JWT_authorization/model"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func CheckUserPermissionsHandle(c *gin.Context) {

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
				Code:    code.PermissionParamsError,
				Message: code.PermissionParamsError.Message(),
				Error:   err,
			}
			return
		}

		isAllowed, apiError := service.CheckPermission(userID, permissionNumber)
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

func GetUserPermissionHandle(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan *model.ApiError)
	resultChannel := make(chan int)

	go func() {
		userID := GetUserID(c)

		userPermissions, apiError := service.GetUserPermissions(userID)
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

func AddUserPermissionHandle(c *gin.Context) {
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
				Code:    code.PermissionParamsError,
				Message: code.PermissionParamsError.Message(),
				Error:   err,
			}
			return
		}

		apiError := service.AddPermission(userID, permissionNumber)
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

func DeleteUserPermissionHandle(c *gin.Context) {
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
				Code:    code.PermissionParamsError,
				Message: code.PermissionParamsError.Message(),
				Error:   err,
			}
			return
		}

		apiError := service.DeletePermission(userID, permissionNumber)
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
