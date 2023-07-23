package postgres

import (
	"context"
	sqlmaker "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	commonpb "github.com/openfms/protos/gen/common/v1"
	userpb "github.com/openfms/protos/gen/user/v1"
	"time"
)

type UserDBPgConn interface {
	GetPgConn() *pgxpool.Pool
	GetSQLBuilder() sqlmaker.StatementBuilderType
	GetUserByEmailUserName(ctx context.Context, userNameEmail string) (*userpb.User, error)
	CreateUser(ctx context.Context, ownerID uint32, user *userpb.User) error
	UpdateUser(ctx context.Context, userRole commonpb.UserRole, userID uint32, user *userpb.User) error
}

var _ UserDBPgConn = &UserDB{}

type UserDB struct {
	postgresConn  *pgxpool.Pool
	selectBuilder sqlmaker.StatementBuilderType
}

func (udb *UserDB) GetPgConn() *pgxpool.Pool {
	return udb.postgresConn
}

func (udb *UserDB) GetSQLBuilder() sqlmaker.StatementBuilderType {
	udb.selectBuilder = sqlmaker.StatementBuilder.PlaceholderFormat(sqlmaker.Dollar)
	return udb.selectBuilder
}

func NewUserDB(rawConn *pgxpool.Pool) *UserDB {
	return &UserDB{
		selectBuilder: sqlmaker.StatementBuilder.PlaceholderFormat(sqlmaker.Dollar),
		postgresConn:  rawConn,
	}
}

func ConnectToUserDB(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rawConn, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	return rawConn, nil
}
