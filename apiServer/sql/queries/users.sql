-- name: CreateUser :one
insert into users
    (id, username, email, password, created_at, updated_at) 
        values ($1, $2, $3, $4, $5, $6) returning *;

-- name: FindUsernameOrEmail :one
select username, email from users
    where username=$1 or email=$2 limit 1;
