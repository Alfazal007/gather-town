// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: members.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addNewRoomMember = `-- name: AddNewRoomMember :one
insert into room_members
    (room_id, user_id)
        values ($1, $2) returning room_id, user_id
`

type AddNewRoomMemberParams struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) AddNewRoomMember(ctx context.Context, arg AddNewRoomMemberParams) (RoomMember, error) {
	row := q.db.QueryRowContext(ctx, addNewRoomMember, arg.RoomID, arg.UserID)
	var i RoomMember
	err := row.Scan(&i.RoomID, &i.UserID)
	return i, err
}

const getAllMembersOfRoom = `-- name: GetAllMembersOfRoom :many
select rm.user_id, u.username from room_members rm join users u on rm.user_id = u.id where room_id=$1
`

type GetAllMembersOfRoomRow struct {
	UserID   uuid.UUID
	Username string
}

func (q *Queries) GetAllMembersOfRoom(ctx context.Context, roomID uuid.UUID) ([]GetAllMembersOfRoomRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMembersOfRoom, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllMembersOfRoomRow
	for rows.Next() {
		var i GetAllMembersOfRoomRow
		if err := rows.Scan(&i.UserID, &i.Username); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExistingPerson = `-- name: GetExistingPerson :one
select room_id, user_id from room_members where room_id=$1 and user_id=$2 limit 1
`

type GetExistingPersonParams struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) GetExistingPerson(ctx context.Context, arg GetExistingPersonParams) (RoomMember, error) {
	row := q.db.QueryRowContext(ctx, getExistingPerson, arg.RoomID, arg.UserID)
	var i RoomMember
	err := row.Scan(&i.RoomID, &i.UserID)
	return i, err
}

const removeExistingPersonFromRoom = `-- name: RemoveExistingPersonFromRoom :one
delete from room_members where room_id=$1 and user_id=$2 returning room_id, user_id
`

type RemoveExistingPersonFromRoomParams struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) RemoveExistingPersonFromRoom(ctx context.Context, arg RemoveExistingPersonFromRoomParams) (RoomMember, error) {
	row := q.db.QueryRowContext(ctx, removeExistingPersonFromRoom, arg.RoomID, arg.UserID)
	var i RoomMember
	err := row.Scan(&i.RoomID, &i.UserID)
	return i, err
}

const userInRoom = `-- name: UserInRoom :one
select room_id, user_id from room_members where room_id=$1 and user_id=$2
`

type UserInRoomParams struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) UserInRoom(ctx context.Context, arg UserInRoomParams) (RoomMember, error) {
	row := q.db.QueryRowContext(ctx, userInRoom, arg.RoomID, arg.UserID)
	var i RoomMember
	err := row.Scan(&i.RoomID, &i.UserID)
	return i, err
}
