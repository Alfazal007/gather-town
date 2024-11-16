package types

import (
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/google/uuid"
)

type CustomRoomData struct {
	RoomID      string
	RoomName    string
	AdminID     string
	RoomMembers []uuid.UUID
}

func ReturnRoomInformationData(roomFromDatabase database.Room, roomMembers []uuid.UUID) CustomRoomData {
	return CustomRoomData{
		RoomID:      roomFromDatabase.ID.String(),
		RoomName:    roomFromDatabase.RoomName,
		AdminID:     roomFromDatabase.AdminID.UUID.String(),
		RoomMembers: roomMembers,
	}
}
