package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/Alfazal007/gather-town/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (apiCfg *ApiConf) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}
	validate := validator.New()
	var deleteRoomParams validators.DeleteRoomValidate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&deleteRoomParams)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue decoding the request body", []string{})
		return
	}
	deleteRoomParams.RoomId = strings.TrimSpace(deleteRoomParams.RoomId)
	err = validate.Struct(deleteRoomParams)
	if err != nil {
		var errMessageArray []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			errMessageArray = append(errMessageArray, fmt.Sprintf("Validation failed on %v at the tag %v", fieldError.Field(), fieldError.Tag()))
		}
		helpers.RespondWithError(w, 400, "Validation errors", errMessageArray)
		return
	}
	roomId, err := uuid.Parse(deleteRoomParams.RoomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Roomid is not valid", []string{})
		return
	}
	roomToBeDeleted, err := apiCfg.DB.GetRoomFromId(r.Context(), roomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Room couldn't be found", []string{})
		return
	}
	if !roomToBeDeleted.AdminID.Valid || roomToBeDeleted.AdminID.UUID != user.ID {
		helpers.RespondWithError(w, 400, "You are not the owner of the room", []string{})
		return
	}
	_, err = apiCfg.DB.DeleteRoomFromId(r.Context(), roomId)
	helpers.RespondWithJSON(w, 200, map[string]string{})
}
