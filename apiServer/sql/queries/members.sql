-- name: AddNewRoomMember :one
insert into room_members
    (room_id, user_id)
        values ($1, $2) returning *;

-- name: GetExistingPerson :one
select * from room_members where room_id=$1 and user_id=$2 limit 1;

-- name: RemoveExistingPersonFromRoom :one
delete from room_members where room_id=$1 and user_id=$2 returning *;

-- name: UserInRoom :one
select * from room_members where room_id=$1 and user_id=$2;

-- name: GetAllMembersOfRoom :many
select rm.user_id, u.username from room_members rm join users u on rm.user_id = u.id where room_id=$1;

