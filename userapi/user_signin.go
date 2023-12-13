package userapi

import (
	"context"
	"errors"

	userpb "github.com/irisco88/protos/gen/user/v1"
	userdb "github.com/irisco88/user-api/db/postgres"
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
		if errors.Is(err, userdb.ErrUserNotFound) {
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
	us.logger.Error("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&",
		zap.Error(err),
		zap.String("***********************************", user.FirstName+" "+user.LastName+"_"+user.Avatar+"_"+user.Role.String()),
	)
	return &userpb.SignInResponse{
		Token:    token,
		Fullname: user.FirstName + " " + user.LastName,
		Avatar:   user.Avatar,
		Roll:     user.Role.String(),
	}, nil
}
