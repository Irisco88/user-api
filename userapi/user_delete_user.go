package userapi

import (
	"context"
	"errors"
	"github.com/openfms/authutil"
	commonpb "github.com/openfms/protos/gen/common/v1"
	userv1pb "github.com/openfms/protos/gen/user/v1"
	userdb "github.com/openfms/user-api/db/postgres"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUser removes an existed user
func (us *UserService) DeleteUser(ctx context.Context, req *userv1pb.DeleteUserRequest) (*userv1pb.DeleteUserResponse, error) {
	claims, ok := authutil.TokenClaimsFromCtx(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid context")
	}
	if err := us.ValidateDeleteUser(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !(claims.Role == commonpb.UserRole_USER_ROLE_ADMIN ||
		(claims.Role == commonpb.UserRole_USER_ROLE_NORMAL && claims.UserID == req.UserId)) {
		return nil, status.Error(codes.Unauthenticated, "invalid access")
	}
	if err := us.userDB.DeleteUser(ctx, req.UserId); err != nil {
		if errors.Is(err, userdb.ErrUserNameEmailExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, userdb.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		us.logger.Error("failed to update user",
			zap.Error(err),
			zap.Uint32("userID", claims.UserID),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userv1pb.DeleteUserResponse{}, nil
}
