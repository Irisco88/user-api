package envconfig

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type UserEnvConfig struct {
	UserDatabase       string        `env:"USERDB_POSTGRES,notEmpty"`
	DebugMode          bool          `env:"DEBUG_MODE" envDefault:"false"`
	Host               string        `env:"HOST" envDefault:"0.0.0.0"`
	Port               string        `env:"PORT" envDefault:"3000"`
	LogRequests        bool          `env:"LOG_REQUESTS" envDefault:"false"`
	JWTSecret          string        `env:"JWT_SECRET"`
	Domain             string        `env:"DOMAIN"`
	JwtValidTime       time.Duration `env:"JWT_VALID_TIME"`
	MinioEndpoint      string        `env:"MINIO_ENDPOINT,notEmpty"`
	MinioAccessKey     string        `env:"MINIO_ACCESS_KEY,notEmpty"`
	MinioSecretKey     string        `env:"MINIO_SECRET_KEY,notEmpty"`
	MinioAvatarsBucket string        `env:"MINIO_AVATARS_BUCKET" envDefault:"user-avatars"`
	UserAvatarMaxSize  int           `env:"USER_AVATAR_MAX_SIZE" envDefault:"5242880"` //bytes
}

func ReadUserEnvironment() (*UserEnvConfig, error) {
	cfg := &UserEnvConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
