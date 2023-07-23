package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

const deleteUserQuery = `
	DELETE FROM users  
		WHERE id=$1
	RETURNING id;
`

// DeleteUser updates a new user
func (udb *UserDB) DeleteUser(ctx context.Context, userID uint32) error {
	err := udb.GetPgConn().QueryRow(ctx, deleteUserQuery,
		userID,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
