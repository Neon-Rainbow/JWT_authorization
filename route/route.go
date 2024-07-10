package route

import (
	"JWT_authorization/internal/controller"
	"JWT_authorization/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	rootRouter := router.Group("/api/v1")

	authGroup := rootRouter.Group("/auth")
	{
		authGroup.POST("/login", controller.LoginHandler)
		authGroup.POST("/admin_login", controller.AdminLoginHandle)

		authGroup.POST("/register", controller.RegisterHandle)

		authGroup.GET("/refresh", controller.RefreshTokenHandle)

	}

	userGroup := rootRouter.Group("/user")
	userGroup.Use(middleware.JWTMiddleware())
	{
		userGroup.POST("/logout", middleware.JWTMiddleware(), controller.LogoutHandle)

		userGroup.POST("/frozen_account", controller.FreezeUserHandle)
		userGroup.POST("/delete_account", controller.DeleteUserHandle)

		userGroup.GET("/check_permission", controller.CheckUserPermissionsHandle)
		userGroup.GET("/get_user_permission", controller.GetUserPermissionHandle)
	}

	adminGroup := rootRouter.Group("/admin")
	adminGroup.Use(middleware.JWTMiddleware(), middleware.AdminMiddleware())
	{
		adminGroup.POST("/frozen_account", controller.FreezeUserHandle)
		adminGroup.POST("/thaw_account", controller.ThawUserHandle)
		adminGroup.POST("/delete_account", controller.DeleteUserHandle)

		adminGroup.GET("/check_permission", controller.CheckUserPermissionsHandle)
		adminGroup.GET("/get_user_permission", controller.GetUserPermissionHandle)

		adminGroup.POST("/add_permission", controller.AddUserPermissionHandle)
		adminGroup.POST("/delete_permission", controller.DeleteUserPermissionHandle)
	}

	return router

}
