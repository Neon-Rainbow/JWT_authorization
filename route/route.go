package route

import (
	"JWT_authorization/internal/controller"
	"JWT_authorization/internal/service"
	"JWT_authorization/middleware"
	"JWT_authorization/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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
		userGroup.POST("/logout", controller.LogoutHandle)
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

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.AuthInterceptor(),
		),
	)

	proto.RegisterJwtAuthorizationServiceServer(grpcServer, service.NewJwtAuthorizationServiceServer())
	reflection.Register(grpcServer)

	log.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
