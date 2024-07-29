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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		inputPermission := c.Query("permission_number")
		userID := GetUserID(c)

		permissionNumber, err := strconv.Atoi(inputPermission)
		if err != nil {
			ctx = context.WithValue(ctx, "error", &model.ApiError{
				Code:         code.PermissionParamsError,
				Message:      code.PermissionParamsError.Message(),
				ErrorMessage: err,
			})
			cancel()
			return
		}

		isAllowed, apiError := ctrl.CheckPermission(ctx, userID, permissionNumber)
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "result", isAllowed)
		cancel()
		return
	}()

	select {
	case <-ctx.Done():
		if ctx.Err().Error() == context.DeadlineExceeded.Error() {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		if ctx.Value("error") != nil {
			ResponseErrorWithApiError(c, ctx.Value("error").(*model.ApiError))
			return
		}
		if ctx.Value("result") != nil {
			ResponseSuccess(c, ctx.Value("result"))
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	}
}

func (ctrl *UserControllerImpl) GetUserPermissionHandle(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		userID := GetUserID(c)

		userPermissions, apiError := ctrl.GetUserPermissions(ctx, userID)
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "result", userPermissions)
		cancel()
		return
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		if ctx.Value("error") != nil {
			ResponseErrorWithApiError(c, ctx.Value("error").(*model.ApiError))
			return
		}
		if ctx.Value("result") != nil {
			permissions := gin.H{}
			userPermissions := ctx.Value("result").(int)
			for i := 1; i <= 7; i++ {
				// 相当于添加 "Permission1": userPermissions&(1<<0) != 0 这个字段
				permissions[fmt.Sprintf("Permission%d", i)] = userPermissions&(1<<(i-1)) != 0
			}
			ResponseSuccess(c, permissions)
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	}
}

func (ctrl *UserControllerImpl) AddUserPermissionHandle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	go func() {
		inputPermission := c.Query("permission_number")
		userID := GetUserID(c)

		permissionNumber, err := strconv.Atoi(inputPermission)
		if err != nil {
			ctx = context.WithValue(ctx, "error", &model.ApiError{
				Code:         code.PermissionParamsError,
				Message:      code.PermissionParamsError.Message(),
				ErrorMessage: err,
			})
			cancel()
			return
		}

		apiError := ctrl.AddPermission(ctx, userID, permissionNumber)
		if apiError != nil {
			ctx = context.WithValue(ctx, "error", apiError)
			cancel()
			return
		}
		ctx = context.WithValue(ctx, "result", true)
		cancel()
		return
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			ResponseErrorWithCode(c, code.RequestTimeout)
			return
		}
		if ctx.Value("error") != nil {
			ResponseErrorWithApiError(c, ctx.Value("error").(*model.ApiError))
			return
		}
		if ctx.Value("result").(bool) {
			ResponseSuccess(c, nil)
			return
		}
		ResponseErrorWithCode(c, code.ServerBusy)
		return
	}
}

func (ctrl *UserControllerImpl) DeleteUserPermissionHandle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
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

		apiError := ctrl.DeletePermission(ctx, userID, permissionNumber)
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
