package httpserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/openfms/authutil"
	commonpb "github.com/openfms/protos/gen/common/v1"
	userdb "github.com/openfms/user-api/db/postgres"
	"github.com/openfms/user-api/envconfig"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type UserHTTPServer struct {
	Router      *mux.Router
	db          userdb.UserDBPgConn
	minioClient *minio.Client
	log         *zap.Logger
	envConfig   *envconfig.UserEnvConfig
	authManager *authutil.AuthManager
}

func NewUserHTTPServer(logger *zap.Logger, dbConn userdb.UserDBPgConn, env *envconfig.UserEnvConfig,
	minioCli *minio.Client, auth *authutil.AuthManager) *UserHTTPServer {
	server := &UserHTTPServer{
		Router:      mux.NewRouter(),
		log:         logger,
		db:          dbConn,
		envConfig:   env,
		minioClient: minioCli,
		authManager: auth,
	}
	server.InitializeRoutes()
	return server
}

func (uhs *UserHTTPServer) Run(host string, port string) error {
	srv := &http.Server{
		Addr:         net.JoinHostPort(host, port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      uhs.Router,
	}
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (uhs *UserHTTPServer) InitializeRoutes() {
	uhs.Router.HandleFunc("/api/v1/user/avatar/upload", uhs.UploadAvatarHandler).Methods("POST")
	uhs.Router.HandleFunc("/api/v1/user/avatar/download", uhs.DownloadAvatarHandler).Methods("GET")
	uhs.Router.Use(authutil.MuxAuthMiddleware(uhs))
}

func (uhs *UserHTTPServer) GetAuthManager() *authutil.AuthManager {
	return uhs.authManager
}

func (uhs *UserHTTPServer) GetRoleAccess(path string) []commonpb.UserRole {
	methodsPerms := map[string][]commonpb.UserRole{
		"/api/v1/user/avatar/upload":   {commonpb.UserRole_USER_ROLE_NORMAL, commonpb.UserRole_USER_ROLE_ADMIN},
		"/api/v1/user/avatar/download": {commonpb.UserRole_USER_ROLE_NORMAL, commonpb.UserRole_USER_ROLE_READER, commonpb.UserRole_USER_ROLE_ADMIN},
	}
	roles, ok := methodsPerms[path]
	if ok {
		return roles
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
