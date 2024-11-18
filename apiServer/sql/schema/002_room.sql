-- +goose Up
create table rooms (
    id uuid primary key,
    room_name text unique not null,
    admin_id uuid REFERENCES users(id) ON DELETE CASCADE
);

create table room_members (
    room_id uuid REFERENCES rooms(id) ON DELETE CASCADE,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (room_id, user_id)
);

-- +goose Down
drop table rooms;
drop table room_members;
