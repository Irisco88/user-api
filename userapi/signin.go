package userapi

import (
	"context"
	"errors"
	userpb "github.com/openfms/protos/gen/user/v1"
	"github.com/openfms/user-api/db/postgres"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *UserService) SignIn(ctx context.Context, req *userpb.SignInRequest) (*userpb.SignInResponse, error) {
	if err := us.ValidateSignInUser(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user, err := us.userDB.GetUserByEmailUserName(ctx, req.UserNameEmail)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid username or email or password")
		}
		us.logger.Error("failed to get user by email or username",
			zap.Error(err),
			zap.String("emailUserName", req.GetUserNameEmail()),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); e != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid username or email or password")
	}
	token, err := us.authManager.GenerateNewToken(user)
	if err != nil {
		us.logger.Error("create token failed",
			zap.Error(err),
			zap.String("emailUserName", req.UserNameEmail),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userpb.SignInResponse{
		Token: token,
	}, nil
}
