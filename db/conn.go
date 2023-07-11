package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openfms/user-api/db/postgres"
)

//go:generate mockgen -source=$GOFILE -destination=mock_db/conn.go -package=$GOPACKAG
type UserDBConn interface {
	GetPgConn() *pgxpool.Pool
	postgres.UserDBPgConn
}

var _ UserDBConn = &UserDataBase{}

type UserDataBase struct {
	pgConn *pgxpool.Pool
	*postgres.UserDB
}

func (tdb *UserDataBase) GetPgConn() *pgxpool.Pool {
	return tdb.pgConn
}

func NewUserDB(pgURL string) (*UserDataBase, error) {
	fmsConn, err := postgres.ConnectToUserDB(pgURL)
	if err != nil {
		return nil, err
	}
	return &UserDataBase{
		UserDB: postgres.NewUserDB(fmsConn),
	}, nil
}
