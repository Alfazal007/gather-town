-- name: AddNewRoomMember :one
insert into room_members
    (room_id, user_id)
        values ($1, $2) returning *;

-- name: GetExistingPerson :one
select * from room_members where room_id=$1 and user_id=$2 limit 1;
