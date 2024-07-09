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

	authRouter := rootRouter.Group("/auth")
	{
		authRouter.POST("/login", controller.LoginHandler)
		authRouter.POST("/register", controller.RegisterHandle)
		authRouter.POST("/refresh", controller.RefreshTokenHandle)
	}

	userRouter := rootRouter.Group("/user")
	{
		userRouter.GET("/info", middleware.JWTMiddleware(), controller.GetUserInfo)
	}

	adminGroup := rootRouter.Group("/admin")
	{
		adminGroup.GET("/info", middleware.JWTMiddleware(), middleware.AdminMiddleware(), controller.GetAdminInfo)
	}

	return router

}
