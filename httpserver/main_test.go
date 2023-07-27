package httpserver

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("USERDB_POSTGRES", "postgres://${PGDB_USER:-admin}:${PGDB_PASSWORD:-12345678}@postgresdb:5432/${PGDB_NAME:-pgdb}?sslmode=disable")
	os.Setenv("MINIO_ENDPOINT", "localhost:9000")
	os.Setenv("MINIO_ACCESS_KEY", "StTsqPyvT8n34qD1m45U")
	os.Setenv("MINIO_SECRET_KEY", "PuMDCJFYbEKNVd7mUYgzAWfPxMW9t8LmsmqZ3Ev3")
	os.Setenv("MINIO_AVATARS_BUCKET", "user-avatars")
	os.Setenv("JWT_SECRET", "5E5A188D066EC3EFC640")
	os.Setenv("DOMAIN", "openfms.com")
	exitCode := m.Run()
	os.Exit(exitCode)
}
