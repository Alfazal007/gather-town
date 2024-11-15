package validators

type LeaveRoomValidators struct {
	RoomId string `json:"roomId" validate:"required"`
}
