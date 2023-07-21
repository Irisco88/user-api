package userapi

import (
	"github.com/openfms/authutil"
	commonpb "github.com/openfms/protos/gen/common/v1"
	userpb "github.com/openfms/protos/gen/user/v1"
	"github.com/openfms/user-api/db"
	"go.uber.org/zap"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	logger      *zap.Logger
	userDB      db.UserDBConn
	authManager *authutil.AuthManager
}

func NewUserService(logger *zap.Logger, dbConn db.UserDBConn, auth *authutil.AuthManager) *UserService {
	return &UserService{
		logger:      logger,
		userDB:      dbConn,
		authManager: auth,
	}
}

func (us *UserService) GetAuthManager() *authutil.AuthManager {
	return us.authManager
}

func (us *UserService) GetRoleAccess(fullMethod string) []commonpb.UserRole {
	methodPerms, ok := userpb.UserServicePermission.MethodStreams[fullMethod]
	if ok {
		return methodPerms.Roles
	}
	return nil
}
