// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: rooms.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addNewRoom = `-- name: AddNewRoom :one
insert into rooms 
    (id, room_name, admin_id)
        values ($1, $2, $3) returning id, room_name, admin_id
`

type AddNewRoomParams struct {
	ID       uuid.UUID
	RoomName string
	AdminID  uuid.NullUUID
}

func (q *Queries) AddNewRoom(ctx context.Context, arg AddNewRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, addNewRoom, arg.ID, arg.RoomName, arg.AdminID)
	var i Room
	err := row.Scan(&i.ID, &i.RoomName, &i.AdminID)
	return i, err
}

const findExistingRoom = `-- name: FindExistingRoom :one
select room_name, id from rooms
    where room_name=$1 and admin_id=$2 limit 1
`

type FindExistingRoomParams struct {
	RoomName string
	AdminID  uuid.NullUUID
}

type FindExistingRoomRow struct {
	RoomName string
	ID       uuid.UUID
}

func (q *Queries) FindExistingRoom(ctx context.Context, arg FindExistingRoomParams) (FindExistingRoomRow, error) {
	row := q.db.QueryRowContext(ctx, findExistingRoom, arg.RoomName, arg.AdminID)
	var i FindExistingRoomRow
	err := row.Scan(&i.RoomName, &i.ID)
	return i, err
}
