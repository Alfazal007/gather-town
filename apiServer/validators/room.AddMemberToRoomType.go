package validators

type AddMemberToRoomValidatorsViaName struct {
	RoomId   string `json:"roomid" validate:"required"`
	UserName string `json:"username" validate:"required"`
}
