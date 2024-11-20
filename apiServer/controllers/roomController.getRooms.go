package controllers

import (
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/Alfazal007/gather-town/types"
	"github.com/google/uuid"
)

func (apiCfg *ApiConf) GetRoomsUserIsPartOf(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}
	userId := uuid.NullUUID{UUID: user.ID, Valid: true}
	roomsUserIsPartOf, err := apiCfg.DB.GetRoomsOfUser(r.Context(), userId)
	if err != nil {
		helpers.RespondWithJSON(w, 200, []types.CustomRoom{})
		return
	}
	roomsToBeReturned := make([]types.CustomRoom, 0)
	for i := 0; i < len(roomsUserIsPartOf); i++ {
		roomsToBeReturned = append(roomsToBeReturned, types.ReturnCreatedRoom(roomsUserIsPartOf[i]))
	}

	helpers.RespondWithJSON(w, 200, roomsToBeReturned)
}
