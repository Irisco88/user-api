package userapi

import (
	"context"
	"errors"
	commonpb "github.com/openfms/protos/gen/common/v1"
	userpb "github.com/openfms/protos/gen/user/v1"
	userdb "github.com/openfms/user-api/db/postgres"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *UserService) SignUp(ctx context.Context, req *userpb.SignUpRequest) (*userpb.SignInResponse, error) {
	if err := us.ValidateSignUpUser(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user := &userpb.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserName:  req.UserName,
		Email:     req.Email,
		Password:  req.Password,
		Avatar:    req.Avatar,
		Role:      commonpb.UserRole_USER_ROLE_NORMAL,
	}
	if err := us.userDB.CreateUser(ctx, 0, user); err != nil {
		if errors.Is(err, userdb.ErrUserNameEmailExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		us.logger.Error("failed to create new user",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("userName", req.UserName),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	token, err := us.authManager.GenerateNewToken(user)
	if err != nil {
		us.logger.Error("create token failed",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("userName", req.UserName),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userpb.SignInResponse{
		Token: token,
	}, nil
}
