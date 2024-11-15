package validators

type DeleteRoomValidate struct {
	RoomId string `json:"roomId" validate:"required"`
}
