package service

import (
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
)

type UserService interface {
	ProcessLoginRequest(req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError)
	ProcessAdminLoginRequest(req model.UserLoginRequest) (*model.UserLoginResponse, *model.ApiError)
	ProcessLogoutRequest(userID string) *model.ApiError
	CheckPermission(userID string, permissionNumber int) (isAllowed bool, apiError *model.ApiError)
	AddPermission(userID string, permissionNumber int) *model.ApiError
	DeletePermission(userID string, permissionNumber int) *model.ApiError
	GetUserPermissions(userID string) (int, *model.ApiError)
	ChangeUserPermissions(userID string, permission int) *model.ApiError
	ProcessRefreshToken(refreshToken string) (newAccessToken string, apiError *model.ApiError)
	ProcessRegisterRequest(req *model.UserRegisterRequest) (*model.UserRegisterResponse, *model.ApiError)
	ProcessDeleteUser(userID string) *model.ApiError
	ProcessFreezeUser(userID string) (apiError *model.ApiError)
	ProcessThawUser(userID string) *model.ApiError
}

type UserServiceImpl struct {
	dao.UserDAOImpl
}

func NewUserService(userDAO dao.UserDAOImpl) *UserServiceImpl {
	return &UserServiceImpl{userDAO}
}
