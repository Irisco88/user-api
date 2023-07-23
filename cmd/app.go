package main

import (
	"fmt"
	"github.com/openfms/authutil"
	userpb "github.com/openfms/protos/gen/user/v1"
	"github.com/openfms/user-api/db"
	"github.com/openfms/user-api/userapi"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	HostAddress    string
	PortNumber     uint
	DebugMode      bool
	LogRequests    bool
	UserDBPostgres string
	SecretKey      string
	TokenValidTime time.Duration
	Domain         string
)

func main() {
	randSecret, err := authutil.GenerateRandomSecretKey(10)
	if err != nil {
		log.Fatal(err)
	}
	app := &cli.App{
		Name:  "userapi",
		Usage: "user service",
		Commands: []*cli.Command{
			{
				Name:  "user",
				Usage: "starts user api",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "host",
						Usage:       "host address",
						Value:       "0.0.0.0",
						DefaultText: "0.0.0.0",
						Destination: &HostAddress,
						EnvVars:     []string{"HOST"},
					},
					&cli.UintFlag{
						Name:        "port",
						Usage:       "server port number",
						Value:       5000,
						DefaultText: "5000",
						Aliases:     []string{"p"},
						Destination: &PortNumber,
						EnvVars:     []string{"PORT"},
					},
					&cli.BoolFlag{
						Name:        "debug",
						Usage:       "enable debug mode",
						Value:       false,
						DefaultText: "false",
						Destination: &DebugMode,
						EnvVars:     []string{"DEBUG_MODE"},
						Required:    false,
					},
					&cli.BoolFlag{
						Name:        "logreqs",
						Usage:       "enable logging requests",
						Value:       false,
						DefaultText: "false",
						Destination: &LogRequests,
						EnvVars:     []string{"LOG_REQUESTS"},
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "fmsdb",
						Usage:       "fmsdb postgres url",
						Value:       "postgres://admin:password@127.0.0.1:9423/default",
						DefaultText: "postgres://admin:password@127.0.0.1:9423/default",
						Destination: &UserDBPostgres,
						EnvVars:     []string{"USERDB_POSTGRES"},
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "secret",
						Usage:       "jwt secret",
						Value:       randSecret,
						DefaultText: randSecret,
						EnvVars:     []string{"JWT_SECRET"},
						Destination: &SecretKey,
					},
					&cli.StringFlag{
						Name:        "domain",
						Usage:       "server domain name",
						Required:    true,
						EnvVars:     []string{"DOMAIN"},
						Destination: &Domain,
					},
					&cli.DurationFlag{
						Name:        "valid-time",
						Usage:       "jwt toke valid time duration",
						Value:       time.Hour * 48,
						DefaultText: "48 hour",
						EnvVars:     []string{"JWT_VALID_TIME"},
						Destination: &TokenValidTime,
					},
				},
				Action: func(ctx *cli.Context) error {
					loggerConfig := zap.NewProductionConfig()
					if DebugMode {
						loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
					}
					logger, err := loggerConfig.Build()
					if err != nil {
						return err
					}
					server := grpcServer(logger, LogRequests)
					reflection.Register(server)
					addr := net.JoinHostPort(HostAddress, fmt.Sprintf("%d", PortNumber))
					lis, err := net.Listen("tcp", addr)
					if err != nil {
						return fmt.Errorf("faild to make listen address:%v", err)
					}
					userDB, err := db.NewUserDB(UserDBPostgres)
					if err != nil {
						return err
					}
					authManager := authutil.NewAuthManager(SecretKey, Domain, TokenValidTime)
					userSrv := userapi.NewUserService(logger, userDB, authManager)
					userpb.RegisterUserServiceServer(server, userSrv)
					go func() {
						logger.Info("Server running ",
							zap.String("host", HostAddress),
							zap.Uint("port", PortNumber),
						)
						if err := server.Serve(lis); err != nil {
							logger.Fatal("Failed to serve",
								zap.Error(err))
							return
						}
					}()
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					server.Stop()
					return nil
				},
			},
		},
	}

	if e := app.Run(os.Args); e != nil {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("create new logger failed:%v\n", err)
		}
		logger.Error("failed to run app", zap.Error(e))
	}
}