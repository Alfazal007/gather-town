-- name: FindExistingRoom :one
select room_name, id from rooms
    where room_name=$1 and admin_id=$2 limit 1;

-- name: AddNewRoom :one
insert into rooms 
    (id, room_name, admin_id)
        values ($1, $2, $3) returning *;

-- name: GetRoomFromId :one
select * from rooms where id=$1 limit 1;

-- name: DeleteRoomFromId :one
delete from rooms where id=$1 returning *;

