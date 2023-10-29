package httpserver

import (
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/irisco88/authutil"
	"github.com/irisco88/user-api/db/mock_db"
	"github.com/irisco88/user-api/envconfig"
	"go.uber.org/zap"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewUserHTTPServer(t *testing.T) {
	logger, _ := zap.NewProduction()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbConn := mock_db.NewMockUserDBConn(ctrl)
	env, err := envconfig.ReadUserEnvironment()
	assert.NilError(t, err)
	client, err := minio.New(env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioAccessKey, env.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	authManager := authutil.NewAuthManager(env.JWTSecret, env.Domain, env.JwtValidTime)
	httpServer := NewUserHTTPServer(logger, dbConn, env, client, authManager)
	assert.NilError(t, httpServer.Run("127.0.0.1", 6060))
}
