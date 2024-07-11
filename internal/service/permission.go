package service

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/dao"
	"JWT_authorization/model"
)

// CheckPermission checks if a user has a specific permission
// used to check if the user has the permission of the permissionNumber
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

// AddPermission adds a permission to a user
// used to add the permission of the permissionNumber to the user
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

// DeletePermission deletes a permission from a user
// used to delete the permission of the permissionNumber from the user
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

// GetUserPermissions gets the permissions of a user
// used to get the permissions of the user
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

func ChangeUserPermissions(userID string, permission int) *model.ApiError {
	err := dao.ChangeUserPermissions(userID, permission)
	if err != nil {
		return &model.ApiError{
			Code:    code.PermissionChangeError,
			Message: "change user permissions error",
			Error:   err,
		}
	}
	return nil
}
