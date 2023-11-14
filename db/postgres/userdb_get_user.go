package postgres

import (
	"context"
	"errors"
	userpb "github.com/irisco88/protos/gen/user/v1"
	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotFound = errors.New("user not found by username or email")
)

const getUserByEmailUserNameQuery = `
	SELECT 
	    id,
	    first_name,
	    COALESCE(last_name,''),
	    user_name,
	    COALESCE(email,''),
	    password,
	    role,
	    COALESCE(avatar,'')
	FROM 
	    users 
	WHERE 
	    email=$1 OR user_name=$1;
`

// GetUserByEmailUserName returns a user by email or username
func (udb *UserDB) GetUserByEmailUserName(ctx context.Context, userNameEmail string) (*userpb.User, error) {
	user := &userpb.User{}
	err := udb.GetPgConn().QueryRow(ctx, getUserByEmailUserNameQuery, userNameEmail).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Avatar,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

const getUserQuery = `
	SELECT 
	    id,
	    first_name,
	    COALESCE(last_name,''),
	    user_name,
	    COALESCE(email,''),
	    password,
	    role,
	    COALESCE(avatar,'')
	FROM 
	    users 
	WHERE 
	    id=$1;
`

// GetUser returns a user by its id
func (udb *UserDB) GetUser(ctx context.Context, userID uint32) (*userpb.User, error) {
	user := &userpb.User{}
	err := udb.GetPgConn().QueryRow(ctx, getUserQuery, userID).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Avatar,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
