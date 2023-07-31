package userapi

import (
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/minio/minio-go/v7"
	userdb "github.com/openfms/user-api/db/postgres"
	"github.com/openfms/user-api/envconfig"
	"go.uber.org/zap"
	"net/http"
)

type UserHTTPServer struct {
	db          userdb.UserDBPgConn
	minioClient *minio.Client
	log         *zap.Logger
	envConfig   *envconfig.UserEnvConfig
	mux         *runtime.ServeMux
}

func NewUserHTTPServer(logger *zap.Logger, dbConn userdb.UserDBPgConn, env *envconfig.UserEnvConfig,
	minioCli *minio.Client, mux *runtime.ServeMux) (*UserHTTPServer, error) {
	server := &UserHTTPServer{
		log:         logger,
		db:          dbConn,
		envConfig:   env,
		minioClient: minioCli,
		mux:         mux,
	}
	err := server.InitializeRoutes()
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (uhs *UserHTTPServer) InitializeRoutes() error {
	if e := uhs.mux.HandlePath("POST", "/api/v1/user/avatar/upload", uhs.UploadAvatarHandler()); e != nil {
		return e
	}
	if e := uhs.mux.HandlePath("GET", "/api/v1/user/avatar/download", uhs.DownloadAvatarHandler()); e != nil {
		return e
	}
	return nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
