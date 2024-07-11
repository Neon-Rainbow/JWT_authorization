package gRPCController

import (
	"JWT_authorization/code"
	"JWT_authorization/internal/service"
	"JWT_authorization/model"
	"JWT_authorization/proto"
	"context"
	"fmt"
	"strconv"
)

type JwtAuthorizationServiceServer struct {
	proto.UnimplementedJwtAuthorizationServiceServer
	service.UserServiceImpl
}

func NewJwtAuthorizationServiceServer(userService service.UserServiceImpl) *JwtAuthorizationServiceServer {
	return &JwtAuthorizationServiceServer{UserServiceImpl: userService}
}

func (s *JwtAuthorizationServiceServer) UserLogin(ctx context.Context, req *proto.UserLoginRequest) (*proto.UserLoginResponse, error) {
	resp, apiError := s.UserServiceImpl.ProcessLoginRequest(model.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if apiError != nil {
		return nil, apiError
	}
	return &proto.UserLoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		UserId:       uint32(resp.UserID),
		Username:     resp.Username,
		IsAdmin:      resp.IsAdmin,
		IsFrozen:     resp.IsFrozen,
	}, nil
}

func (s *JwtAuthorizationServiceServer) AdminLogin(ctx context.Context, req *proto.AdminLoginRequest) (*proto.AdminLoginResponse, error) {
	resp, apiError := s.UserServiceImpl.ProcessLoginRequest(model.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if apiError != nil {
		return nil, apiError
	}
	return &proto.AdminLoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		AdminId:      uint32(resp.UserID),
		Username:     resp.Username,
	}, nil
}

func (s *JwtAuthorizationServiceServer) UserRegister(ctx context.Context, req *proto.UserRegisterRequest) (*proto.UserRegisterResponse, error) {
	resp, apiError := s.UserServiceImpl.ProcessRegisterRequest(&model.UserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if apiError != nil {
		return nil, apiError
	}
	return &proto.UserRegisterResponse{
		UserId:   uint32(resp.UserID),
		Username: resp.Username,
	}, nil
}

func (s *JwtAuthorizationServiceServer) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	resp, apiError := s.UserServiceImpl.ProcessRefreshToken(req.RefreshToken)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.RefreshTokenResponse{
		AccessToken: resp,
	}, nil
}

func (s *JwtAuthorizationServiceServer) UserLogout(ctx context.Context, req *proto.UserLogoutRequest) (*proto.UserLogoutResponse, error) {
	jwtToken := ctx.Value("AuthorizationToken").(string)
	apiError := s.UserServiceImpl.ProcessLogoutRequest(jwtToken)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.UserLogoutResponse{}, nil
}

func (s *JwtAuthorizationServiceServer) UserFrozen(ctx context.Context, req *proto.UserFrozenRequest) (*proto.UserFrozenResponse, error) {
	userID := fmt.Sprintf("%v", ctx.Value("UserID"))
	apiError := s.UserServiceImpl.ProcessFreezeUser(userID)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.UserFrozenResponse{}, nil
}

func (s *JwtAuthorizationServiceServer) CheckUserPermission(ctx context.Context, req *proto.UserCheckPermissionRequest) (*proto.UserCheckPermissionResponse, error) {
	userID := fmt.Sprintf("%v", ctx.Value("UserID"))
	permission, _ := strconv.Atoi(req.Permission)
	hasPermission, apiError := s.UserServiceImpl.CheckPermission(userID, permission)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.UserCheckPermissionResponse{
		HasPermission: hasPermission,
	}, nil
}

func (s *JwtAuthorizationServiceServer) GetUserPermissions(ctx context.Context, req *proto.UserGetUserPermissionRequest) (*proto.UserGetUserPermissionResponse, error) {
	userID := fmt.Sprintf("%v", ctx.Value("UserID"))
	permissions, apiError := s.UserServiceImpl.GetUserPermissions(userID)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.UserGetUserPermissionResponse{
		Permissions: strconv.Itoa(permissions),
	}, nil
}

func (s *JwtAuthorizationServiceServer) AdminFrozenAccount(ctx context.Context, req *proto.AdminFrozenAccountRequest) (*proto.AdminFrozenAccountResponse, error) {

	if !ctx.Value("IsAdmin").(bool) {
		return nil, &model.ApiError{
			Code:    code.PermissionDenied,
			Message: "user is not an admin",
		}
	}

	apiError := s.UserServiceImpl.ProcessFreezeUser(req.UserId)

	if apiError != nil {
		return nil, apiError
	}
	return &proto.AdminFrozenAccountResponse{}, nil
}

func (s *JwtAuthorizationServiceServer) AdminThawAccount(ctx context.Context, req *proto.AdminThawAccountRequest) (*proto.AdminThawAccountResponse, error) {
	if !ctx.Value("IsAdmin").(bool) {
		return nil, &model.ApiError{
			Code:    code.PermissionDenied,
			Message: "user is not an admin",
		}
	}
	apiError := s.UserServiceImpl.ProcessThawUser(req.UserId)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.AdminThawAccountResponse{}, nil
}

func (s *JwtAuthorizationServiceServer) AdminDeleteAccount(ctx context.Context, req *proto.AdminDeleteAccountRequest) (*proto.AdminDeleteAccountResponse, error) {
	if !ctx.Value("IsAdmin").(bool) {
		return nil, &model.ApiError{
			Code:    code.PermissionDenied,
			Message: "user is not an admin",
		}
	}
	apiError := s.UserServiceImpl.ProcessDeleteUser(req.UserId)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.AdminDeleteAccountResponse{}, nil
}

func (s *JwtAuthorizationServiceServer) AdminCheckPermission(ctx context.Context, req *proto.AdminCheckPermissionRequest) (*proto.AdminCheckPermissionResponse, error) {
	if !ctx.Value("IsAdmin").(bool) {
		return nil, &model.ApiError{
			Code:    code.PermissionDenied,
			Message: "user is not an admin",
		}
	}
	pms, _ := strconv.Atoi(req.Permission)
	hasPermission, apiError := s.UserServiceImpl.CheckPermission(req.UserId, pms)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.AdminCheckPermissionResponse{
		HasPermission: hasPermission,
	}, nil
}

func (s *JwtAuthorizationServiceServer) AdminAddUserPermission(ctx context.Context, req *proto.AdminAddUserPermissionRequest) (*proto.AdminAddUserPermissionResponse, error) {
	if !ctx.Value("IsAdmin").(bool) {
		return nil, &model.ApiError{
			Code:    code.PermissionDenied,
			Message: "user is not an admin",
		}
	}
	pms, _ := strconv.Atoi(req.Permission)
	apiError := s.UserServiceImpl.AddPermission(req.UserId, pms)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.AdminAddUserPermissionResponse{}, nil
}

func (s *JwtAuthorizationServiceServer) AdminDeleteUserPermission(ctx context.Context, req *proto.AdminDeleteUserPermissionRequest) (*proto.AdminDeleteUserPermissionResponse, error) {
	if !ctx.Value("IsAdmin").(bool) {
		return nil, &model.ApiError{
			Code:    code.PermissionDenied,
			Message: "user is not an admin",
		}
	}
	pms, _ := strconv.Atoi(req.Permission)
	apiError := s.UserServiceImpl.DeletePermission(req.UserId, pms)
	if apiError != nil {
		return nil, apiError
	}
	return &proto.AdminDeleteUserPermissionResponse{}, nil
}
