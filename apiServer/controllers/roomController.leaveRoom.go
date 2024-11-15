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

func (apiCfg *ApiConf) LeaveTheRoom(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}

	validate := validator.New()
	var leaveRoomParams validators.LeaveRoomValidators

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&leaveRoomParams)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Invalid json object provided %v", err.Error()), []string{})
		return
	}

	leaveRoomParams.RoomId = strings.TrimSpace(leaveRoomParams.RoomId)

	err = validate.Struct(leaveRoomParams)
	if err != nil {
		var errMessageArray []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			errMessageArray = append(errMessageArray, fmt.Sprintf("Validation failed on %v at the tag %v", fieldError.Field(), fieldError.Tag()))
		}
		helpers.RespondWithError(w, 400, "Validation errors", errMessageArray)
		return
	}
	roomId, err := uuid.Parse(leaveRoomParams.RoomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid roomid", []string{})
		return
	}
	_, err = apiCfg.DB.RemoveExistingPersonFromRoom(r.Context(), database.RemoveExistingPersonFromRoomParams{
		RoomID: roomId,
		UserID: user.ID,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Either room not found or you are not part of it", []string{})
		return
	}
	helpers.RespondWithJSON(w, 200, map[string]string{})
}
