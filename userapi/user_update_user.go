package userapi

import (
	"context"
	"errors"
	"github.com/openfms/authutil"
	userv1pb "github.com/openfms/protos/gen/user/v1"
	userdb "github.com/openfms/user-api/db/postgres"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *UserService) UpdateUser(ctx context.Context, req *userv1pb.UpdateUserRequest) (*userv1pb.UpdateUserResponse, error) {
	claims, ok := authutil.TokenClaimsFromCtx(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid context")
	}
	if err := us.ValidateUpdateUser(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := us.userDB.UpdateUser(ctx, claims.Role, claims.UserID, req.User); err != nil {
		if errors.Is(err, userdb.ErrUserNameEmailExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		us.logger.Error("failed to update user",
			zap.Error(err),
			zap.Uint32("userID", claims.UserID),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userv1pb.UpdateUserResponse{}, nil
}
