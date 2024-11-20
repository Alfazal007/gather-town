package controllers

import (
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/Alfazal007/gather-town/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *ApiConf) GetRoomInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}

	roomIdString := chi.URLParam(r, "roomId")
	roomId, err := uuid.Parse(roomIdString)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid roomid", []string{})
		return
	}

	roomData, err := apiCfg.DB.GetRoomFromId(r.Context(), roomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the room", []string{})
		return
	}

	if roomData.AdminID.Valid && roomData.AdminID.UUID != user.ID {
		_, err = apiCfg.DB.GetExistingPerson(r.Context(), database.GetExistingPersonParams{
			RoomID: roomId,
			UserID: user.ID,
		})
		if err != nil {
			helpers.RespondWithError(w, 400, "Either room not found or you are not part of it", []string{})
			return
		}
	}

	roomMembers, err := apiCfg.DB.GetAllMembersOfRoom(r.Context(), roomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the data of the required room", []string{})
		return
	}
	helpers.RespondWithJSON(w, 200, types.ReturnRoomInformationData(roomData, roomMembers))
}
