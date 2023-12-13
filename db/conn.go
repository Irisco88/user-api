package db

import (
	userdb "github.com/irisco88/user-api/db/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=$GOFILE -destination=mock_db/conn.go -package=$GOPACKAG
type UserDBConn interface {
	GetPgConn() *pgxpool.Pool
	userdb.UserDBPgConn
}

var _ UserDBConn = &UserDataBase{}

type UserDataBase struct {
	pgConn *pgxpool.Pool
	*userdb.UserDB
}

func (tdb *UserDataBase) GetPgConn() *pgxpool.Pool {
	return tdb.pgConn
}

func NewUserDB(pgURL string) (*UserDataBase, error) {
	fmsConn, err := userdb.ConnectToUserDB(pgURL)
	if err != nil {
		return nil, err
	}
	return &UserDataBase{
		UserDB: userdb.NewUserDB(fmsConn),
	}, nil
}
