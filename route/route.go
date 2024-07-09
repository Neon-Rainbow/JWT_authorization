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

	rootRouter := router.Group("/api")

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
		userGroup.GET("/info", controller.GetUserInfo)
		userGroup.POST("/frozen", controller.FreezeUserHandle)
		userGroup.POST("/delete_account", controller.DeleteUserHandle)
	}

	adminGroup := rootRouter.Group("/admin")
	adminGroup.Use(middleware.JWTMiddleware(), middleware.AdminMiddleware())
	{
		adminGroup.GET("/info", controller.GetAdminInfo)
		adminGroup.POST("/frozen", controller.FreezeUserHandle)
		adminGroup.POST("/thaw", controller.ThawUserHandle)
		adminGroup.POST("/delete_account", controller.DeleteUserHandle)
	}

	return router

}
