package postgres

import (
	"context"
	"errors"
	commonpb "github.com/irisco88/protos/gen/common/v1"
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

const ListUsersQuery = `
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
	    users;
`

func (udb *UserDB) ListUsers(ctx context.Context) ([]*userpb.User, error) {
	rows, err := udb.GetPgConn().Query(ctx, ListUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	usersDatas := make([]*userpb.User, 0)
	for rows.Next() {
		var (
			//	usersData  = &userpb.User{}
			id         uint32
			first_name string
			last_name  string
			user_name  string
			email      string
			password   string
			role       int32
			avatar     string
		)

		err := rows.Scan(
			&id,
			&first_name,
			&last_name,
			&user_name,
			&email,
			&password,
			&role,
			&avatar,
		)

		if err != nil {
			return nil, err
		}
		//usersData.Role = intToUserRole(role)
		usersData := &userpb.User{
			Id:        id,
			FirstName: first_name,
			LastName:  last_name,
			UserName:  user_name,
			Email:     email,
			Password:  password,
			Role:      intToUserRole(role),
			Avatar:    avatar,
		}
		usersDatas = append(usersDatas, usersData)
	}

	return usersDatas, nil
}
func intToUserRole(role int32) commonpb.UserRole {
	switch role {
	case 0:
		return commonpb.UserRole_USER_ROLE_NORMAL
	case 1:
		return commonpb.UserRole_USER_ROLE_READER
	case 2:
		return commonpb.UserRole_USER_ROLE_ADMIN
	default:
		return commonpb.UserRole_USER_ROLE_NORMAL
	}
}
