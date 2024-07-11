package service

import (
	"JWT_authorization/code"
	"JWT_authorization/model"
)

// CheckPermission checks if a user has a specific permission
// used to check if the user has the permission of the permissionNumber
func (s *UserServiceImpl) CheckPermission(userID string, permissionNumber int) (isAllowed bool, apiError *model.ApiError) {
	userPermissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return false, &model.ApiError{
			Code:         code.PermissionGetError,
			Message:      "get user permissions error",
			ErrorMessage: err,
		}
	}
	return (1<<(permissionNumber-1))&userPermissions != 0, nil
}

// AddPermission adds a permission to a user
// used to add the permission of the permissionNumber to the user
func (s *UserServiceImpl) AddPermission(userID string, permissionNumber int) *model.ApiError {
	userPermissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.PermissionGetError,
			Message:      "get user permissions error",
			ErrorMessage: err,
		}
	}

	newPermissions := userPermissions | (1 << (permissionNumber - 1))

	err = s.ChangeUserPermissions(userID, newPermissions)
	if err != nil {
		return &model.ApiError{
			Code:         code.PermissionChangeError,
			Message:      "add user permissions error",
			ErrorMessage: err,
		}
	}
	return nil
}

// DeletePermission deletes a permission from a user
// used to delete the permission of the permissionNumber from the user
func (s *UserServiceImpl) DeletePermission(userID string, permissionNumber int) *model.ApiError {
	userPermissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return &model.ApiError{
			Code:         code.PermissionGetError,
			Message:      "get user permissions error",
			ErrorMessage: err,
		}
	}

	newPermissions := userPermissions &^ (1 << (permissionNumber - 1))

	err = s.ChangeUserPermissions(userID, newPermissions)
	if err != nil {
		return &model.ApiError{
			Code:         code.PermissionChangeError,
			Message:      "add user permissions error",
			ErrorMessage: err,
		}
	}
	return nil
}

// GetUserPermissions gets the permissions of a user
// used to get the permissions of the user
func (s *UserServiceImpl) GetUserPermissions(userID string) (int, *model.ApiError) {
	userPermissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return 0, &model.ApiError{
			Code:         code.PermissionGetError,
			Message:      "get user permissions error",
			ErrorMessage: err,
		}
	}
	return userPermissions, nil
}

func (s *UserServiceImpl) ChangeUserPermissions(userID string, permission int) *model.ApiError {
	err := s.ChangeUserPermissions(userID, permission)
	if err != nil {
		return &model.ApiError{
			Code:         code.PermissionChangeError,
			Message:      "change user permissions error",
			ErrorMessage: err,
		}
	}
	return nil
}
