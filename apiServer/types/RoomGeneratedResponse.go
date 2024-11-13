package types

import "github.com/Alfazal007/gather-town/internal/database"

type CustomRoom struct {
	ID       string
	RoomName string
	AdminID  string
}

func ReturnCreatedRoom(roomFromDatabase database.Room) CustomRoom {
	return CustomRoom{
		ID:       roomFromDatabase.ID.String(),
		RoomName: roomFromDatabase.RoomName,
		AdminID:  roomFromDatabase.AdminID.UUID.String(),
	}
}
