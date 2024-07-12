package httpController

import (
	"JWT_authorization/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	LoginHandler(c *gin.Context)
	AdminLoginHandle(c *gin.Context)
	LogoutHandle(c *gin.Context)
	CheckUserPermissionsHandle(c *gin.Context)
	GetUserPermissionHandle(c *gin.Context)
	AddUserPermissionHandle(c *gin.Context)
	DeleteUserPermissionHandle(c *gin.Context)
	RefreshTokenHandle(c *gin.Context)
	RegisterHandle(c *gin.Context)
	DeleteUserHandle(c *gin.Context)
	FreezeUserHandle(c *gin.Context)
	ThawUserHandle(c *gin.Context)
}

type UserControllerImpl struct {
	service.UserServiceImpl
}

func NewUserController(userService service.UserServiceImpl) *UserControllerImpl {
	return &UserControllerImpl{userService}
}
