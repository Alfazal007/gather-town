package types

import (
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/google/uuid"
)

type RoomMemberData struct {
	Id       uuid.UUID `json:"roomMemberId"`
	Username string    `json:"roomMemberUsername"`
}

type CustomRoomData struct {
	RoomID      string           `json:"roomId"`
	RoomName    string           `json:"roomName"`
	AdminID     string           `json:"adminId"`
	RoomMembers []RoomMemberData `json:"roomMembers"`
}

func ReturnRoomInformationData(roomFromDatabase database.Room, roomMembers []database.GetAllMembersOfRoomRow) CustomRoomData {
	dataToBeReturned := make([]RoomMemberData, 0)
	for i := 0; i < len(roomMembers); i++ {
		dataToBeReturned = append(dataToBeReturned, RoomMemberData{
			Id:       roomMembers[i].UserID,
			Username: roomMembers[i].Username,
		})
	}
	return CustomRoomData{
		RoomID:      roomFromDatabase.ID.String(),
		RoomName:    roomFromDatabase.RoomName,
		AdminID:     roomFromDatabase.AdminID.UUID.String(),
		RoomMembers: dataToBeReturned,
	}
}
