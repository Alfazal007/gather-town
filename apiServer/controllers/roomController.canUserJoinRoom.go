package controllers

import (
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CanJoin struct {
	CanJoin bool `json:"canJoin"`
}

func (apiCfg *ApiConf) UserCanJoinRoom(w http.ResponseWriter, r *http.Request) {
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
	if roomData.AdminID.Valid && roomData.AdminID.UUID == user.ID {
		helpers.RespondWithJSON(w, 200, CanJoin{CanJoin: true})
		return
	}
	_, err = apiCfg.DB.UserInRoom(r.Context(), database.UserInRoomParams{
		RoomID: roomId,
		UserID: user.ID,
	})
	if err != nil {
		helpers.RespondWithJSON(w, 200, CanJoin{CanJoin: false})
		return
	}
	helpers.RespondWithJSON(w, 200, CanJoin{CanJoin: true})
}
