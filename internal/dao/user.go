package dao

import (
	"JWT_authorization/model"
	"JWT_authorization/util/MySQL"
	"context"
	"errors"
	"time"
)

// GetUserInformationByID gets the user information from MySQL
func GetUserInformationByID(userID int) (*model.User, error) {
	db := MySQL.GetMySQL()
	var user *model.User
	result := db.Where("id = ?", userID).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func GetUserInformationByUsername(username string) (*model.User, error) {
	db := MySQL.GetMySQL()
	var user *model.User
	result := db.Where("username  = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func CheckUserExists(username string, telephone string) (usernameExists bool, telephoneExists bool, err error) {
	db := MySQL.GetMySQL()

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
func CreateUser(username string, password string, telephone string) (*model.User, error) {
	db := MySQL.GetMySQL()
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

func CheckUserFrozen(userID string) (bool, error) {
	db := MySQL.GetMySQL()
	var user *model.User
	result := db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return false, result.Error
	}
	return user.IsFrozen, nil
}

// FreezeUser freezes a user in MySQL
func FreezeUser(userID string) error {
	db := MySQL.GetMySQL()
	result := db.Model(&model.User{}).Where("id = ?", userID).Update("is_frozen", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ThawUser thaws a user in MySQL
func ThawUser(userID string) error {
	db := MySQL.GetMySQL()
	result := db.Model(&model.User{}).Where("id = ?", userID).Update("is_frozen", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteUser(userID string) error {
	db := MySQL.GetMySQL()
	result := db.Where("id = ?", userID).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserPermissions(userID string) (int, error) {
	db := MySQL.GetMySQL()
	var user *model.User
	result := db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.Permission, nil
}

func ChangeUserPermissions(userID string, newPermissions int) error {
	db := MySQL.GetMySQL()
	result := db.Model(&model.User{}).Where("id = ?", userID).Update("permission", newPermissions)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
