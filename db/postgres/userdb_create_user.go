package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	userpb "github.com/openfms/protos/gen/user/v1"
)

var (
	ErrUserNameEmailExists = errors.New("email or username already exists")
)

const createUserQuery = `
	INSERT INTO users (first_name, last_name, user_name, email, password, role,avatar,owner_id)
		VALUES ($1,$2,$3,$4,crypt($5, gen_salt('bf')),$6,$7,$8) RETURNING id;
`

// CreateUser creates a new user
func (udb *UserDB) CreateUser(ctx context.Context, ownerID uint32, user *userpb.User) error {
	err := udb.GetPgConn().QueryRow(ctx, createUserQuery,
		user.GetFirstName(),
		user.GetLastName(),
		user.GetUserName(),
		user.GetEmail(),
		user.GetPassword(),
		user.GetRole(),
		user.GetAvatar(),
		ownerID,
	).Scan(&user.Id)
	pgErr, ok := err.(*pgconn.PgError)
	if ok && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "idx_users_unique_email" ||
			pgErr.ConstraintName == "idx_users_unique_user_name" {
			return ErrUserNameEmailExists
		}
	}
	return nil
}
