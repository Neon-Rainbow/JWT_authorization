package dao

import (
	"JWT_authorization/model"
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserDAO interface {
	GetUserInformationByID(userID int) (*model.User, error)
	GetUserInformationByUsername(username string) (*model.User, error)
	CheckUserExists(username string, telephone string) (usernameExists bool, telephoneExists bool, err error)
	CreateUser(username string, password string, telephone string) (*model.User, error)
	CheckUserFrozen(userID string) (bool, error)
	FreezeUser(userID string) error
	ThawUser(userID string) error
	DeleteUser(userID string) error
	GetUserPermissions(userID string) (int, error)
	ChangeUserPermissions(userID string, newPermissions int) error
	SetTokenToRedis(userID string, refreshToken string) error
	GetTokenFromRedis(userID string) (refreshToken string, err error)
	DeleteTokenFromRedis(userID string) error
}

type UserDAOImpl struct {
	*gorm.DB
	*redis.Client
	context.Context
}

func NewUserDAOImpl(db *gorm.DB, rdb *redis.Client) *UserDAOImpl {
	return &UserDAOImpl{db, rdb, context.Background()}
}
