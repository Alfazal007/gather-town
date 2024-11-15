package validators

type RemoveMemberFromRoomValidators struct {
	RoomId string `json:"roomId" validate:"required"`
	UserId string `json:"userId" validate:"required"`
}
