package dao

import (
	"JWT_authorization/model"
	"context"
	"errors"
	"sync"
)

// GetUserInformationByID gets the user information from MySQL
func (dao *UserDAOImpl) GetUserInformationByID(ctx context.Context, userID int) (*model.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		db := dao.DB
		var user *model.User
		result := db.Where("id = ?", userID).First(user)
		if result.Error != nil {
			return nil, result.Error
		}
		return user, nil
	}
}

func (dao *UserDAOImpl) GetUserInformationByUsername(ctx context.Context, username string) (*model.User, error) {
	db := dao.DB
	var user *model.User
	result := db.WithContext(ctx).Where("username  = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (dao *UserDAOImpl) CheckUserExists(ctx context.Context, username string, telephone string) (usernameExists bool, telephoneExists bool, err error) {

	db := dao.DB

	type result struct {
		usernameExists  bool
		telephoneExists bool
		err             error
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		var usernameCount int64
		res := db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&usernameCount)
		if res.Error != nil {
			ctx = context.WithValue(ctx, "error", res.Error)
			return
		}
		if usernameCount > 0 {
			ctx = context.WithValue(ctx, "usernameExists", true)
			return
		}
		ctx = context.WithValue(ctx, "usernameExists", false)
		return
	}()

	go func() {
		defer wg.Done()
		var telephoneCount int64
		res := db.WithContext(ctx).Model(&model.User{}).Where("telephone = ?", telephone).Count(&telephoneCount)
		if res.Error != nil {
			ctx = context.WithValue(ctx, "error", res.Error)
			return
		}
		if telephoneCount > 0 {
			ctx = context.WithValue(ctx, "telephoneExists", true)
			return
		}
		ctx = context.WithValue(ctx, "telephoneExists", false)
		return
	}()

	go func() {
		wg.Wait()
		usernameExists = ctx.Value("usernameExists").(bool)
		telephoneExists = ctx.Value("telephoneExists").(bool)
		err = ctx.Value("error").(error)
		return
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return false, false, ctx.Err()
		}
		if ctx.Value("error") != nil {
			return false, false, ctx.Value("error").(error)
		}
		return usernameExists, telephoneExists, nil
	}
}

// CreateUser creates a new user in MySQL
func (dao *UserDAOImpl) CreateUser(ctx context.Context, username string, password string, telephone string) (*model.User, error) {
	db := dao.DB
	user := &model.User{
		Username:  username,
		Password:  password,
		Telephone: telephone,
		IsFrozen:  false,
		IsAdmin:   false,
	}
	result := db.WithContext(ctx).Model(&model.User{}).Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

func (dao *UserDAOImpl) CheckUserFrozen(ctx context.Context, userID string) (bool, error) {
	db := dao.DB
	var user *model.User
	result := db.WithContext(ctx).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return false, result.Error
	}
	return user.IsFrozen, nil
}

// FreezeUser freezes a user in MySQL
func (dao *UserDAOImpl) FreezeUser(ctx context.Context, userID string) error {
	db := dao.DB
	result := db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("is_frozen", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ThawUser thaws a user in MySQL
func (dao *UserDAOImpl) ThawUser(ctx context.Context, userID string) error {
	db := dao.DB
	result := db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("is_frozen", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao *UserDAOImpl) DeleteUser(ctx context.Context, userID string) error {
	db := dao.DB
	result := db.WithContext(ctx).Where("id = ?", userID).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao *UserDAOImpl) GetUserPermissions(ctx context.Context, userID string) (int, error) {
	db := dao.DB
	var user *model.User
	result := db.WithContext(ctx).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.Permission, nil
}

// ChangeUserPermissions changes the permissions of a user in MySQL
// @param userID: the ID of the user
// @param newPermissions: the new permissions of the user, represented as an integer, with each bit representing a permission
func (dao *UserDAOImpl) ChangeUserPermissions(ctx context.Context, userID string, newPermissions int) error {
	db := dao.DB
	result := db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("permission", newPermissions)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
