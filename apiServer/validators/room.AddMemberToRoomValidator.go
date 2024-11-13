package validators

type AddMemberToRoomValidators struct {
	RoomId string `json:"roomId" validate:"required"`
	UserId string `json:"userId" validate:"required"`
}
