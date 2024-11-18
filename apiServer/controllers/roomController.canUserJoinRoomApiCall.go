package controllers

import (
	"fmt"
	"net/http"

	"github.com/Alfazal007/gather-town/helpers"
	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/Alfazal007/gather-town/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (apiCfg *ApiConf) UserCanJoinRoomApiCall(w http.ResponseWriter, r *http.Request) {
	usernameRequested := chi.URLParam(r, "username")
	roomIdString := chi.URLParam(r, "roomId")
	roomId, err := uuid.Parse(roomIdString)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid roomid", []string{})
		return
	}
	jwtToken := chi.URLParam(r, "token")
	jwtSecret := utils.LoadEnvVariables().AccessTokenSecret
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid token here", []string{})
		return
	}
	if !token.Valid {
		helpers.RespondWithError(w, 400, "Invalid token", []string{})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.RespondWithError(w, 400, "Invalid claims login again", []string{})
		return
	}

	username := claims["username"].(string)
	id := claims["user_id"].(string)

	if username != usernameRequested {
		helpers.RespondWithError(w, 400, "Invalid username trying to be accessed", []string{})
		return
	}

	user, err := apiCfg.DB.GetUserByName(r.Context(), username)
	if err != nil {
		helpers.RespondWithError(w, 400, "Some manpulation done with the token", []string{})
		return
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		helpers.RespondWithError(w, 400, "Some manpulation done with the token", []string{})
		return
	}
	if idUUID != user.ID {
		helpers.RespondWithError(w, 400, "Some manipulations done with the token try again", []string{})
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
		helpers.RespondWithJSON(w, 400, CanJoin{CanJoin: false})
		return
	}
	helpers.RespondWithJSON(w, 200, CanJoin{CanJoin: true})
}
