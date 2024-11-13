package controllers

import (
	"database/sql"
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

func (apiCfg *ApiConf) AddMembersToRoom(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}

	validate := validator.New()
	var addMembersToRoomParams validators.AddMemberToRoomValidators

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&addMembersToRoomParams)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Invalid json object provided %v", err.Error()), []string{})
		return
	}

	addMembersToRoomParams.RoomId = strings.TrimSpace(addMembersToRoomParams.RoomId)
	addMembersToRoomParams.UserId = strings.TrimSpace(addMembersToRoomParams.UserId)

	err = validate.Struct(addMembersToRoomParams)
	if err != nil {
		var errMessageArray []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			errMessageArray = append(errMessageArray, fmt.Sprintf("Validation failed on %v at the tag %v", fieldError.Field(), fieldError.Tag()))
		}
		helpers.RespondWithError(w, 400, "Validation errors", errMessageArray)
		return
	}
	roomId, err := uuid.Parse(addMembersToRoomParams.RoomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid roomid", []string{})
		return
	}
	userToBeAdded, err := uuid.Parse(addMembersToRoomParams.UserId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid userid", []string{})
		return
	}
	if user.ID == userToBeAdded {
		helpers.RespondWithError(w, 400, "Admin cannot add himself as member", []string{})
		return
	}

	existingRoom, err := apiCfg.DB.GetRoomFromId(r.Context(), roomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the room", []string{})
		return
	}
	if existingRoom.AdminID.UUID != user.ID {
		helpers.RespondWithError(w, 400, "Only admins can perform this task", []string{})
		return
	}
	existinUser, err := apiCfg.DB.GetUseFromId(r.Context(), userToBeAdded)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the user", []string{})
		return
	}
	_, err = apiCfg.DB.GetExistingPerson(r.Context(), database.GetExistingPersonParams{
		RoomID: roomId,
		UserID: userToBeAdded,
	})
	if err != nil && err != sql.ErrNoRows {
		helpers.RespondWithError(w, 400, "User might already be present in the room or there was an issue talking to the database", []string{})
		return
	}
	_, err = apiCfg.DB.AddNewRoomMember(r.Context(), database.AddNewRoomMemberParams{
		UserID: existinUser.ID,
		RoomID: existingRoom.ID,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue adding the user to the room", []string{})
		return
	}
	helpers.RespondWithJSON(w, 200, map[string]interface{}{})
}
