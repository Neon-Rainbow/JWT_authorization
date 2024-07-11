package middleware

import (
	"JWT_authorization/util/jwt"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// AuthInterceptor 解析Authorization头的拦截器
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		var token string
		if auths, ok := md["authorization"]; ok && len(auths) > 0 {
			token = strings.TrimPrefix(auths[0], "Bearer ")
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %v", err)
		}
		if claims.TokenType != "access_token" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token type, must be access_token")
		}

		// 将token添加到上下文中
		newCtx := context.WithValue(ctx, "AuthorizationToken", token)
		newCtx = context.WithValue(newCtx, "UserID", claims.UserID)
		newCtx = context.WithValue(newCtx, "Username", claims.Username)
		newCtx = context.WithValue(newCtx, "IsAdmin", claims.IsAdmin)

		return handler(newCtx, req)
	}
}

// AdminInterceptor 检查用户是否是管理员
func AdminInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		isAdmin, ok := ctx.Value("IsAdmin").(bool)
		if !ok || !isAdmin {
			return nil, status.Errorf(codes.PermissionDenied, "user is not an admin")
		}

		return handler(ctx, req)
	}
}

func RequiresAuthAndAdminInterceptor(method string) bool {
	switch method {
	case "/proto.JwtAuthorizationService/AdminLogin",
		"/proto.JwtAuthorizationService/AdminFrozenAccount",
		"/proto.JwtAuthorizationService/AdminThawAccount",
		"/proto.JwtAuthorizationService/AdminDeleteAccount",
		"/proto.JwtAuthorizationService/AdminCheckPermission",
		"/proto.JwtAuthorizationService/AdminGetUserPermission",
		"/proto.JwtAuthorizationService/AdminAddUserPermission",
		"/proto.JwtAuthorizationService/AdminDeleteUserPermission":
		return true
	default:
		return false
	}
}

func requiresAuthInterceptor(method string) bool {
	switch method {
	case "/proto.JwtAuthorizationService/UserLogout",
		"/proto.JwtAuthorizationService/UserFrozen",
		"/proto.JwtAuthorizationService/CheckUserPermission",
		"/proto.JwtAuthorizationService/GetUserPermission":
		return true
	default:
		return false
	}
}
