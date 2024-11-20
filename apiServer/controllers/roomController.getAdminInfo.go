package controllers

import (
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *ApiConf) GetAdminInformation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database", []string{})
		return
	}
	roomId := chi.URLParam(r, "roomId")
	roomIdInUUID, err := uuid.Parse(roomId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid room id", []string{})
		return
	}
	roomData, err := apiCfg.DB.GetRoomFromId(r.Context(), roomIdInUUID)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the room", []string{})
		return
	}
	if !roomData.AdminID.Valid {
		helpers.RespondWithError(w, 400, "Issue finding the room", []string{})
		return
	}
	roomMembers, err := apiCfg.DB.GetAllMembersOfRoom(r.Context(), roomIdInUUID)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the room data", []string{})
		return
	}

	type AdminData struct {
		AdminName string `json:"adminName"`
		AdminId   string `json:"adminId"`
	}

	adminData, err := apiCfg.DB.GetUseFromId(r.Context(), roomData.AdminID.UUID)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue finding the room data", []string{})
		return
	}

	if user.ID == roomData.AdminID.UUID || isMember(user.ID, roomMembers) {
		helpers.RespondWithJSON(w, 200, AdminData{
			AdminName: adminData.Username,
			AdminId:   roomData.AdminID.UUID.String(),
		})
		return
	}
	helpers.RespondWithError(w, 400, "Cannot respond with data", []string{})
}

func isMember(userId uuid.UUID, roomMembers []database.GetAllMembersOfRoomRow) bool {
	for i := 0; i < len(roomMembers); i++ {
		if roomMembers[i].UserID == userId {
			return true
		}
	}
	return false
}
