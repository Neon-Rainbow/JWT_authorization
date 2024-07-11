package route

import (
	"JWT_authorization/config"
	"JWT_authorization/internal/controller/gRPCController"
	"JWT_authorization/internal/controller/httpController"
	"JWT_authorization/internal/dao"
	"JWT_authorization/internal/service"
	"JWT_authorization/middleware"
	"JWT_authorization/proto"
	"JWT_authorization/util/MySQL"
	"JWT_authorization/util/Redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	userDAO := dao.NewUserDAOImpl(MySQL.GetMySQL(), Redis.GetRedis())
	userService := service.NewUserService(*userDAO)
	ctrl := httpController.NewUserController(*userService)

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	rootRouter := router.Group("/api/v1")

	authGroup := rootRouter.Group("/auth")
	{
		authGroup.POST("/login", ctrl.LoginHandler)
		authGroup.POST("/admin_login", ctrl.AdminLoginHandle)

		authGroup.POST("/register", ctrl.RegisterHandle)

		authGroup.GET("/refresh", ctrl.RefreshTokenHandle)

	}

	userGroup := rootRouter.Group("/user")
	userGroup.Use(middleware.JWTMiddleware())
	{
		userGroup.POST("/logout", ctrl.LogoutHandle)
		userGroup.POST("/frozen_account", ctrl.FreezeUserHandle)
		userGroup.POST("/delete_account", ctrl.DeleteUserHandle)

		userGroup.GET("/check_permission", ctrl.CheckUserPermissionsHandle)
		userGroup.GET("/get_user_permission", ctrl.GetUserPermissionHandle)
	}

	adminGroup := rootRouter.Group("/admin")
	adminGroup.Use(middleware.JWTMiddleware(), middleware.AdminMiddleware())
	{
		adminGroup.POST("/frozen_account", ctrl.FreezeUserHandle)
		adminGroup.POST("/thaw_account", ctrl.ThawUserHandle)
		adminGroup.POST("/delete_account", ctrl.DeleteUserHandle)

		adminGroup.GET("/check_permission", ctrl.CheckUserPermissionsHandle)
		adminGroup.GET("/get_user_permission", ctrl.GetUserPermissionHandle)

		adminGroup.POST("/add_permission", ctrl.AddUserPermissionHandle)
		adminGroup.POST("/delete_permission", ctrl.DeleteUserPermissionHandle)
	}

	return router

}

func StartGRPCServer() {
	address := fmt.Sprintf("%v:%v", config.GetConfig().GRPC.Address, config.GetConfig().GRPC.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.InterceptorSelector()),
	)

	userDAO := dao.NewUserDAOImpl(MySQL.GetMySQL(), Redis.GetRedis())
	userService := service.NewUserService(*userDAO)
	proto.RegisterJwtAuthorizationServiceServer(grpcServer, gRPCController.NewJwtAuthorizationServiceServer(*userService))
	reflection.Register(grpcServer)

	log.Println(fmt.Sprintf("gRPC server is running on %v", address))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
