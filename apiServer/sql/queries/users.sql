-- name: CreateUser :one
insert into users 
    (id, username, email, password, role, created_at, updated_at) 
        values ($1, $2, $3, $4, $5, $6, $7) returning *;

