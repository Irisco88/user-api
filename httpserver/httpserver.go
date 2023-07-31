package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/openfms/authutil"
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

func (uhs *UserHTTPServer) Run(host string, port uint) error {
	uhs.log.Info("running http server",
		zap.String("addr", net.JoinHostPort(host, fmt.Sprintf("%d", port))),
	)
	srv := &http.Server{
		Addr:         net.JoinHostPort(host, fmt.Sprintf("%d", port)),
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
	uhs.Router.HandleFunc("/api/v1/user/avatar/download/{code}", uhs.DownloadAvatarHandler).Methods("GET")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]any{"message": message, "code": code})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
