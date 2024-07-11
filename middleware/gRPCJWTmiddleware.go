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
	case "/proto.JwtAuthorizationService/UserLogin",
		"/proto.JwtAuthorizationService/AdminLogin",
		"/proto.JwtAuthorizationService/UserRegister",
		"/proto.JwtAuthorizationService/RefreshToken":
		return false
	default:
		return true
	}
}

// InterceptorSelector 根据方法名称选择拦截器
func InterceptorSelector() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var interceptors []grpc.UnaryServerInterceptor

		method := info.FullMethod
		if requiresAuthInterceptor(method) {
			interceptors = append(interceptors, AuthInterceptor())
		}

		// 如果没有选择拦截器，则直接调用处理器
		if len(interceptors) == 0 {
			return handler(ctx, req)
		}

		// 链接拦截器
		chainedInterceptor := chainInterceptors(interceptors...)
		return chainedInterceptor(ctx, req, info, handler)
	}
}

func chainInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var lastHandler grpc.UnaryHandler
		lastHandler = handler

		for i := len(interceptors) - 1; i >= 0; i-- {
			next := lastHandler
			currentInterceptor := interceptors[i]
			lastHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return currentInterceptor(currentCtx, currentReq, info, next)
			}
		}

		return lastHandler(ctx, req)
	}
}
