package userapi

import (
	"context"
	userpb "github.com/irisco88/protos/gen/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *UserService) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	users, err := us.userDB.ListUsers(ctx) // Use the embedded interface directly
	if err != nil {
		us.logger.Error("failed to update users", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userpb.ListUsersResponse{
		Users: users, // Assuming GetUsrs() is a method in the GetUsersResponse
	}, nil

	// }

}
