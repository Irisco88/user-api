package userapi

import (
	"github.com/irisco88/authutil"
	commonpb "github.com/irisco88/protos/gen/common/v1"
	userpb "github.com/irisco88/protos/gen/user/v1"
	"github.com/irisco88/user-api/db"
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
