package userapi

import (
	userpb "github.com/openfms/protos/gen/user/v1"
	"github.com/openfms/user-api/db"
	"go.uber.org/zap"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	logger *zap.Logger
	userDB db.UserDBConn
}

func NewUserService(logger *zap.Logger, dbConn db.UserDBConn) *UserService {
	return &UserService{
		logger: logger,
		userDB: dbConn,
	}
}
