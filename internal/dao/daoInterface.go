package dao

import (
	"JWT_authorization/model"
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserDAO interface {
	GetUserInformationByID(ctx context.Context, userID int) (*model.User, error)
	GetUserInformationByUsername(ctx context.Context, username string) (*model.User, error)
	CheckUserExists(ctx context.Context, username string, telephone string) (usernameExists bool, telephoneExists bool, err error)
	CreateUser(ctx context.Context, username string, password string, telephone string) (*model.User, error)
	CheckUserFrozen(ctx context.Context, userID string) (bool, error)
	FreezeUser(ctx context.Context, userID string) error
	ThawUser(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) error
	GetUserPermissions(ctx context.Context, userID string) (int, error)
	ChangeUserPermissions(ctx context.Context, userID string, newPermissions int) error
	SetTokenToRedis(ctx context.Context, userID string, refreshToken string) error
	GetTokenFromRedis(ctx context.Context, userID string) (refreshToken string, err error)
	DeleteTokenFromRedis(ctx context.Context, userID string) error
}

type UserDAOImpl struct {
	*gorm.DB
	*redis.Client
	context.Context
}

func NewUserDAOImpl(db *gorm.DB, rdb *redis.Client) *UserDAOImpl {
	return &UserDAOImpl{db, rdb, context.Background()}
}
