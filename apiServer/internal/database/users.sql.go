// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
insert into users
    (id, username, email, password, created_at, updated_at) 
        values ($1, $2, $3, $4, $5, $6) returning id, username, password, email, refresh_token, role, created_at, updated_at
`

type CreateUserParams struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.RefreshToken,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUserViaId = `-- name: DeleteUserViaId :one
delete from users where id=$1 returning id, username, password, email, refresh_token, role, created_at, updated_at
`

func (q *Queries) DeleteUserViaId(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, deleteUserViaId, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.RefreshToken,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUsernameOrEmail = `-- name: FindUsernameOrEmail :one
select username, email, id from users
    where username=$1 or email=$2 limit 1
`

type FindUsernameOrEmailParams struct {
	Username string
	Email    string
}

type FindUsernameOrEmailRow struct {
	Username string
	Email    string
	ID       uuid.UUID
}

func (q *Queries) FindUsernameOrEmail(ctx context.Context, arg FindUsernameOrEmailParams) (FindUsernameOrEmailRow, error) {
	row := q.db.QueryRowContext(ctx, findUsernameOrEmail, arg.Username, arg.Email)
	var i FindUsernameOrEmailRow
	err := row.Scan(&i.Username, &i.Email, &i.ID)
	return i, err
}

const findUsernameOrEmailForLogin = `-- name: FindUsernameOrEmailForLogin :one
select id, username, password, email, refresh_token, role, created_at, updated_at from users
    where username=$1 or email=$1
`

func (q *Queries) FindUsernameOrEmailForLogin(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUsernameOrEmailForLogin, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.RefreshToken,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUseFromId = `-- name: GetUseFromId :one
select id, username, password, email, refresh_token, role, created_at, updated_at from users where id=$1 limit 1
`

func (q *Queries) GetUseFromId(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUseFromId, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.RefreshToken,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
select id, username, password, email, refresh_token, role, created_at, updated_at from users where username=$1
`

func (q *Queries) GetUserByName(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.RefreshToken,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateRefreshToken = `-- name: UpdateRefreshToken :one
update users set refresh_token=$1 where username=$2 returning id, username, password, email, refresh_token, role, created_at, updated_at
`

type UpdateRefreshTokenParams struct {
	RefreshToken sql.NullString
	Username     string
}

func (q *Queries) UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateRefreshToken, arg.RefreshToken, arg.Username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.RefreshToken,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
