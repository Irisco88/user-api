package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	userpb "github.com/irisco88/protos/gen/user/v1"
)

const updateUserQuery = `
	UPDATE users SET 
		first_name = $1,
		last_name=$2, 
		user_name=$3,
		email=$4,
		password=crypt($5, gen_salt('bf')),
		role=$6,
		avatar=$7
	WHERE 
	    id=$8
	RETURNING id;
`

// UpdateUser updates a new user
func (udb *UserDB) UpdateUser(ctx context.Context, user *userpb.User) error {
	err := udb.GetPgConn().QueryRow(ctx, updateUserQuery,
		user.GetFirstName(),
		user.GetLastName(),
		user.GetUserName(),
		user.GetEmail(),
		user.GetPassword(),
		user.GetRole(),
		user.GetAvatar(),
		user.GetId(),
	).Scan(&user.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "idx_users_unique_email" ||
				pgErr.ConstraintName == "idx_users_unique_user_name" {
				return ErrUserNameEmailExists
			}
		}
		return err
	}
	return nil
}
