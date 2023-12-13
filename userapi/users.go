package userapi

import (
	"github.com/irisco88/authutil"
	commonpb "github.com/irisco88/protos/gen/common/v1"
	userspb "github.com/irisco88/protos/gen/users/v1"
	"github.com/irisco88/user-api/db"
	"go.uber.org/zap"
)

type UsersService struct {
	userspb.UnimplementedUsersServiceServer
	logger      *zap.Logger
	userDB      db.UserDBConn
	authManager *authutil.AuthManager
}

func NewUsersService(logger *zap.Logger, dbConn db.UserDBConn, auth *authutil.AuthManager) *UsersService {
	return &UsersService{
		logger:      logger,
		userDB:      dbConn,
		authManager: auth,
	}
}

func (us *UsersService) GetAuthManager() *authutil.AuthManager {
	return us.authManager
}

func (us *UsersService) GetRoleAccess(fullMethod string) []commonpb.UserRole {
	methodPerms, ok := userspb.UsersServicePermission.MethodStreams[fullMethod]
	if ok {
		return methodPerms.Roles
	}
	return nil
}
