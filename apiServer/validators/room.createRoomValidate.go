package validators

type CreateRoomValidate struct {
	Name string `json:"name" validate:"required,min=4,max=20"`
}
