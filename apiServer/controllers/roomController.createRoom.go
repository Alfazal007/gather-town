package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/Alfazal007/gather-town/types"
	"github.com/Alfazal007/gather-town/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (apiCfg *ApiConf) CreateRoom(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "User not logged in", []string{})
		return
	}
	validate := validator.New()
	var createRoomParams validators.CreateRoomValidate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createRoomParams)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue decoding the request body", []string{})
		return
	}
	createRoomParams.Name = strings.TrimSpace(createRoomParams.Name)
	err = validate.Struct(createRoomParams)
	if err != nil {
		var errMessageArray []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			errMessageArray = append(errMessageArray, fmt.Sprintf("Validation failed on %v at the tag %v", fieldError.Field(), fieldError.Tag()))
		}
		helpers.RespondWithError(w, 400, "Validation errors", errMessageArray)
		return
	}
	// cheeck if the database already has the same named room from the same user
	_, err = apiCfg.DB.FindExistingRoom(r.Context(), database.FindExistingRoomParams{
		RoomName: createRoomParams.Name,
		AdminID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	if err != nil && err != sql.ErrNoRows {
		helpers.RespondWithError(w, 400, "Try to change to room name or try again later", []string{})
		return
	}
	if err == nil {
		helpers.RespondWithError(w, 400, "Try to change to room name", []string{})
		return
	}
	createdRoom, err := apiCfg.DB.AddNewRoom(r.Context(), database.AddNewRoomParams{
		ID:       uuid.New(),
		RoomName: createRoomParams.Name,
		AdminID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error creating the room", []string{})
		return
	}
	helpers.RespondWithJSON(w, 201, types.ReturnCreatedRoom(createdRoom))
}
