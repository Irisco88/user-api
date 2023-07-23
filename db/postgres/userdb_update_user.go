package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	commonpb "github.com/openfms/protos/gen/common/v1"
	userpb "github.com/openfms/protos/gen/user/v1"
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
	    id=$8 AND ($9=2 OR ($9=0 AND id=$10))
	RETURNING id;
`

// UpdateUser updates a new user
func (udb *UserDB) UpdateUser(ctx context.Context, userRole commonpb.UserRole, userID uint32, user *userpb.User) error {
	err := udb.GetPgConn().QueryRow(ctx, updateUserQuery,
		user.GetFirstName(),
		user.GetLastName(),
		user.GetUserName(),
		user.GetEmail(),
		user.GetPassword(),
		user.GetRole(),
		user.GetAvatar(),
		user.GetId(),
		userRole,
		userID,
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
