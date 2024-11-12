-- name: CreateUser :one
insert into users
    (id, username, email, password, created_at, updated_at) 
        values ($1, $2, $3, $4, $5, $6) returning *;

-- name: FindUsernameOrEmail :one
select username, email from users
    where username=$1 or email=$2 limit 1;


-- name: FindUsernameOrEmailForLogin :one
select * from users
    where username=$1 or email=$1;

-- name: UpdateRefreshToken :one
update users set refresh_token=$1 where username=$2 returning *;

-- name: GetUserByName :one
select * from users where username=$1;
