package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
)

func CheckPermission(userID string, permissionNumber int) (isAllowed bool, apiError *model.ApiError) {
	userPermissions, err := dao.GetUserPermissions(userID)
	if err != nil {
		return false, &model.ApiError{
			Code:    code.PermissionGetError,
			Message: "get user permissions error",
			Error:   err,
		}
	}
	return (1<<(permissionNumber-1))&userPermissions != 0, nil
}

func AddPermission(userID string, permissionNumber int) *model.ApiError {
	userPermissions, err := dao.GetUserPermissions(userID)
	if err != nil {
		return &model.ApiError{
			Code:    code.PermissionGetError,
			Message: "get user permissions error",
			Error:   err,
		}
	}

	newPermissions := userPermissions | (1 << (permissionNumber - 1))

	err = dao.ChangeUserPermissions(userID, newPermissions)
	if err != nil {
		return &model.ApiError{
			Code:    code.PermissionChangeError,
			Message: "add user permissions error",
			Error:   err,
		}
	}
	return nil
}

func DeletePermission(userID string, permissionNumber int) *model.ApiError {
	userPermissions, err := dao.GetUserPermissions(userID)
	if err != nil {
		return &model.ApiError{
			Code:    code.PermissionGetError,
			Message: "get user permissions error",
			Error:   err,
		}
	}

	newPermissions := userPermissions &^ (1 << (permissionNumber - 1))

	err = dao.ChangeUserPermissions(userID, newPermissions)
	if err != nil {
		return &model.ApiError{
			Code:    code.PermissionChangeError,
			Message: "add user permissions error",
			Error:   err,
		}
	}
	return nil
}

func GetUserPermissions(userID string) (int, *model.ApiError) {
	userPermissions, err := dao.GetUserPermissions(userID)
	if err != nil {
		return 0, &model.ApiError{
			Code:    code.PermissionGetError,
			Message: "get user permissions error",
			Error:   err,
		}
	}
	return userPermissions, nil
}
