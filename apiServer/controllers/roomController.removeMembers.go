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

func (apiCfg *ApiConf) RemoveMemberFromTheRoomByAdmin(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}

	validate := validator.New()
	var removeMembersFromRoomParams validators.RemoveMemberFromRoomValidators

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&removeMembersFromRoomParams)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Invalid json object provided %v", err.Error()), []string{})
		return
	}

	removeMembersFromRoomParams.RoomId = strings.TrimSpace(removeMembersFromRoomParams.RoomId)
	removeMembersFromRoomParams.UserId = strings.TrimSpace(removeMembersFromRoomParams.UserId)

	err = validate.Struct(removeMembersFromRoomParams)
	if err != nil {
		var errMessageArray []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			errMessageArray = append(errMessageArray, fmt.Sprintf("Validation failed on %v at the tag %v", fieldError.Field(), fieldError.Tag()))
		}
		helpers.RespondWithError(w, 400, "Validation errors", errMessageArray)
		return
	}
	roomId, err := uuid.Parse(removeMembersFromRoomParams.RoomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid roomid", []string{})
		return
	}
	userToBeRemoved, err := uuid.Parse(removeMembersFromRoomParams.UserId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid userid", []string{})
		return
	}
	if user.ID == userToBeRemoved {
		helpers.RespondWithError(w, 400, "If you are the admin, then delete the room or else just leave the room", []string{})
		return
	}
	roomFromTheDatabase, err := apiCfg.DB.GetRoomFromId(r.Context(), roomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the room", []string{})
		return
	}
	if !roomFromTheDatabase.AdminID.Valid || roomFromTheDatabase.AdminID.UUID != user.ID {
		helpers.RespondWithError(w, 400, "You are probably not the admin of this room", []string{})
		return
	}
	// check if this room has the user specifed
	_, err = apiCfg.DB.GetExistingPerson(r.Context(), database.GetExistingPersonParams{
		RoomID: roomId,
		UserID: userToBeRemoved,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Maybe the user is not in the room", []string{})
		return
	}
	_, err = apiCfg.DB.RemoveExistingPersonFromRoom(r.Context(), database.RemoveExistingPersonFromRoomParams{
		RoomID: roomId,
		UserID: userToBeRemoved,
	})
	helpers.RespondWithJSON(w, 200, map[string]string{})
}
