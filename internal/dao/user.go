package dao

import (
	"JWT_authorization/model"
	"context"
	"errors"
	"time"
)

// GetUserInformationByID gets the user information from MySQL
func (dao *UserDAOImpl) GetUserInformationByID(ctx context.Context, userID int) (*model.User, error) {
	db := dao.DB
	var user *model.User
	result := db.Where("id = ?", userID).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (dao *UserDAOImpl) GetUserInformationByUsername(ctx context.Context, username string) (*model.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	db := dao.DB
	var user *model.User
	result := db.Where("username  = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (dao *UserDAOImpl) CheckUserExists(ctx context.Context, username string, telephone string) (usernameExists bool, telephoneExists bool, err error) {
	db := dao.DB

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		usernameExists  bool
		telephoneExists bool
		err             error
	}

	resultChan := make(chan result, 2)

	go func() {
		var usernameCount int64
		res := db.Model(&model.User{}).Where("username = ?", username).Count(&usernameCount)
		if res.Error != nil {
			resultChan <- result{false, false, res.Error}
			cancel()
			return
		}
		if usernameCount > 0 {
			resultChan <- result{true, false, nil}
			cancel()
			return
		}
		resultChan <- result{false, false, nil}
	}()

	go func() {
		var telephoneCount int64
		res := db.Model(&model.User{}).Where("telephone = ?", telephone).Count(&telephoneCount)
		if res.Error != nil {
			resultChan <- result{false, false, res.Error}
			cancel()
			return
		}
		if telephoneCount > 0 {
			resultChan <- result{false, true, nil}
			cancel()
			return
		}
		resultChan <- result{false, false, nil}
	}()

	var usernameResult, telephoneResult result
	for i := 0; i < 2; i++ {
		select {
		case res := <-resultChan:
			if res.err != nil {
				return false, false, res.err
			}
			if res.usernameExists {
				usernameResult.usernameExists = true
			}
			if res.telephoneExists {
				telephoneResult.telephoneExists = true
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return false, false, ctx.Err()
			}
		}
	}

	return usernameResult.usernameExists, telephoneResult.telephoneExists, nil
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
	result := db.Model(&model.User{}).Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

func (dao *UserDAOImpl) CheckUserFrozen(ctx context.Context, userID string) (bool, error) {
	db := dao.DB
	var user *model.User
	result := db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return false, result.Error
	}
	return user.IsFrozen, nil
}

// FreezeUser freezes a user in MySQL
func (dao *UserDAOImpl) FreezeUser(ctx context.Context, userID string) error {
	db := dao.DB
	result := db.Model(&model.User{}).Where("id = ?", userID).Update("is_frozen", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ThawUser thaws a user in MySQL
func (dao *UserDAOImpl) ThawUser(ctx context.Context, userID string) error {
	db := dao.DB
	result := db.Model(&model.User{}).Where("id = ?", userID).Update("is_frozen", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao *UserDAOImpl) DeleteUser(ctx context.Context, userID string) error {
	db := dao.DB
	result := db.Where("id = ?", userID).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao *UserDAOImpl) GetUserPermissions(ctx context.Context, userID string) (int, error) {
	db := dao.DB
	var user *model.User
	result := db.Where("id = ?", userID).First(&user)
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
	result := db.Model(&model.User{}).Where("id = ?", userID).Update("permission", newPermissions)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
