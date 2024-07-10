package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
	"context"
	"time"
)

func ProcessDeleteUser(userID string) *model.ApiError {

	errorChannel := make(chan *model.ApiError, 1)
	resultChannel := make(chan bool, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		err := dao.DeleteUser(userID)
		if err != nil {
			errorChannel <- &model.ApiError{
				Code:    code.DeleteUserTokenError,
				Message: "delete user error",
				Error:   err,
			}
			return
		}
		resultChannel <- true
	}()

	go func() {
		err := dao.DeleteTokenFromRedis(userID)
		if err != nil {
			errorChannel <- &model.ApiError{
				Code:    code.DeleteUserTokenError,
				Message: "delete token from redis error",
				Error:   err,
			}
			return
		}
		resultChannel <- true
	}()

	for i := 0; i < 2; i++ {
		select {
		case apiError := <-errorChannel:
			return apiError
		case <-resultChannel:
			continue
		case <-ctx.Done():
			return &model.ApiError{
				Code:    code.RequestTimeout,
				Message: "request timeout",
				Error:   nil,
			}
		}
	}
	return nil
}
