package service

import (
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
	"context"
)

type UserService interface {
	ProcessLoginRequest(ctx context.Context, req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError)
	ProcessAdminLoginRequest(ctx context.Context, req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError)
	ProcessLogoutRequest(ctx context.Context, userID string) *model.ApiError
	CheckPermission(ctx context.Context, userID string, permissionNumber int) (isAllowed bool, apiError *model.ApiError)
	AddPermission(ctx context.Context, userID string, permissionNumber int) *model.ApiError
	DeletePermission(ctx context.Context, userID string, permissionNumber int) *model.ApiError
	GetUserPermissions(ctx context.Context, userID string) (int, *model.ApiError)
	ChangeUserPermissions(ctx context.Context, userID string, permission int) *model.ApiError
	ProcessRefreshToken(ctx context.Context, refreshToken string) (newAccessToken string, apiError *model.ApiError)
	ProcessRegisterRequest(ctx context.Context, req *model.UserRegisterRequest) (*model.UserRegisterResponse, *model.ApiError)
	ProcessDeleteUser(ctx context.Context, userID string) *model.ApiError
	ProcessFreezeUser(ctx context.Context, userID string) (apiError *model.ApiError)
	ProcessThawUser(ctx context.Context, userID string) *model.ApiError
}

type UserServiceImpl struct {
	dao.UserDAOImpl
}

func NewUserService(userDAO dao.UserDAOImpl) *UserServiceImpl {
	return &UserServiceImpl{userDAO}
}
