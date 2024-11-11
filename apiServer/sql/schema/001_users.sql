-- +goose Up
create type user_role AS enum ('admin', 'user');

create table users (
    id uuid primary key,
    username text not null unique,
    password text not null,
    email text not null unique,
    refresh_token text,
    role user_role,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
drop table users;
drop type user_role;
